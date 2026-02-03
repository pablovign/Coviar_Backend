-- Tabla para tokens de recuperación de contraseña
CREATE TABLE IF NOT EXISTS restaurar_contrasenas (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES cuentas(id_cuenta),
    token VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_restaurar_contrasenas_token ON restaurar_contrasenas(token);
CREATE INDEX IF NOT EXISTS idx_restaurar_contrasenas_user ON restaurar_contrasenas(user_id);
