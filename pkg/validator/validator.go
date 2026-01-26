package validator

import (
	"fmt"
	"regexp"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var msgs []string
	for _, err := range v {
		msgs = append(msgs, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(msgs, "; ")
}

var (
	cuitRegex     = regexp.MustCompile(`^[0-9]{11}$`)
	dniRegex      = regexp.MustCompile(`^[0-9]{7,8}$`)
	telefonoRegex = regexp.MustCompile(`^[0-9]+$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func ValidateCUIT(cuit string) error {
	if !cuitRegex.MatchString(cuit) {
		return fmt.Errorf("CUIT debe tener exactamente 11 dígitos numéricos")
	}
	return nil
}

func ValidateDNI(dni string) error {
	if dni == "" {
		return nil // DNI es opcional
	}
	if !dniRegex.MatchString(dni) {
		return fmt.Errorf("DNI debe tener 7 u 8 dígitos numéricos")
	}
	return nil
}

func ValidateTelefono(telefono string) error {
	if !telefonoRegex.MatchString(telefono) {
		return fmt.Errorf("teléfono solo debe contener dígitos")
	}
	return nil
}

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("formato de email inválido")
	}
	return nil
}

func ValidateNotEmpty(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s no puede estar vacío", fieldName)
	}
	return nil
}

func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("la contraseña debe tener al menos 8 caracteres")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("la contraseña debe tener al menos 6 caracteres")
	}
	return nil
}
