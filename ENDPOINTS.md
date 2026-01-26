# Endpoints API - Coviar Backend

## Base URL
```
http://localhost:8080
```

## Endpoints Implementados

### 1. Registro
**POST** `/api/registro`

Crea una bodega, cuenta y responsable en una sola transacción.

**Body:**
```json
{
  "bodega": {
    "razon_social": "Bodega Ejemplo SA",
    "nombre_fantasia": "Bodega Ejemplo",
    "cuit": "20-12345678-9",
    "inv_bod": "INV001",
    "inv_vin": "VIN001",
    "calle": "Calle Principal",
    "numeracion": "123",
    "id_localidad": 1,
    "telefono": "+54 261 1234567",
    "email_institucional": "contacto@bodega.com"
  },
  "cuenta": {
    "email_login": "admin@bodega.com",
    "password": "Password123!"
  },
  "responsable": {
    "nombre": "Juan",
    "apellido": "Pérez",
    "cargo": "Gerente",
    "dni": "12345678"
  }
}
```

**Response 201:**
```json
{
  "id_bodega": 1,
  "id_cuenta": 1,
  "id_responsable": 1,
  "mensaje": "Registro exitoso"
}
```

---

### 2. Login
**POST** `/api/login`

Autentica un usuario con email y contraseña.

**Body:**
```json
{
  "email_login": "admin@bodega.com",
  "password": "Password123!"
}
```

**Response 200:**
```json
{
  "id_cuenta": 1,
  "tipo": "BODEGA",
  "email_login": "admin@bodega.com",
  "fecha_registro": "2024-01-15T10:30:00Z",
  "bodega": {
    "id_bodega": 1,
    "razon_social": "Bodega Ejemplo SA",
    "nombre_fantasia": "Bodega Ejemplo",
    "cuit": "20-12345678-9",
    "telefono": "+54 261 1234567",
    "email_institucional": "contacto@bodega.com"
  }
}
```

---

### 3. Obtener Cuenta por ID
**GET** `/api/cuentas/{id}`

Obtiene una cuenta con sus datos de bodega asociada.

**Response 200:**
```json
{
  "id_cuenta": 1,
  "tipo": "BODEGA",
  "email_login": "admin@bodega.com",
  "fecha_registro": "2024-01-15T10:30:00Z",
  "bodega": {
    "id_bodega": 1,
    "razon_social": "Bodega Ejemplo SA",
    "nombre_fantasia": "Bodega Ejemplo"
  }
}
```

---

### 4. Modificar Contraseña
**PUT** `/api/cuentas/{id}`

Actualiza la contraseña de una cuenta.

**Body:**
```json
{
  "password": "NewPassword123!"
}
```

**Response 200:**
```json
{
  "mensaje": "Contraseña actualizada"
}
```

---

### 5. Obtener Bodega por ID
**GET** `/api/bodegas/{id}`

Obtiene los datos completos de una bodega.

**Response 200:**
```json
{
  "id_bodega": 1,
  "razon_social": "Bodega Ejemplo SA",
  "nombre_fantasia": "Bodega Ejemplo",
  "cuit": "20-12345678-9",
  "inv_bod": "INV001",
  "inv_vin": "VIN001",
  "calle": "Calle Principal",
  "numeracion": "123",
  "id_localidad": 1,
  "telefono": "+54 261 1234567",
  "email_institucional": "contacto@bodega.com",
  "fecha_registro": "2024-01-15T10:30:00Z"
}
```

---

### 6. Modificar Datos de Bodega
**PUT** `/api/bodegas/{id}`

Actualiza teléfono, email institucional y nombre de fantasía.

**Body:**
```json
{
  "telefono": "+54 261 9876543",
  "email_institucional": "nuevo@bodega.com",
  "nombre_fantasia": "Nuevo Nombre"
}
```

**Response 200:**
```json
{
  "mensaje": "Bodega actualizada"
}
```

---

### 7. Listar Provincias
**GET** `/api/provincias`

Obtiene todas las provincias.

**Response 200:**
```json
[
  {
    "id": 1,
    "nombre": "Mendoza"
  },
  {
    "id": 2,
    "nombre": "San Juan"
  }
]
```

---

### 8. Listar Departamentos
**GET** `/api/departamentos?provincia={id}`

Obtiene departamentos de una provincia específica.

**Query Params:**
- `provincia` (opcional): ID de la provincia

**Response 200:**
```json
[
  {
    "id": 1,
    "id_provincia": 1,
    "nombre": "Luján de Cuyo"
  },
  {
    "id": 2,
    "id_provincia": 1,
    "nombre": "Maipú"
  }
]
```

---

### 9. Listar Localidades
**GET** `/api/localidades?departamento={id}`

Obtiene localidades de un departamento específico.

**Query Params:**
- `departamento` (opcional): ID del departamento

**Response 200:**
```json
[
  {
    "id": 1,
    "id_departamento": 1,
    "nombre": "Chacras de Coria"
  },
  {
    "id": 2,
    "id_departamento": 1,
    "nombre": "Vistalba"
  }
]
```

---

## Códigos de Error

- `400` - Bad Request (datos inválidos)
- `404` - Not Found (recurso no encontrado)
- `500` - Internal Server Error

**Ejemplo de error:**
```json
{
  "error": "credenciales inválidas"
}
```

## Notas

1. La conexión con Supabase usa **Session Pooler** en el puerto 5432
2. Todas las contraseñas se hashean con bcrypt
3. El registro crea bodega, cuenta y responsable en una transacción atómica
4. Los endpoints de ubicación soportan filtrado por query parameters
