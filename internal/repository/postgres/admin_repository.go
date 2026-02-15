package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"coviar_backend/internal/domain"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// GetStats obtiene estadísticas generales del sistema
func (r *AdminRepository) GetStats(ctx context.Context) (*domain.AdminStatsResponse, error) {
	stats := &domain.AdminStatsResponse{}

	// Total de bodegas
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM bodegas
	`).Scan(&stats.TotalBodegas)
	if err != nil {
		return nil, fmt.Errorf("error counting bodegas: %w", err)
	}

	// Evaluaciones completadas
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM autoevaluaciones
		WHERE estado = 'COMPLETADA'
	`).Scan(&stats.EvaluacionesCompletadas)
	if err != nil {
		return nil, fmt.Errorf("error counting completed evaluaciones: %w", err)
	}

	// Promedio de sostenibilidad (establecido en 0 porque porcentaje no está en la tabla)
	stats.PromedioSostenibilidad = 0

	return stats, nil
}

// GetAllEvaluaciones obtiene todas las autoevaluaciones con filtros opcionales
func (r *AdminRepository) GetAllEvaluaciones(ctx context.Context, estado string, idBodega int) ([]domain.EvaluacionListItem, error) {
	query := `
		SELECT
			a.id_autoevaluacion,
			a.id_bodega,
			b.nombre_fantasia,
			b.razon_social,
			a.estado,
			a.fecha_inicio,
			a.fecha_fin,
			COALESCE(resp.nombre || ' ' || resp.apellido, 'Sin responsable') as responsable
		FROM autoevaluaciones a
		INNER JOIN bodegas b ON a.id_bodega = b.id_bodega
		LEFT JOIN cuentas c ON b.id_bodega = c.id_bodega
		LEFT JOIN responsables resp ON c.id_cuenta = resp.id_cuenta AND resp.activo = true
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	// Filtro por estado
	if estado != "" && estado != "TODOS" {
		query += fmt.Sprintf(" AND a.estado = $%d", argCount)
		args = append(args, strings.ToUpper(estado))
		argCount++
	}

	// Filtro por bodega
	if idBodega > 0 {
		query += fmt.Sprintf(" AND a.id_bodega = $%d", argCount)
		args = append(args, idBodega)
		argCount++
	}

	query += " ORDER BY a.fecha_inicio DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying evaluaciones: %w", err)
	}
	defer rows.Close()

	var evaluaciones []domain.EvaluacionListItem
	for rows.Next() {
		var eval domain.EvaluacionListItem
		err := rows.Scan(
			&eval.IDAutoevaluacion,
			&eval.IDBodega,
			&eval.NombreBodega,
			&eval.RazonSocial,
			&eval.Estado,
			&eval.FechaInicio,
			&eval.FechaFin,
			&eval.Responsable,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning evaluacion: %w", err)
		}

		// Porcentaje se establece en nil porque no está en la tabla
		eval.Porcentaje = nil

		evaluaciones = append(evaluaciones, eval)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return evaluaciones, nil
}
