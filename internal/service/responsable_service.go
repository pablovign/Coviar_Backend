package service

import (
	"context"
	"time"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
	"coviar_backend/pkg/validator"
)

type ResponsableService struct {
	responsableRepo    repository.ResponsableRepository
	cuentaRepo         repository.CuentaRepository
	autoevaluacionRepo repository.AutoevaluacionRepository
}

func NewResponsableService(
	responsableRepo repository.ResponsableRepository,
	cuentaRepo repository.CuentaRepository,
	autoevaluacionRepo repository.AutoevaluacionRepository,
) *ResponsableService {
	return &ResponsableService{
		responsableRepo:    responsableRepo,
		cuentaRepo:         cuentaRepo,
		autoevaluacionRepo: autoevaluacionRepo,
	}
}

func (s *ResponsableService) GetByID(ctx context.Context, id int) (*domain.Responsable, error) {
	return s.responsableRepo.FindByID(ctx, id)
}

func (s *ResponsableService) GetByCuentaID(ctx context.Context, cuentaID int) ([]*domain.Responsable, error) {
	return s.responsableRepo.FindByCuentaID(ctx, cuentaID)
}

func (s *ResponsableService) Update(ctx context.Context, id int, dto *domain.ResponsableUpdateDTO) error {
	if err := validator.ValidateNotEmpty(dto.Nombre, "nombre"); err != nil {
		return domain.ErrValidation
	}
	if err := validator.ValidateNotEmpty(dto.Apellido, "apellido"); err != nil {
		return domain.ErrValidation
	}
	if err := validator.ValidateNotEmpty(dto.Cargo, "cargo"); err != nil {
		return domain.ErrValidation
	}
	if err := validator.ValidateDNI(dto.DNI); err != nil {
		return domain.ErrValidation
	}

	responsable, err := s.responsableRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	responsable.Nombre = dto.Nombre
	responsable.Apellido = dto.Apellido
	responsable.Cargo = dto.Cargo
	responsable.DNI = dto.DNI

	return s.responsableRepo.Update(ctx, nil, responsable)
}

func (s *ResponsableService) DarDeBaja(ctx context.Context, id int) error {
	responsable, err := s.responsableRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if !responsable.Activo {
		return domain.ErrResponsableYaDadoDeBaja
	}

	cuenta, err := s.cuentaRepo.FindByID(ctx, responsable.IDCuenta)
	if err != nil {
		return err
	}

	if cuenta.IDBodega != nil {
		hasPending, err := s.autoevaluacionRepo.HasPendingByBodega(ctx, *cuenta.IDBodega)
		if err != nil {
			return err
		}
		if hasPending {
			return domain.ErrAutoevaluacionesPendientes
		}
	}

	now := time.Now()
	responsable.Activo = false
	responsable.FechaBaja = &now

	return s.responsableRepo.Update(ctx, nil, responsable)
}

// CanAccess verifica si el usuario autenticado puede acceder al responsable
func (s *ResponsableService) CanAccess(ctx context.Context, responsableID int, userCuentaID int, userTipo string) (bool, error) {
	if userTipo == string(domain.TipoCuentaAdministradorApp) {
		return true, nil
	}

	responsable, err := s.responsableRepo.FindByID(ctx, responsableID)
	if err != nil {
		return false, err
	}

	return responsable.IDCuenta == userCuentaID, nil
}
