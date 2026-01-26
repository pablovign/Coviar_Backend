package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

// ConnectSupabase establece conexión con Supabase usando PostgreSQL
func ConnectSupabase(supabaseURL, supabaseKey, dbPassword string) (*DB, error) {
	// Extraer project ref de la URL
	var projectRef string
	if len(supabaseURL) > 8 {
		urlWithoutProtocol := supabaseURL[8:]
		if len(urlWithoutProtocol) > 0 {
			for i, c := range urlWithoutProtocol {
				if c == '.' {
					projectRef = urlWithoutProtocol[:i]
					break
				}
			}
		}
	}

	if projectRef == "" {
		return nil, fmt.Errorf("no se pudo extraer el project ref de la URL de Supabase")
	}

	// Configuración del Session Pooler de Supabase
	// Basado en: postgresql://postgres.jibisagabcbajwgliero:[PASSWORD]@aws-0-us-west-2.pooler.supabase.com:5432/postgres
	host := "aws-0-us-west-2.pooler.supabase.com"
	port := 5432
	user := fmt.Sprintf("postgres.%s", projectRef)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=require connect_timeout=15",
		host, port, user, dbPassword)

	fmt.Printf("🔍 Conectando a Supabase Session Pooler...\n")
	fmt.Printf("   Host: %s:%d\n", host, port)
	fmt.Printf("   Usuario: %s\n", user)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error abriendo conexión: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error conectando a Supabase: %w", err)
	}

	fmt.Printf("✅ Conexión exitosa a Supabase Session Pooler\n")
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
