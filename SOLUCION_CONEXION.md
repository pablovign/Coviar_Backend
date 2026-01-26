# Solución al Problema de Conexión a Supabase

## ⚠️ Problema Actual

El error `dial tcp: lookup db.jibisagabcbajwgliero.supabase.co: no such host` indica que el host de la base de datos no existe en DNS.

## 🔍 Por qué ocurre esto

Supabase ha cambiado su infraestructura y la forma de conectarse a PostgreSQL directamente. Los proyectos nuevos pueden:
1. No tener habilitado el acceso directo a PostgreSQL
2. Usar diferentes formatos de conexión según la región
3. Requerir el uso del pooler de conexiones

## ✅ Soluciones

### Opción 1: Obtener la Cadena de Conexión Correcta (RECOMENDADO)

1. Ve a tu dashboard de Supabase: https://app.supabase.com/
2. Selecciona tu proyecto `jibisagabcbajwgliero`
3. Ve a **Settings** > **Database**
4. Busca la sección **Connection string** o **Connection pooling**
5. Verás algo como:

```
postgresql://postgres.jibisagabcbajwgliero:[YOUR-PASSWORD]@aws-0-us-east-1.pooler.supabase.com:6543/postgres
```

6. Copia el **host** de esa cadena (por ejemplo: `aws-0-us-east-1.pooler.supabase.com`)

### Opción 2: Usar la API de Supabase en lugar de PostgreSQL Directo

Si no puedes conectarte directamente a PostgreSQL, puedes usar la biblioteca `supabase-go` que se conecta via API REST.

## 🛠️ Implementación de la Solución

### Para usar la cadena de conexión correcta:

Necesitamos modificar el código para aceptar un host personalizado. Agrega esta variable en tu `.env`:

```env
SUPABASE_DB_HOST=aws-0-us-east-1.pooler.supabase.com
```

O el host que hayas obtenido del dashboard.

## 📝 Pasos para Resolver

1. **Obtén la información de conexión:**
   ```bash
   # Ve a Supabase Dashboard
   # Settings > Database > Connection string
   ```

2. **Actualiza el .env con la información correcta:**
   ```env
   SUPABASE_DB_HOST=el_host_correcto_aqui
   SUPABASE_DB_PORT=6543
   ```

3. **Verifica que tienes acceso directo a PostgreSQL:**
   - Algunos planes de Supabase no permiten conexión directa
   - Puede que necesites habilitar "Direct connections" en el dashboard

## 🔄 Alternativa: Usar Supabase Client (API REST)

Si la conexión directa a PostgreSQL no funciona, podemos volver a usar la librería de Supabase que se conecta via HTTP:

```bash
cd coviar-backend-final
# Actualizar go.mod para usar supabase-go
```

Esta opción es más fácil pero tiene limitaciones:
- No soporta transacciones complejas
- No permite queries SQL directas
- Todas las operaciones deben ir a través de la API REST

## 📞 ¿Qué necesitas hacer AHORA?

1. Abre tu dashboard de Supabase
2. Ve a Settings > Database
3. Copia la **Connection string**
4. Envíame el HOST que aparece en esa cadena (sin la contraseña)
5. Yo actualizaré el código para usar ese host

### Ejemplo de lo que deberías ver:

```
Session pooler (recommended for serverless/edge functions)
postgresql://postgres.jibisagabcbajwgliero:[YOUR-PASSWORD]@aws-0-us-east-1.pooler.supabase.com:6543/postgres
                                                              ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑
                                                              Este es el HOST que necesito
```

O puede ser algo como:

```
Direct connection
postgresql://postgres:[YOUR-PASSWORD]@db.jibisagabcbajwgliero.supabase.co:5432/postgres
                                       ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑
                                       Este sería el HOST correcto
```

## 💡 Información Adicional

- **Puerto 6543**: Pooler de conexiones (recomendado para aplicaciones)
- **Puerto 5432**: Conexión directa (puede no estar disponible en todos los planes)
- **SSL**: Siempre es requerido (`sslmode=require`)

Dime qué ves en tu dashboard de Supabase y te ayudaré a configurarlo correctamente.
