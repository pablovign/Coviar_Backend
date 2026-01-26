-- ============================================
-- SCRIPT DE CREACIÓN DE TABLAS - COVIAR BACKEND
-- ============================================
-- Este script crea todas las tablas necesarias para el backend integrado

-- Crear extensión para UUID (opcional)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- TABLA: PROVINCIA
-- ============================================
CREATE TABLE IF NOT EXISTS provincia (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- TABLA: DEPARTAMENTO
-- ============================================
CREATE TABLE IF NOT EXISTS departamento (
    id SERIAL PRIMARY KEY,
    id_provincia INTEGER NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_provincia) REFERENCES provincia(id) ON DELETE CASCADE,
    UNIQUE(id_provincia, nombre)
);

-- ============================================
-- TABLA: LOCALIDAD
-- ============================================
CREATE TABLE IF NOT EXISTS localidad (
    id SERIAL PRIMARY KEY,
    id_departamento INTEGER NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_departamento) REFERENCES departamento(id) ON DELETE CASCADE,
    UNIQUE(id_departamento, nombre)
);

-- ============================================
-- TABLA: BODEGA
-- ============================================
CREATE TABLE IF NOT EXISTS bodega (
    id_bodega SERIAL PRIMARY KEY,
    razon_social VARCHAR(200) NOT NULL,
    nombre_fantasia VARCHAR(200) NOT NULL,
    cuit VARCHAR(11) NOT NULL UNIQUE,
    inv_bod VARCHAR(50),
    inv_vin VARCHAR(50),
    calle VARCHAR(200) NOT NULL,
    numeracion VARCHAR(20) NOT NULL,
    id_localidad INTEGER NOT NULL,
    telefono VARCHAR(20) NOT NULL,
    email_institucional VARCHAR(100) NOT NULL UNIQUE,
    fecha_registro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_localidad) REFERENCES localidad(id) ON DELETE SET NULL
);

-- ============================================
-- TABLA: CUENTA
-- ============================================
CREATE TABLE IF NOT EXISTS cuenta (
    id_cuenta SERIAL PRIMARY KEY,
    tipo VARCHAR(50) NOT NULL DEFAULT 'BODEGA',
    id_bodega INTEGER,
    email_login VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    fecha_registro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_bodega) REFERENCES bodega(id_bodega) ON DELETE SET NULL
);

-- ============================================
-- TABLA: RESPONSABLE
-- ============================================
CREATE TABLE IF NOT EXISTS responsable (
    id_responsable SERIAL PRIMARY KEY,
    id_bodega INTEGER NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    cargo VARCHAR(100) NOT NULL,
    dni VARCHAR(12),
    activo BOOLEAN DEFAULT TRUE,
    fecha_registro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_bodega) REFERENCES bodega(id_bodega) ON DELETE CASCADE
);

-- ============================================
-- TABLA: USUARIO
-- ============================================
CREATE TABLE IF NOT EXISTS usuario (
    id_usuario SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    rol VARCHAR(50) NOT NULL DEFAULT 'bodega',
    activo BOOLEAN DEFAULT TRUE,
    fecha_registro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ultimo_acceso TIMESTAMP
);

-- ============================================
-- ÍNDICES PARA MEJOR RENDIMIENTO
-- ============================================
CREATE INDEX idx_bodega_cuit ON bodega(cuit);
CREATE INDEX idx_bodega_email ON bodega(email_institucional);
CREATE INDEX idx_cuenta_email ON cuenta(email_login);
CREATE INDEX idx_cuenta_bodega ON cuenta(id_bodega);
CREATE INDEX idx_responsable_bodega ON responsable(id_bodega);
CREATE INDEX idx_usuario_email ON usuario(email);
CREATE INDEX idx_usuario_rol ON usuario(rol);
CREATE INDEX idx_usuario_activo ON usuario(activo);
CREATE INDEX idx_departamento_provincia ON departamento(id_provincia);
CREATE INDEX idx_localidad_departamento ON localidad(id_departamento);
CREATE INDEX idx_bodega_localidad ON bodega(id_localidad);

-- ============================================
-- COMENTARIOS DE TABLAS
-- ============================================
COMMENT ON TABLE provincia IS 'Provincias/Estados del país';
COMMENT ON TABLE departamento IS 'Departamentos dentro de provincias';
COMMENT ON TABLE localidad IS 'Localidades dentro de departamentos';
COMMENT ON TABLE bodega IS 'Bodegas registradas en el sistema';
COMMENT ON TABLE cuenta IS 'Cuentas de acceso al sistema';
COMMENT ON TABLE responsable IS 'Responsables de las bodegas';
COMMENT ON TABLE usuario IS 'Usuarios del sistema con diferentes roles';
