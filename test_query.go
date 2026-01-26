package main

import (
	"context"
	"coviar_backend/pkg/config"
	"coviar_backend/pkg/database"
	"fmt"
	"log"
)

func testQuery() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.ConnectSupabase(cfg.Supabase.URL, cfg.Supabase.Key, cfg.Supabase.DBPassword)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.QueryContext(context.Background(), "SELECT id_provincia, nombre FROM provincias LIMIT 5")
	if err != nil {
		log.Fatalf("Error query: %v", err)
	}
	defer rows.Close()

	fmt.Println("Provincias:")
	for rows.Next() {
		var id int
		var nombre string
		rows.Scan(&id, &nombre)
		fmt.Printf("  %d: %s\n", id, nombre)
	}
}
