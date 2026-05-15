package main

import (
	"log"

	"tech-challenge-users/internal/adapter/config"
	"tech-challenge-users/internal/adapter/database"
	"tech-challenge-users/internal/adapter/database/migrations"
	"tech-challenge-users/internal/adapter/database/repository"
	httpAdapter "tech-challenge-users/internal/adapter/http"
	"tech-challenge-users/internal/adapter/http/handlers"
	"tech-challenge-users/internal/application/services"
	// future modules will import additional services here
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

	// Repositories
	personRepo := repository.NewPersonRepository(db)
	userRepo := repository.NewUserRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	transactor := repository.NewTransactor(db)

	// Services
	userService := services.NewUserService(transactor, personRepo, userRepo, employeeRepo)
	customerService := services.NewCustomerService(transactor, personRepo, customerRepo)

	// Bootstrap admin user (idempotent)
	userService.CreateAdminUser(services.AdminConfig{
		Email:    cfg.AdminEmail,
		Password: cfg.AdminPassword,
		Document: cfg.AdminDocument,
	})

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	internalUserHandler := handlers.NewInternalUserHandler(userService)
	customerHandler := handlers.NewCustomerHandler(customerService)

	// Router
	router := httpAdapter.NewRouter(userHandler, internalUserHandler, customerHandler)
	router.Setup()

	addr := ":" + cfg.HTTPPort
	log.Printf("starting tech-challenge-users on %s (env=%s)", addr, cfg.ENV)

	if err := router.Engine().Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
