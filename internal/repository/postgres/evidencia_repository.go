package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
)

type EvidenciaRepository struct {
	db *sql.DB
}

func NewEvidenciaRepository(db *sql.DB) repository.EvidenciaRepository {
	return &EvidenciaRepository{db: db}
}

func (r *EvidenciaRepository) Create(ctx context.Context, tx repository.Transaction, evidencia *domain.Evidencia) (int, error) {
	query := `
		INSERT INTO evidencias (id_respuesta, nombre, ubicacion)
		VALUES ($1, $2, $3)
		RETURNING id_evidencia
	`

	var id int
	err := r.db.QueryRowContext(ctx, query, evidencia.IDRespuesta, evidencia.Nombre, evidencia.Ubicacion).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating evidencia: %w", err)
	}

	evidencia.ID = id
	return id, nil
}

func (r *EvidenciaRepository) FindByRespuesta(ctx context.Context, idRespuesta int) (*domain.Evidencia, error) {
	query := `
		SELECT id_evidencia, id_respuesta, nombre, ubicacion
		FROM evidencias
		WHERE id_respuesta = $1
	`

	evidencia := &domain.Evidencia{}
	err := r.db.QueryRowContext(ctx, query, idRespuesta).Scan(&evidencia.ID, &evidencia.IDRespuesta, &evidencia.Nombre, &evidencia.Ubicacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying evidencia: %w", err)
	}

	return evidencia, nil
}

func (r *EvidenciaRepository) FindByAutoevaluacion(ctx context.Context, idAutoevaluacion int) ([]*domain.Evidencia, error) {
	query := `
		SELECT e.id_evidencia, e.id_respuesta, e.nombre, e.ubicacion
		FROM evidencias e
		INNER JOIN respuestas r ON e.id_respuesta = r.id_respuesta
		WHERE r.id_autoevaluacion = $1
	`

	rows, err := r.db.QueryContext(ctx, query, idAutoevaluacion)
	if err != nil {
		return nil, fmt.Errorf("error querying evidencias: %w", err)
	}
	defer rows.Close()

	var evidencias []*domain.Evidencia
	for rows.Next() {
		ev := &domain.Evidencia{}
		if err := rows.Scan(&ev.ID, &ev.IDRespuesta, &ev.Nombre, &ev.Ubicacion); err != nil {
			return nil, fmt.Errorf("error scanning evidencia: %w", err)
		}
		evidencias = append(evidencias, ev)
	}

	return evidencias, rows.Err()
}

func (r *EvidenciaRepository) Delete(ctx context.Context, tx repository.Transaction, id int) error {
	query := `DELETE FROM evidencias WHERE id_evidencia = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting evidencia: %w", err)
	}

	return nil
}

func (r *EvidenciaRepository) CountEvidenciasByAutoevaluacion(ctx context.Context, idAutoevaluacion int) (int, error) {
	query := `
		SELECT COUNT(DISTINCT e.id_evidencia)
		FROM evidencias e
		INNER JOIN respuestas r ON e.id_respuesta = r.id_respuesta
		WHERE r.id_autoevaluacion = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, idAutoevaluacion).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting evidencias: %w", err)
	}

	return count, nil
}
