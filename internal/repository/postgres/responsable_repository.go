package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
)

type ResponsableRepository struct {
	db *sql.DB
}

func NewResponsableRepository(db *sql.DB) repository.ResponsableRepository {
	return &ResponsableRepository{db: db}
}

func (r *ResponsableRepository) Create(ctx context.Context, tx repository.Transaction, responsable *domain.Responsable) (int, error) {
	query := `
		INSERT INTO responsables (id_bodega, nombre, apellido, cargo, dni, activo, fecha_registro)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id_responsable
	`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		responsable.IDBodega, responsable.Nombre, responsable.Apellido, responsable.Cargo, responsable.DNI, responsable.Activo,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("error creating responsable: %w", err)
	}

	responsable.ID = id
	responsable.FechaRegistro = time.Now()
	return id, nil
}

func (r *ResponsableRepository) FindByID(ctx context.Context, id int) (*domain.Responsable, error) {
	query := `
		SELECT id_responsable, id_bodega, nombre, apellido, cargo, dni, activo, fecha_registro
		FROM responsables WHERE id_responsable = $1
	`

	responsable := &domain.Responsable{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&responsable.ID, &responsable.IDBodega, &responsable.Nombre, &responsable.Apellido,
		&responsable.Cargo, &responsable.DNI, &responsable.Activo, &responsable.FechaRegistro,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error finding responsable: %w", err)
	}

	return responsable, nil
}

func (r *ResponsableRepository) FindByBodegaID(ctx context.Context, bodegaID int) ([]*domain.Responsable, error) {
	query := `
		SELECT id_responsable, id_bodega, nombre, apellido, cargo, dni, activo, fecha_registro
		FROM responsables WHERE id_bodega = $1 ORDER BY id_responsable
	`

	rows, err := r.db.QueryContext(ctx, query, bodegaID)
	if err != nil {
		return nil, fmt.Errorf("error finding responsables: %w", err)
	}
	defer rows.Close()

	var responsables []*domain.Responsable
	for rows.Next() {
		responsable := &domain.Responsable{}
		err := rows.Scan(
			&responsable.ID, &responsable.IDBodega, &responsable.Nombre, &responsable.Apellido,
			&responsable.Cargo, &responsable.DNI, &responsable.Activo, &responsable.FechaRegistro,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning responsable: %w", err)
		}
		responsables = append(responsables, responsable)
	}

	return responsables, rows.Err()
}

func (r *ResponsableRepository) Update(ctx context.Context, tx repository.Transaction, responsable *domain.Responsable) error {
	query := `
		UPDATE responsable
		SET nombre = $1, apellido = $2, cargo = $3, dni = $4, activo = $5
		WHERE id_responsable = $6
	`

	_, err := r.db.ExecContext(ctx, query,
		responsable.Nombre, responsable.Apellido, responsable.Cargo, responsable.DNI, responsable.Activo, responsable.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating responsable: %w", err)
	}

	return nil
}

func (r *ResponsableRepository) Delete(ctx context.Context, tx repository.Transaction, id int) error {
	query := `DELETE FROM responsabless WHERE id_responsable = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting responsable: %w", err)
	}

	return nil
}
