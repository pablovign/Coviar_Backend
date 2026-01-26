# 🚨 Instrucciones Finales - Conexión a Supabase

## Problema Actual

El host `db.jibisagabcbajwgliero.supabase.co` no existe en DNS. El código está intentando diferentes hosts automáticamente, pero ninguno funciona.

## ✅ Solución

Necesitas obtener la **Connection String** correcta desde tu Dashboard de Supabase.

### Pasos a Seguir:

1. **Abre tu Dashboard de Supabase**
   - Ve a: https://app.supabase.com/
   - Selecciona tu proyecto: `jibisagabcbajwgliero`

2. **Obtén la Connection String**
   - Ve a **Settings** (⚙️) en el menú lateral
   - Click en **Database**
   - Busca la sección **Connection string** o **Connection pooling**

3. **Copia el HOST correcto**

   Verás algo como esto:

   **Opción A - Session Pooler (Recomendado):**
   ```
   postgresql://postgres.jibisagabcbajwgliero:[YOUR-PASSWORD]@aws-0-us-east-1.pooler.supabase.com:6543/postgres
   ```

   **Opción B - Direct Connection:**
   ```
   postgresql://postgres:[YOUR-PASSWORD]@db.jibisagabcbajwgliero.supabase.co:5432/postgres
   ```

   **Opción C - Transaction Pooler:**
   ```
   postgresql://postgres:[YOUR-PASSWORD]@aws-0-us-east-1.pooler.supabase.com:6543/postgres?pgbouncer=true
   ```

4. **Extrae el HOST y PUERTO**

   Del ejemplo de Session Pooler:
   - **HOST**: `aws-0-us-east-1.pooler.supabase.com`
   - **PUERTO**: `6543`

5. **Actualiza tu archivo .env**

   Agrega estas dos líneas al archivo `.env`:

   ```env
   # Host y puerto obtenidos del dashboard
   SUPABASE_DB_HOST=aws-0-us-east-1.pooler.supabase.com
   SUPABASE_DB_PORT=6543
   ```

   Tu archivo `.env` completo debería verse así:

   ```env
   # Supabase Configuration
   SUPABASE_URL=https://jibisagabcbajwgliero.supabase.co
   SUPABASE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImppYmlzYWdhYmNiYWp3Z2xpZXJvIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjczOTUxOTAsImV4cCI6MjA4Mjk3MTE5MH0.DKlobWHvhLTtoaVDIiWiMsm3iw_pBKQliH6C_-gFe9k
   SUPABASE_DB_PASSWORD=e4LtpIroGJP9wZ8L
   # IMPORTANTE: Agrega estas dos líneas con los valores de tu dashboard
   SUPABASE_DB_HOST=el_host_que_copiaste
   SUPABASE_DB_PORT=6543

   # Server Configuration
   SERVER_HOST=0.0.0.0
   SERVER_PORT=8080

   # JWT Configuration
   JWT_SECRET=your_jwt_secret_key_here_change_in_production

   # Application
   APP_ENV=development
   ```

6. **Ejecuta el servidor**

   ```bash
   cd coviar-backend-final
   go run cmd/api/main.go
   ```

   Deberías ver:
   ```
   ✓ Configuración cargada
   🔍 Intentando conectar a: el_host_correcto:6543
   ✅ Conexión exitosa a el_host_correcto:6543
   ✓ Conexión a Supabase establecida
   ✓ Repositorios inicializados
   ✓ Servicios inicializados
   ✓ Handlers inicializados
   🚀 Servidor iniciando en http://0.0.0.0:8080
   ```

## 📸 Screenshots de dónde encontrar la información

En el Dashboard de Supabase verás:

```
Settings > Database

┌─────────────────────────────────────────────────────────────┐
│ Connection string                                            │
├─────────────────────────────────────────────────────────────┤
│ Session pooler (recommended for serverless/edge)            │
│ postgresql://postgres.xxx:[PASSWORD]@aws-0-xxx.pooler...    │
│                                                              │
│ Transaction pooler                                           │
│ postgresql://postgres.xxx:[PASSWORD]@aws-0-xxx.pooler...    │
│                                                              │
│ Direct connection                                            │
│ postgresql://postgres:[PASSWORD]@db.xxx.supabase.co:5432... │
└─────────────────────────────────────────────────────────────┘
```

## ❓ ¿Qué opción elegir?

- **Session Pooler** (Puerto 6543): ✅ RECOMENDADO para aplicaciones
  - Mejor rendimiento
  - Maneja muchas conexiones simultáneas

- **Direct Connection** (Puerto 5432): ⚠️ Solo si Session Pooler no funciona
  - Conexión directa a PostgreSQL
  - Puede no estar disponible en todos los planes

## 🔧 Si aún no funciona

Si después de configurar el host correcto sigue sin funcionar:

1. **Verifica que tu proyecto de Supabase esté activo**
   - No debe estar pausado

2. **Revisa las reglas de firewall**
   - En Settings > Database > Connection pooling
   - Verifica que tu IP esté permitida

3. **Prueba con psql directamente**
   ```bash
   psql "postgresql://postgres:e4LtpIroGJP9wZ8L@EL_HOST:PUERTO/postgres?sslmode=require"
   ```

4. **Contacta al soporte de Supabase**
   - Puede que el proyecto tenga restricciones especiales

## 📝 Notas Importantes

- ✅ El código ya está configurado para aceptar host personalizado
- ✅ El código intentará automáticamente diferentes configuraciones si no proporcionas host
- ✅ La contraseña ya está en tu .env
- ⚠️ NUNCA compartas tu contraseña de base de datos públicamente
- ⚠️ El host puede ser diferente dependiendo de tu región (us-east-1, us-west-1, eu-west-1, etc.)

## 🎯 Próximos Pasos

Una vez que te conectes exitosamente:

1. Ejecutar las migraciones SQL en Supabase (crear tablas)
2. Probar los endpoints de la API
3. ¡Empezar a desarrollar!

---

**¿Necesitas ayuda?** Comparte un screenshot de la sección "Database" > "Connection string" de tu dashboard (ocultando la contraseña) y te ayudaré a configurarlo.
