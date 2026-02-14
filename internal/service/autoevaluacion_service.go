package service

import (
	"context"
	"fmt"
	"os"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
)

type AutoevaluacionService struct {
	autoevaluacionRepo repository.AutoevaluacionRepository
	segmentoRepo       repository.SegmentoRepository
	capituloRepo       repository.CapituloRepository
	indicadorRepo      repository.IndicadorRepository
	nivelRespuestaRepo repository.NivelRespuestaRepository
	respuestaRepo      repository.RespuestaRepository
	evidenciaRepo      repository.EvidenciaRepository
}

func NewAutoevaluacionService(
	autoevaluacionRepo repository.AutoevaluacionRepository,
	segmentoRepo repository.SegmentoRepository,
	capituloRepo repository.CapituloRepository,
	indicadorRepo repository.IndicadorRepository,
	nivelRespuestaRepo repository.NivelRespuestaRepository,
	respuestaRepo repository.RespuestaRepository,
	evidenciaRepo repository.EvidenciaRepository,
) *AutoevaluacionService {
	return &AutoevaluacionService{
		autoevaluacionRepo: autoevaluacionRepo,
		segmentoRepo:       segmentoRepo,
		capituloRepo:       capituloRepo,
		indicadorRepo:      indicadorRepo,
		nivelRespuestaRepo: nivelRespuestaRepo,
		respuestaRepo:      respuestaRepo,
		evidenciaRepo:      evidenciaRepo,
	}
}

// CreateAutoevaluacion crea una nueva autoevaluación para una bodega
/*func (s *AutoevaluacionService) CreateAutoevaluacion(ctx context.Context, idBodega int) (*domain.Autoevaluacion, error) {
	auto := &domain.Autoevaluacion{
		IDBodega: idBodega,
	}

	id, err := s.autoevaluacionRepo.Create(ctx, nil, auto)
	if err != nil {
		return nil, fmt.Errorf("error creating autoevaluacion: %w", err)
	}

	auto.ID = id
	return auto, nil
}*/

// CreateAutoevaluacion crea una nueva autoevaluación para una bodega o retorna la pendiente
func (s *AutoevaluacionService) CreateAutoevaluacion(ctx context.Context, idBodega int) (*domain.AutoevaluacionPendienteResponse, error) {
	// Verificar si ya existe una autoevaluación pendiente para esta bodega
	autoPendiente, err := s.autoevaluacionRepo.FindPendienteByBodega(ctx, idBodega)
	if err != nil {
		return nil, fmt.Errorf("error checking for pending autoevaluacion: %w", err)
	}

	// Si existe una autoevaluación pendiente, retornarla con sus respuestas
	if autoPendiente != nil {
		// Obtener las respuestas existentes
		respuestas, err := s.respuestaRepo.FindByAutoevaluacion(ctx, autoPendiente.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting respuestas: %w", err)
		}

		// Convertir respuestas al formato DTO
		respuestasDTO := make([]domain.GuardarRespuestaRequest, len(respuestas))
		for i, resp := range respuestas {
			respuestasDTO[i] = domain.GuardarRespuestaRequest{
				IDIndicador:      resp.IDIndicador,
				IDNivelRespuesta: resp.IDNivelRespuesta,
			}
		}

		return &domain.AutoevaluacionPendienteResponse{
			AutoevaluacionPendiente: autoPendiente,
			Respuestas:              respuestasDTO,
			Mensaje:                 "Ya existe una autoevaluación pendiente. Para crear una nueva, primero debe cancelar la actual.",
		}, nil
	}

	// No existe autoevaluación pendiente, crear una nueva
	auto := &domain.Autoevaluacion{
		IDBodega: idBodega,
	}

	id, err := s.autoevaluacionRepo.Create(ctx, nil, auto)
	if err != nil {
		return nil, fmt.Errorf("error creating autoevaluacion: %w", err)
	}

	auto.ID = id
	return &domain.AutoevaluacionPendienteResponse{
		AutoevaluacionPendiente: auto,
		Respuestas:              []domain.GuardarRespuestaRequest{},
		Mensaje:                 "Autoevaluación creada correctamente",
	}, nil
}

// GetSegmentos obtiene todos los segmentos disponibles
func (s *AutoevaluacionService) GetSegmentos(ctx context.Context) ([]*domain.Segmento, error) {
	segmentos, err := s.segmentoRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting segmentos: %w", err)
	}
	return segmentos, nil
}

// SeleccionarSegmento selecciona un segmento para la autoevaluación
func (s *AutoevaluacionService) SeleccionarSegmento(ctx context.Context, idAutoevaluacion int, idSegmento int) error {
	// Verificar que la autoevaluación existe
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil {
		return domain.ErrNotFound
	}

	// Verificar que el segmento existe
	seg, err := s.segmentoRepo.FindByID(ctx, idSegmento)
	if err != nil {
		return fmt.Errorf("error finding segmento: %w", err)
	}

	if seg == nil {
		return domain.ErrNotFound
	}

	// Actualizar la autoevaluación con el segmento
	err = s.autoevaluacionRepo.UpdateSegmento(ctx, idAutoevaluacion, idSegmento)
	if err != nil {
		return fmt.Errorf("error selecting segmento: %w", err)
	}

	return nil
}

// GetEstructura obtiene la estructura del cuestionario con indicadores habilitados según el segmento
func (s *AutoevaluacionService) GetEstructura(ctx context.Context, idAutoevaluacion int) (*domain.EstructuraAutoevaluacion, error) {
	// Obtener autoevaluación
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return nil, fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil || auto.IDSegmento == nil {
		return nil, fmt.Errorf("autoevaluacion not found or segmento not selected")
	}

	// Obtener indicadores habilitados para este segmento
	habilitadosIds, err := s.indicadorRepo.FindBySegmento(ctx, *auto.IDSegmento)
	if err != nil {
		return nil, fmt.Errorf("error getting enabled indicators: %w", err)
	}

	// Convertir a mapa para búsqueda rápida
	habilitadosMap := make(map[int]bool)
	for _, id := range habilitadosIds {
		habilitadosMap[id] = true
	}

	// Obtener todos los capítulos
	capitulos, err := s.capituloRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting capitulos: %w", err)
	}

	estructura := &domain.EstructuraAutoevaluacion{
		Capitulos: make([]*domain.CapituloEstructura, 0),
	}

	for _, cap := range capitulos {
		// Obtener indicadores del capítulo
		indicadores, err := s.indicadorRepo.FindByCapitulo(ctx, cap.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting indicadores: %w", err)
		}

		capEstructura := &domain.CapituloEstructura{
			Capitulo:    cap,
			Indicadores: make([]*domain.IndicadorConHabilitacion, 0),
		}

		for _, ind := range indicadores {
			// Obtener niveles de respuesta
			niveles, err := s.nivelRespuestaRepo.FindByIndicador(ctx, ind.ID)
			if err != nil {
				return nil, fmt.Errorf("error getting niveles_respuesta: %w", err)
			}

			habilitado := habilitadosMap[ind.ID]

			indConHab := &domain.IndicadorConHabilitacion{
				Indicador:        ind,
				NivelesRespuesta: niveles,
				Habilitado:       habilitado,
			}

			capEstructura.Indicadores = append(capEstructura.Indicadores, indConHab)
		}

		estructura.Capitulos = append(estructura.Capitulos, capEstructura)
	}

	return estructura, nil
}

// GuardarRespuestas guarda las respuestas de la autoevaluación y retorna las respuestas con sus IDs
func (s *AutoevaluacionService) GuardarRespuestas(ctx context.Context, idAutoevaluacion int, respuestas []domain.GuardarRespuestaRequest) ([]*domain.Respuesta, error) {
	// Verificar que la autoevaluación existe
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return nil, fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil {
		return nil, domain.ErrNotFound
	}

	// Obtener respuestas existentes para detectar cambios de nivel
	existentes, err := s.respuestaRepo.FindByAutoevaluacion(ctx, idAutoevaluacion)
	if err != nil {
		return nil, fmt.Errorf("error getting existing respuestas: %w", err)
	}

	// Mapa de id_indicador -> respuesta existente
	existentesMap := make(map[int]*domain.Respuesta)
	for _, r := range existentes {
		existentesMap[r.IDIndicador] = r
	}

	// Guardar nuevas respuestas o actualizar existentes
	resultado := make([]*domain.Respuesta, 0, len(respuestas))
	for _, respReq := range respuestas {
		// Si la respuesta ya existía con un nivel diferente, eliminar su evidencia
		if existente, ok := existentesMap[respReq.IDIndicador]; ok {
			if existente.IDNivelRespuesta != respReq.IDNivelRespuesta {
				evidencia, err := s.evidenciaRepo.FindByRespuesta(ctx, existente.ID)
				if err == nil && evidencia != nil {
					os.Remove(evidencia.Ubicacion)
					s.evidenciaRepo.Delete(ctx, nil, evidencia.ID)
				}
			}
		}

		respuesta := &domain.Respuesta{
			IDNivelRespuesta: respReq.IDNivelRespuesta,
			IDIndicador:      respReq.IDIndicador,
			IDAutoevaluacion: idAutoevaluacion,
		}

		id, err := s.respuestaRepo.Upsert(ctx, nil, respuesta)
		if err != nil {
			return nil, fmt.Errorf("error guarding respuesta: %w", err)
		}
		respuesta.ID = id
		resultado = append(resultado, respuesta)
	}

	return resultado, nil
}

// CompletarAutoevaluacion marca la autoevaluación como completada
/*func (s *AutoevaluacionService) CompletarAutoevaluacion(ctx context.Context, idAutoevaluacion int) error {
	// Verificar que la autoevaluación existe
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil {
		return domain.ErrNotFound
	}

	// Obtener respuestas para validar que todas las preguntas fueron respondidas
	respuestas, err := s.respuestaRepo.FindByAutoevaluacion(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error getting respuestas: %w", err)
	}

	// Validación básica: debe haber al menos una respuesta
	if len(respuestas) == 0 {
		return fmt.Errorf("autoevaluacion must have at least one respuesta")
	}

	// Marcar como completada
	err = s.autoevaluacionRepo.Complete(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error completing autoevaluacion: %w", err)
	}

	return nil
}*/

// CompletarAutoevaluacion marca la autoevaluación como completada
func (s *AutoevaluacionService) CompletarAutoevaluacion(ctx context.Context, idAutoevaluacion int) error {
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil {
		return domain.ErrNotFound
	}

	// Verificar que tenga segmento seleccionado
	if auto.IDSegmento == nil {
		return fmt.Errorf("autoevaluacion must have segmento selected")
	}

	// Obtener respuestas para validar que todas las preguntas fueron respondidas
	respuestas, err := s.respuestaRepo.FindByAutoevaluacion(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error getting respuestas: %w", err)
	}

	// Validación básica: debe haber al menos una respuesta
	if len(respuestas) == 0 {
		return fmt.Errorf("autoevaluacion must have at least one respuesta")
	}

	// === VALIDACIÓN ESTRICTA DE COMPLETITUD ===
	// Obtener los indicadores requeridos para este segmento
	requiredIndicators, err := s.indicadorRepo.FindBySegmento(ctx, *auto.IDSegmento)
	if err != nil {
		return fmt.Errorf("error getting required indicators: %w", err)
	}

	// Verificar que coincida la cantidad
	if len(respuestas) != len(requiredIndicators) {
		return fmt.Errorf("autoevaluacion incomplete: expected %d answers, got %d", len(requiredIndicators), len(respuestas))
	}
	// ==========================================

	// Calcular puntaje total
	puntajeTotal, err := s.respuestaRepo.CalculateTotalScore(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error calculating total score: %w", err)
	}

	// Obtener niveles de sostenibilidad para el segmento
	niveles, err := s.segmentoRepo.FindNivelesSostenibilidadBySegmento(ctx, *auto.IDSegmento)
	if err != nil {
		return fmt.Errorf("error getting niveles sostenibilidad: %w", err)
	}

	// Si no hay niveles configurados, completar sin asignar nivel
	if len(niveles) == 0 {
		// Marcar como completada solo con puntaje (sin nivel de sostenibilidad)
		err = s.autoevaluacionRepo.CompleteWithScore(ctx, idAutoevaluacion, puntajeTotal, 0)
		if err != nil {
			return fmt.Errorf("error completing autoevaluacion: %w", err)
		}
		return nil
	}

	// Determinar el nivel de sostenibilidad según el puntaje
	var nivelAsignado *domain.NivelSostenibilidad
	for _, nivel := range niveles {
		if puntajeTotal >= nivel.MinPuntaje && puntajeTotal <= nivel.MaxPuntaje {
			nivelAsignado = nivel
			break
		}
	}

	// Si no encontró nivel exacto, asignar el más cercano
	if nivelAsignado == nil {
		// Buscar el nivel más alto si el puntaje supera todos los rangos
		for _, nivel := range niveles {
			if puntajeTotal >= nivel.MinPuntaje {
				nivelAsignado = nivel
			}
		}
		// Si aún no hay nivel, usar el primero (puntaje menor al mínimo)
		if nivelAsignado == nil && len(niveles) > 0 {
			nivelAsignado = niveles[0]
		}
	}

	var idNivelSostenibilidad int
	if nivelAsignado != nil {
		idNivelSostenibilidad = nivelAsignado.ID
	}

	// Marcar como completada con puntaje y nivel de sostenibilidad
	err = s.autoevaluacionRepo.CompleteWithScore(ctx, idAutoevaluacion, puntajeTotal, idNivelSostenibilidad)
	if err != nil {
		return fmt.Errorf("error completing autoevaluacion: %w", err)
	}

	return nil
}

// CancelarAutoevaluacion marca la autoevaluación como cancelada
func (s *AutoevaluacionService) CancelarAutoevaluacion(ctx context.Context, idAutoevaluacion int) error {
	// Verificar que la autoevaluación existe
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil {
		return domain.ErrNotFound
	}

	// Verificar que esté en estado PENDIENTE
	if auto.Estado != domain.EstadoPendiente {
		return fmt.Errorf("solo se pueden cancelar autoevaluaciones en estado PENDIENTE")
	}

	// Marcar como cancelada
	err = s.autoevaluacionRepo.Cancel(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error canceling autoevaluacion: %w", err)
	}

	return nil
}
