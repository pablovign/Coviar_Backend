-- ============================================
-- DATOS DE EJEMPLO - UBICACIONES
-- ============================================
-- Este script inserta datos de ejemplo para las provincias, departamentos y localidades

-- Insertar Provincias de ejemplo (Argentina)
INSERT INTO provincia (nombre) VALUES 
('Buenos Aires'),
('Córdoba'),
('Mendoza'),
('San Juan'),
('La Rioja'),
('Catamarca'),
('Tucumán'),
('Misiones'),
('Corrientes'),
('Salta'),
('Jujuy'),
('Formosa'),
('Chaco'),
('Santiago del Estero'),
('Santa Fe'),
('Entre Ríos'),
('Neuquén'),
('Río Negro'),
('Chubut'),
('Santa Cruz'),
('Tierra del Fuego')
ON CONFLICT (nombre) DO NOTHING;

-- Insertar Departamentos de ejemplo para Buenos Aires
INSERT INTO departamento (id_provincia, nombre) 
SELECT id, 'Acoyte' FROM provincia WHERE nombre = 'Buenos Aires'
ON CONFLICT DO NOTHING;

INSERT INTO departamento (id_provincia, nombre) 
SELECT id, 'Almirante Brown' FROM provincia WHERE nombre = 'Buenos Aires'
ON CONFLICT DO NOTHING;

INSERT INTO departamento (id_provincia, nombre) 
SELECT id, 'Berazategui' FROM provincia WHERE nombre = 'Buenos Aires'
ON CONFLICT DO NOTHING;

-- Insertar Localidades de ejemplo
INSERT INTO localidad (id_departamento, nombre) 
SELECT id, 'La Matanza' FROM departamento WHERE nombre = 'Acoyte' AND id_provincia = (SELECT id FROM provincia WHERE nombre = 'Buenos Aires')
ON CONFLICT DO NOTHING;

INSERT INTO localidad (id_departamento, nombre) 
SELECT id, 'Lomas de Zamora' FROM departamento WHERE nombre = 'Almirante Brown' AND id_provincia = (SELECT id FROM provincia WHERE nombre = 'Buenos Aires')
ON CONFLICT DO NOTHING;

INSERT INTO localidad (id_departamento, nombre) 
SELECT id, 'Berazategui' FROM departamento WHERE nombre = 'Berazategui' AND id_provincia = (SELECT id FROM provincia WHERE nombre = 'Buenos Aires')
ON CONFLICT DO NOTHING;

-- Más departamentos para Córdoba
INSERT INTO departamento (id_provincia, nombre) 
SELECT id, 'Colón' FROM provincia WHERE nombre = 'Córdoba'
ON CONFLICT DO NOTHING;

INSERT INTO departamento (id_provincia, nombre) 
SELECT id, 'Capital' FROM provincia WHERE nombre = 'Córdoba'
ON CONFLICT DO NOTHING;

-- Localidades para Córdoba
INSERT INTO localidad (id_departamento, nombre) 
SELECT id, 'Córdoba' FROM departamento WHERE nombre = 'Capital' AND id_provincia = (SELECT id FROM provincia WHERE nombre = 'Córdoba')
ON CONFLICT DO NOTHING;

-- Departamentos y Localidades para Mendoza
INSERT INTO departamento (id_provincia, nombre) 
SELECT id, 'Capital' FROM provincia WHERE nombre = 'Mendoza'
ON CONFLICT DO NOTHING;

INSERT INTO localidad (id_departamento, nombre) 
SELECT id, 'Mendoza' FROM departamento WHERE nombre = 'Capital' AND id_provincia = (SELECT id FROM provincia WHERE nombre = 'Mendoza')
ON CONFLICT DO NOTHING;

-- ============================================
-- FIN DEL SCRIPT DE DATOS DE EJEMPLO
-- ============================================
