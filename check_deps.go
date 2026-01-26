package main

import (
	"coviar_backend/pkg/config"
	"coviar_backend/pkg/database"
	"fmt"
)

func checkDeps() {
	cfg, _ := config.Load()
	db, _ := database.ConnectSupabase(cfg.Supabase.URL, cfg.Supabase.Key, cfg.Supabase.DBPassword)
	defer db.Close()

	var count int
	db.QueryRow("SELECT COUNT(*) FROM departamentos").Scan(&count)
	fmt.Printf("Departamentos: %d\n", count)

	db.QueryRow("SELECT COUNT(*) FROM localidades").Scan(&count)
	fmt.Printf("Localidades: %d\n", count)
}
