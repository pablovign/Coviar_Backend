package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Server   ServerConfig
	Supabase SupabaseConfig
	JWT      JWTConfig
	App      AppConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type SupabaseConfig struct {
	URL        string
	Key        string
	DBPassword string
}

type JWTConfig struct {
	Secret string
}

type AppConfig struct {
	Environment string
}

// Load carga las variables de entorno desde .env
func Load() (*Config, error) {
	// Cargar .env (opcional)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No se encontró .env, usando variables del sistema")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Supabase: SupabaseConfig{
			URL:        os.Getenv("SUPABASE_URL"),
			Key:        os.Getenv("SUPABASE_KEY"),
			DBPassword: os.Getenv("SUPABASE_DB_PASSWORD"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your_jwt_secret_key_here"),
		},
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
		},
	}

	// Validar variables críticas de Supabase
	if cfg.Supabase.URL == "" || cfg.Supabase.Key == "" || cfg.Supabase.DBPassword == "" {
		return nil, fmt.Errorf("SUPABASE_URL, SUPABASE_KEY y SUPABASE_DB_PASSWORD son requeridas")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
