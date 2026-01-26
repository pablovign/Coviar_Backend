-- ============================================
-- CREAR TABLA USUARIO
-- ============================================
-- Esta tabla es necesaria para el backend pero no existe en el esquema actual

CREATE TABLE IF NOT EXISTS public.usuarios (
    id_usuario SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    rol VARCHAR(50) NOT NULL DEFAULT 'bodega',
    activo BOOLEAN DEFAULT TRUE,
    fecha_registro TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índices para mejorar el rendimiento
CREATE INDEX IF NOT EXISTS idx_usuarios_email ON public.usuarios(email);
CREATE INDEX IF NOT EXISTS idx_usuarios_activo ON public.usuarios(activo);

-- Comentarios
COMMENT ON TABLE public.usuarios IS 'Tabla de usuarios del sistema';
COMMENT ON COLUMN public.usuarios.rol IS 'Rol del usuario: bodega, admin, etc.';
