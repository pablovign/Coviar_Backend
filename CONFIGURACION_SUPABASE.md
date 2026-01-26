# Configuración de Supabase para Coviar Backend Final

Este documento explica cómo configurar y conectar el backend de Coviar con Supabase.

## 📋 Cambios Realizados

Se ha migrado el proyecto de PostgreSQL local a Supabase. Los cambios incluyen:

1. **Configuración actualizada** (`pkg/config/config.go`)
   - Agregado soporte para variables de Supabase
   - Eliminada configuración de PostgreSQL local

2. **Conexión a base de datos** (`pkg/database/postgres.go`)
   - Nueva función `ConnectSupabase()` que se conecta a Supabase usando PostgreSQL
   - Usa el puerto 6543 (pooling con transaction mode) para mejor rendimiento

3. **Variables de entorno** (`.env`)
   - `SUPABASE_URL`: URL de tu proyecto en Supabase
   - `SUPABASE_KEY`: Clave anon/public key de Supabase
   - `SUPABASE_DB_PASSWORD`: Contraseña de la base de datos PostgreSQL

## 🔑 Obtener las Credenciales de Supabase

### 1. SUPABASE_URL y SUPABASE_KEY

Estas credenciales ya están configuradas en el archivo `.env`:

```env
SUPABASE_URL=https://qrixjgbzlxlyqjdfxnae.supabase.co
SUPABASE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 2. SUPABASE_DB_PASSWORD (REQUERIDO)

**Esta contraseña NO está configurada y debes obtenerla tú mismo:**

#### Pasos para obtener la contraseña:

1. Ve a tu dashboard de Supabase: https://app.supabase.com/
2. Selecciona tu proyecto `qrixjgbzlxlyqjdfxnae`
3. Ve a **Settings** (Configuración) en el menú lateral
4. Haz clic en **Database**
5. En la sección **Connection string**, verás la contraseña o podrás resetearla

#### Actualizar el archivo .env:

Una vez que tengas la contraseña, edita el archivo `.env` y reemplaza:

```env
SUPABASE_DB_PASSWORD=your_database_password_here
```

Por tu contraseña real:

```env
SUPABASE_DB_PASSWORD=tu_contraseña_real_aqui
```

⚠️ **IMPORTANTE**: Nunca compartas tu contraseña de base de datos públicamente.

## 🗄️ Crear las Tablas en Supabase

El proyecto incluye migraciones SQL en la carpeta `migrations/`. Debes ejecutar estas migraciones en tu base de datos de Supabase:

### Opción 1: Usando el SQL Editor de Supabase (Recomendado)

1. Ve a tu dashboard de Supabase
2. Abre el **SQL Editor** en el menú lateral
3. Copia y pega el contenido de `migrations/001_create_tables.sql`
4. Ejecuta el script
5. Luego copia y pega el contenido de `migrations/002_insert_data.sql`
6. Ejecuta el script

### Opción 2: Usando psql desde tu terminal

```bash
# Conectarse a Supabase
psql "postgresql://postgres:[TU_PASSWORD]@db.qrixjgbzlxlyqjdfxnae.supabase.co:5432/postgres"

# Dentro de psql, ejecutar:
\i migrations/001_create_tables.sql
\i migrations/002_insert_data.sql
```

## 🚀 Ejecutar el Proyecto

Una vez que hayas configurado la contraseña y creado las tablas:

```bash
cd coviar-backend-final

# Ejecutar en modo desarrollo
go run cmd/api/main.go

# O compilar y ejecutar
go build -o bin/coviar-api.exe cmd/api/main.go
bin/coviar-api.exe
```

Deberías ver:

```
✓ Configuración cargada
✓ Conexión a Supabase establecida
✓ Repositorios inicializados
✓ Servicios inicializados
✓ Handlers inicializados
🚀 Servidor iniciando en http://0.0.0.0:8080
📍 Entorno: development
🔗 Supabase URL: https://qrixjgbzlxlyqjdfxnae.supabase.co
```

## 🧪 Probar la Conexión

Prueba el endpoint de health check:

```bash
curl http://localhost:8080/health
```

Deberías recibir:

```json
{
  "status": "ok",
  "version": "2.0.0",
  "message": "Coviar Backend - Integrado y Funcional"
}
```

## 📊 Endpoints Disponibles

Todos los endpoints del proyecto están disponibles igual que antes:

- **Health Check**: `GET /health`
- **Usuarios**:
  - `GET /api/usuarios` - Listar usuarios
  - `POST /api/usuarios` - Crear usuario
  - `POST /api/usuarios/login` - Login
  - `GET /api/usuarios/{id}` - Obtener por ID
  - `DELETE /api/usuarios/{id}` - Eliminar
- **Ubicaciones**:
  - `GET /api/v1/provincias`
  - `GET /api/v1/departamentos`
  - `GET /api/v1/localidades`
- **Registro**:
  - `POST /api/registro` - Registrar bodega

## 🔧 Solución de Problemas

### Error: "SUPABASE_DB_PASSWORD son requeridas"

Verifica que hayas configurado la contraseña en el archivo `.env`.

### Error: "error connecting to Supabase database"

1. Verifica que la contraseña sea correcta
2. Verifica que tu IP esté en la lista blanca de Supabase (por defecto permite todas)
3. Verifica que el proyecto de Supabase esté activo

### Error: "relation does not exist"

Ejecuta las migraciones SQL en Supabase (ver sección "Crear las Tablas en Supabase").

## 🔐 Seguridad

1. **Nunca** compartas tu archivo `.env`
2. **Nunca** subas tu contraseña de base de datos a GitHub
3. El archivo `.gitignore` ya incluye `.env` para protegerte
4. Usa `.env.example` como plantilla para otros desarrolladores

## 📝 Diferencias con coviar-backend

El proyecto `coviar-backend` usa la librería `supabase-go` para peticiones HTTP/REST a la API de Supabase, mientras que `coviar-backend-final` se conecta directamente a la base de datos PostgreSQL de Supabase usando el driver `lib/pq`.

Ambos enfoques son válidos:
- **coviar-backend**: Más fácil de usar, abstrae la complejidad de SQL
- **coviar-backend-final**: Más control, permite queries SQL complejas, transacciones, etc.

## ✅ Lista de Verificación

Antes de ejecutar el proyecto:

- [ ] Obtener contraseña de base de datos de Supabase
- [ ] Actualizar `SUPABASE_DB_PASSWORD` en `.env`
- [ ] Ejecutar migraciones SQL en Supabase
- [ ] Verificar que las dependencias estén instaladas (`go mod download`)
- [ ] Ejecutar el proyecto y verificar la conexión

## 📚 Recursos

- [Documentación de Supabase](https://supabase.com/docs)
- [Supabase Database](https://supabase.com/docs/guides/database)
- [Connection Pooling](https://supabase.com/docs/guides/database/connecting-to-postgres#connection-pooler)
