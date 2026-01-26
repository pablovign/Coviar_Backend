package domain

import "time"

// ============================================
// MODELOS DE BODEGA
// ============================================

type Bodega struct {
	ID                 int       `json:"id_bodega,omitempty"`
	RazonSocial        string    `json:"razon_social"`
	NombreFantasia     string    `json:"nombre_fantasia"`
	CUIT               string    `json:"cuit"`
	InvBod             *string   `json:"inv_bod,omitempty"`
	InvVin             *string   `json:"inv_vin,omitempty"`
	Calle              string    `json:"calle"`
	Numeracion         string    `json:"numeracion"`
	IDLocalidad        int       `json:"id_localidad"`
	Telefono           string    `json:"telefono"`
	EmailInstitucional string    `json:"email_institucional"`
	FechaRegistro      time.Time `json:"fecha_registro,omitempty"`
}

type BodegaRequest struct {
	RazonSocial        string  `json:"razon_social"`
	NombreFantasia     string  `json:"nombre_fantasia"`
	CUIT               string  `json:"cuit"`
	InvBod             *string `json:"inv_bod,omitempty"`
	InvVin             *string `json:"inv_vin,omitempty"`
	Calle              string  `json:"calle"`
	Numeracion         string  `json:"numeracion"`
	IDLocalidad        int     `json:"id_localidad"`
	Telefono           string  `json:"telefono"`
	EmailInstitucional string  `json:"email_institucional"`
}

type BodegaUpdateDTO struct {
	Telefono           string `json:"telefono"`
	EmailInstitucional string `json:"email_institucional"`
	NombreFantasia     string `json:"nombre_fantasia"`
}

// ============================================
// MODELOS DE CUENTA
// ============================================

type TipoCuenta string

const (
	TipoCuentaBodega           TipoCuenta = "BODEGA"
	TipoCuentaAdministradorApp TipoCuenta = "ADMINISTRADOR_APP"
)

type Cuenta struct {
	ID            int        `json:"id_cuenta,omitempty"`
	Tipo          TipoCuenta `json:"tipo"`
	IDBodega      *int       `json:"id_bodega,omitempty"`
	EmailLogin    string     `json:"email_login"`
	PasswordHash  string     `json:"-"`
	FechaRegistro time.Time  `json:"fecha_registro,omitempty"`
}

type CuentaRequest struct {
	EmailLogin string `json:"email_login"`
	Password   string `json:"password"`
}

// ============================================
// MODELOS DE RESPONSABLE
// ============================================

type Responsable struct {
	ID            int       `json:"id_responsable,omitempty"`
	IDBodega      int       `json:"id_bodega"`
	Nombre        string    `json:"nombre"`
	Apellido      string    `json:"apellido"`
	Cargo         string    `json:"cargo"`
	DNI           *string   `json:"dni,omitempty"`
	Activo        bool      `json:"activo"`
	FechaRegistro time.Time `json:"fecha_registro,omitempty"`
}

type ResponsableRequest struct {
	Nombre   string  `json:"nombre"`
	Apellido string  `json:"apellido"`
	Cargo    string  `json:"cargo"`
	DNI      *string `json:"dni,omitempty"`
}

// ============================================
// MODELOS DE USUARIO
// ============================================

type Usuario struct {
	IdUsuario     int        `json:"idUsuario,omitempty" db:"id_usuario"`
	Email         string     `json:"email" db:"email"`
	PasswordHash  string     `json:"-" db:"password_hash"`
	Nombre        string     `json:"nombre" db:"nombre"`
	Apellido      string     `json:"apellido" db:"apellido"`
	Rol           string     `json:"rol" db:"rol"`
	Activo        bool       `json:"activo" db:"activo"`
	FechaRegistro time.Time  `json:"fecha_registro" db:"fecha_registro"`
	UltimoAcceso  *time.Time `json:"ultimo_acceso" db:"ultimo_acceso"`
}

type UsuarioDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Rol      string `json:"rol"`
}

type UsuarioLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ToPublic limpia campos sensibles antes de enviar al cliente
func (u *Usuario) ToPublic() *Usuario {
	publicUser := *u
	publicUser.PasswordHash = ""
	return &publicUser
}

// ============================================
// MODELOS DE UBICACIÓN
// ============================================

type Provincia struct {
	ID     int    `json:"id_provincia,omitempty"`
	Nombre string `json:"nombre"`
}

type Departamento struct {
	ID          int    `json:"id_departamento,omitempty"`
	IDProvincia int    `json:"id_provincia"`
	Nombre      string `json:"nombre"`
}

type Localidad struct {
	ID             int    `json:"id_localidad,omitempty"`
	IDDepartamento int    `json:"id_departamento"`
	Nombre         string `json:"nombre"`
}

// ============================================
// DTOs DE REGISTRO
// ============================================

type RegistroRequest struct {
	Bodega      BodegaRequest      `json:"bodega"`
	Cuenta      CuentaRequest      `json:"cuenta"`
	Responsable ResponsableRequest `json:"responsable"`
}

type RegistroResponse struct {
	IDBodega      int    `json:"id_bodega"`
	IDCuenta      int    `json:"id_cuenta"`
	IDResponsable int    `json:"id_responsable"`
	Mensaje       string `json:"mensaje"`
}
