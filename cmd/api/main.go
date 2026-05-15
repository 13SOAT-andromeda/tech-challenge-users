package main

import (
	"log"

	"tech-challenge-users/internal/adapter/config"
	"tech-challenge-users/internal/adapter/database"
	"tech-challenge-users/internal/adapter/database/migrations"
	httpAdapter "tech-challenge-users/internal/adapter/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db := database.Connect(cfg)

	if err := migrations.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	router := httpAdapter.NewRouter()
	router.Setup()

	addr := ":" + cfg.HTTPPort
	log.Printf("starting tech-challenge-users on %s (env=%s)", addr, cfg.ENV)

	if err := router.Engine().Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
