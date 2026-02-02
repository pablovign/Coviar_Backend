package service

import (
	"context"
	"fmt"

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
}

func NewAutoevaluacionService(
	autoevaluacionRepo repository.AutoevaluacionRepository,
	segmentoRepo repository.SegmentoRepository,
	capituloRepo repository.CapituloRepository,
	indicadorRepo repository.IndicadorRepository,
	nivelRespuestaRepo repository.NivelRespuestaRepository,
	respuestaRepo repository.RespuestaRepository,
) *AutoevaluacionService {
	return &AutoevaluacionService{
		autoevaluacionRepo: autoevaluacionRepo,
		segmentoRepo:       segmentoRepo,
		capituloRepo:       capituloRepo,
		indicadorRepo:      indicadorRepo,
		nivelRespuestaRepo: nivelRespuestaRepo,
		respuestaRepo:      respuestaRepo,
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

// GuardarRespuestas guarda las respuestas de la autoevaluación
func (s *AutoevaluacionService) GuardarRespuestas(ctx context.Context, idAutoevaluacion int, respuestas []domain.GuardarRespuestaRequest) error {
	// Verificar que la autoevaluación existe
	auto, err := s.autoevaluacionRepo.FindByID(ctx, idAutoevaluacion)
	if err != nil {
		return fmt.Errorf("error finding autoevaluacion: %w", err)
	}

	if auto == nil {
		return domain.ErrNotFound
	}

	// Guardar nuevas respuestas o actualizar existentes
	for _, respReq := range respuestas {
		respuesta := &domain.Respuesta{
			IDNivelRespuesta: respReq.IDNivelRespuesta,
			IDIndicador:      respReq.IDIndicador,
			IDAutoevaluacion: idAutoevaluacion,
		}

		_, err := s.respuestaRepo.Upsert(ctx, nil, respuesta)
		if err != nil {
			return fmt.Errorf("error guarding respuesta: %w", err)
		}
	}

	return nil
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

	// Determinar el nivel de sostenibilidad según el puntaje
	var nivelAsignado *domain.NivelSostenibilidad
	for _, nivel := range niveles {
		if puntajeTotal >= nivel.MinPuntaje && puntajeTotal <= nivel.MaxPuntaje {
			nivelAsignado = nivel
			break
		}
	}

	if nivelAsignado == nil {
		return fmt.Errorf("no se encontró un nivel de sostenibilidad para el puntaje %d en el segmento %d", puntajeTotal, *auto.IDSegmento)
	}

	// Marcar como completada con puntaje y nivel de sostenibilidad
	err = s.autoevaluacionRepo.CompleteWithScore(ctx, idAutoevaluacion, puntajeTotal, nivelAsignado.ID)
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
