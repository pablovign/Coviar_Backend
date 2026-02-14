package service

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
	"coviar_backend/pkg/validator"
)

type RegistroService struct {
	bodegaRepo      repository.BodegaRepository
	cuentaRepo      repository.CuentaRepository
	responsableRepo repository.ResponsableRepository
	txManager       repository.TransactionManager
}

func NewRegistroService(
	bodegaRepo repository.BodegaRepository,
	cuentaRepo repository.CuentaRepository,
	responsableRepo repository.ResponsableRepository,
	txManager repository.TransactionManager,
) *RegistroService {
	return &RegistroService{
		bodegaRepo:      bodegaRepo,
		cuentaRepo:      cuentaRepo,
		responsableRepo: responsableRepo,
		txManager:       txManager,
	}
}

func (s *RegistroService) RegistrarBodega(ctx context.Context, req *domain.RegistroRequest) (*domain.RegistroResponse, error) {
	// Validar datos
	if err := s.validarRegistro(req); err != nil {
		return nil, err
	}

	// Normalizar datos a mayúsculas
	s.normalizarRegistroRequest(req)

	// Verificar duplicados
	if err := s.verificarDuplicados(ctx, req); err != nil {
		return nil, err
	}

	// Hash de contraseña
	passwordHash, err := hashPassword(req.Cuenta.Password)
	if err != nil {
		return nil, fmt.Errorf("error al procesar contraseña: %w", err)
	}

	// Iniciar transacción
	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error iniciando transacción: %w", err)
	}
	defer tx.Rollback()

	// Crear bodega
	bodega := &domain.Bodega{
		RazonSocial:        req.Bodega.RazonSocial,
		NombreFantasia:     req.Bodega.NombreFantasia,
		CUIT:               req.Bodega.CUIT,
		InvBod:             req.Bodega.InvBod,
		InvVin:             req.Bodega.InvVin,
		Calle:              req.Bodega.Calle,
		Numeracion:         req.Bodega.Numeracion,
		IDLocalidad:        req.Bodega.IDLocalidad,
		Telefono:           req.Bodega.Telefono,
		EmailInstitucional: req.Bodega.EmailInstitucional,
	}

	idBodega, err := s.bodegaRepo.Create(ctx, tx, bodega)
	if err != nil {
		return nil, fmt.Errorf("error creando bodega: %w", err)
	}

	// Crear cuenta
	cuenta := &domain.Cuenta{
		Tipo:         domain.TipoCuentaBodega,
		IDBodega:     &idBodega,
		EmailLogin:   req.Cuenta.EmailLogin,
		PasswordHash: passwordHash,
	}

	idCuenta, err := s.cuentaRepo.Create(ctx, tx, cuenta)
	if err != nil {
		return nil, fmt.Errorf("error creando cuenta: %w", err)
	}

	// Crear responsable
	responsable := &domain.Responsable{
		IDCuenta: idCuenta,
		Nombre:   req.Responsable.Nombre,
		Apellido: req.Responsable.Apellido,
		Cargo:    req.Responsable.Cargo,
		DNI:      "",
		Activo:   true,
	}
	if req.Responsable.DNI != nil {
		responsable.DNI = *req.Responsable.DNI
	}

	idResponsable, err := s.responsableRepo.Create(ctx, tx, responsable)
	if err != nil {
		return nil, fmt.Errorf("error creando responsable: %w", err)
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error confirmando transacción: %w", err)
	}

	return &domain.RegistroResponse{
		IDBodega:      idBodega,
		IDCuenta:      idCuenta,
		IDResponsable: idResponsable,
		Mensaje:       "Registro exitoso",
	}, nil
}

func (s *RegistroService) validarRegistro(req *domain.RegistroRequest) error {
	var errs validator.ValidationErrors

	// Validar bodega
	if err := validator.ValidateNotEmpty(req.Bodega.RazonSocial, "razon_social"); err != nil {
		errs = append(errs, validator.ValidationError{Field: "bodega.razon_social", Message: err.Error()})
	}
	if err := validator.ValidateNotEmpty(req.Bodega.NombreFantasia, "nombre_fantasia"); err != nil {
		errs = append(errs, validator.ValidationError{Field: "bodega.nombre_fantasia", Message: err.Error()})
	}
	if err := validator.ValidateCUIT(req.Bodega.CUIT); err != nil {
		errs = append(errs, validator.ValidationError{Field: "bodega.cuit", Message: err.Error()})
	}
	if err := validator.ValidateNotEmpty(req.Bodega.Calle, "calle"); err != nil {
		errs = append(errs, validator.ValidationError{Field: "bodega.calle", Message: err.Error()})
	}
	if err := validator.ValidateTelefono(req.Bodega.Telefono); err != nil {
		errs = append(errs, validator.ValidationError{Field: "bodega.telefono", Message: err.Error()})
	}
	if err := validator.ValidateEmail(req.Bodega.EmailInstitucional); err != nil {
		errs = append(errs, validator.ValidationError{Field: "bodega.email_institucional", Message: err.Error()})
	}
	if err := validator.ValidateInvCode(req.Bodega.InvBod, "inv_bod"); err != nil {
		errs = append(errs, validator.ValidationError{Field: "bodega.inv_bod", Message: err.Error()})
	}
	if err := validator.ValidateInvCode(req.Bodega.InvVin, "inv_vin"); err != nil {
		errs = append(errs, validator.ValidationError{Field: "bodega.inv_vin", Message: err.Error()})
	}

	// Validar cuenta
	if err := validator.ValidateEmail(req.Cuenta.EmailLogin); err != nil {
		errs = append(errs, validator.ValidationError{Field: "cuenta.email_login", Message: err.Error()})
	}
	if err := validator.ValidatePasswordStrength(req.Cuenta.Password); err != nil {
		errs = append(errs, validator.ValidationError{Field: "cuenta.password", Message: err.Error()})
	}

	// Validar responsable
	if err := validator.ValidateNotEmpty(req.Responsable.Nombre, "nombre"); err != nil {
		errs = append(errs, validator.ValidationError{Field: "responsable.nombre", Message: err.Error()})
	}
	if err := validator.ValidateNotEmpty(req.Responsable.Apellido, "apellido"); err != nil {
		errs = append(errs, validator.ValidationError{Field: "responsable.apellido", Message: err.Error()})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (s *RegistroService) verificarDuplicados(ctx context.Context, req *domain.RegistroRequest) error {
	// Verificar CUIT
	bodega, err := s.bodegaRepo.FindByCUIT(ctx, req.Bodega.CUIT)
	if err != nil {
		return fmt.Errorf("error al verificar CUIT: %w", err)
	}
	if bodega != nil {
		return domain.ErrCUITYaRegistrado
	}

	// Verificar email
	cuenta, err := s.cuentaRepo.FindByEmail(ctx, req.Cuenta.EmailLogin)
	if err != nil {
		return fmt.Errorf("error al verificar email: %w", err)
	}
	if cuenta != nil {
		return domain.ErrEmailYaRegistrado
	}

	return nil
}

// normalizarRegistroRequest convierte todos los campos de texto del request a mayúsculas.
// Los emails se convierten a minúsculas para estandarización (estándar de emails).
// Los campos numéricos como CUIT, DNI y teléfono no se modifican.
func (s *RegistroService) normalizarRegistroRequest(req *domain.RegistroRequest) {
	// ===== NORMALIZAR BODEGA =====
	req.Bodega.RazonSocial = validator.NormalizarTexto(req.Bodega.RazonSocial)
	req.Bodega.NombreFantasia = validator.NormalizarTexto(req.Bodega.NombreFantasia)
	req.Bodega.Calle = validator.NormalizarTexto(req.Bodega.Calle)
	req.Bodega.Numeracion = validator.NormalizarTexto(req.Bodega.Numeracion)

	// CUIT solo contiene números, no requiere conversión
	// req.Bodega.CUIT permanece igual

	// Teléfono solo contiene números, no requiere conversión
	// req.Bodega.Telefono permanece igual

	// Email institucional: se convierte a minúsculas (estándar de emails)
	req.Bodega.EmailInstitucional = strings.ToLower(strings.TrimSpace(req.Bodega.EmailInstitucional))

	// Códigos INV: si existen, convertir a mayúsculas
	req.Bodega.InvBod = validator.NormalizarPuntero(req.Bodega.InvBod)
	req.Bodega.InvVin = validator.NormalizarPuntero(req.Bodega.InvVin)

	// ===== NORMALIZAR CUENTA =====
	// Email de login: se convierte a minúsculas (estándar de emails)
	req.Cuenta.EmailLogin = strings.ToLower(strings.TrimSpace(req.Cuenta.EmailLogin))
	// La contraseña NO se modifica - es case-sensitive por seguridad

	// ===== NORMALIZAR RESPONSABLE =====
	req.Responsable.Nombre = validator.NormalizarTexto(req.Responsable.Nombre)
	req.Responsable.Apellido = validator.NormalizarTexto(req.Responsable.Apellido)
	req.Responsable.Cargo = validator.NormalizarTexto(req.Responsable.Cargo)
	// DNI solo contiene números, no requiere conversión
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
