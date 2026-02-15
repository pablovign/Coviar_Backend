package service

import (
	"context"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
	"coviar_backend/pkg/validator"
)

type BodegaService struct {
	bodegaRepo repository.BodegaRepository
}

func NewBodegaService(bodegaRepo repository.BodegaRepository) *BodegaService {
	return &BodegaService{bodegaRepo: bodegaRepo}
}

func (s *BodegaService) GetByID(ctx context.Context, id int) (*domain.Bodega, error) {
	return s.bodegaRepo.FindByID(ctx, id)
}

func (s *BodegaService) Update(ctx context.Context, id int, dto *domain.BodegaUpdateDTO) error {
	if err := validator.ValidateTelefono(dto.Telefono); err != nil {
		return domain.ErrValidation
	}
	if err := validator.ValidateEmail(dto.EmailInstitucional); err != nil {
		return domain.ErrValidation
	}
	if err := validator.ValidateNotEmpty(dto.NombreFantasia, "nombre_fantasia"); err != nil {
		return domain.ErrValidation
	}

	// Normalizar campos de texto a mayúsculas
	dto.NombreFantasia = validator.NormalizarTexto(dto.NombreFantasia)

	bodega, err := s.bodegaRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	bodega.Telefono = dto.Telefono
	bodega.EmailInstitucional = dto.EmailInstitucional
	bodega.NombreFantasia = dto.NombreFantasia

	return s.bodegaRepo.Update(ctx, nil, bodega)
}
