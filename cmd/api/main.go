package main

import (
	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/DataDog/dd-trace-go/v2/profiler"
	"go.uber.org/zap"

	"tech-challenge-users/internal/adapter/config"
	"tech-challenge-users/internal/adapter/database"
	"tech-challenge-users/internal/adapter/database/migrations"
	"tech-challenge-users/internal/adapter/database/repository"
	httpAdapter "tech-challenge-users/internal/adapter/http"
	"tech-challenge-users/internal/adapter/http/handlers"
	"tech-challenge-users/internal/application/services"
	"tech-challenge-users/internal/application/usecases"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer func() {
		tracer.Stop()
		profiler.Stop()
		_ = logger.Sync()
	}()

	sugar := logger.Sugar()

	cfg, err := config.Load()
	if err != nil {
		sugar.Fatalf("failed to load config: %v", err)
	}

	if err = profiler.Start(
		profiler.WithEnv(cfg.ENV),
		profiler.WithService(cfg.DDService),
		profiler.WithVersion(cfg.APIVersion),
		profiler.WithTags("layer:api"),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	); err != nil {
		sugar.Warnf("datadog profiler unavailable: %v", err)
	}

	if err = tracer.Start(
		tracer.WithEnv(cfg.ENV),
		tracer.WithService(cfg.DDService),
		tracer.WithServiceVersion(cfg.APIVersion),
	); err != nil {
		sugar.Fatalf("failed to start datadog tracer: %v", err)
	}

	if !cfg.DogStatsD.Disabled && cfg.DogStatsD.Addr != "" {
		statsdClient, errStatsd := statsd.New(cfg.DogStatsD.Addr,
			statsd.WithNamespace("tech_challenge_users."),
			statsd.WithTags([]string{
				"env:" + cfg.ENV,
				"service:" + cfg.DDService,
				"version:" + cfg.APIVersion,
			}),
		)
		if errStatsd != nil {
			sugar.Warnf("dogstatsd unavailable, metrics disabled: %v", errStatsd)
		} else {
			defer statsdClient.Close()
			sugar.Infof("dogstatsd connected: %s", cfg.DogStatsD.Addr)
		}
	}

	db := database.Connect(cfg)

	if err := migrations.RunMigrations(db); err != nil {
		sugar.Fatalf("failed to run migrations: %v", err)
	}

	// Repositories
	personRepo := repository.NewPersonRepository(db)
	userRepo := repository.NewUserRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	vehicleRepo := repository.NewVehicleRepository(db)
	cvRepo := repository.NewCustomerVehicleRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	transactor := repository.NewTransactor(db)

	// Services / Use cases
	userService := services.NewUserService(transactor, personRepo, userRepo, employeeRepo)
	customerService := services.NewCustomerService(transactor, personRepo, customerRepo)
	vehicleService := services.NewVehicleService(vehicleRepo)
	customerUseCase := usecases.NewCustomerUseCase(customerRepo, vehicleRepo, cvRepo)
	companyService := services.NewCompanyService(companyRepo)
	employeeService := services.NewEmployeeService(employeeRepo)

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
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)
	customerVehicleHandler := handlers.NewCustomerVehicleHandler(customerUseCase)
	companyHandler := handlers.NewCompanyHandler(companyService)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	// Router
	router := httpAdapter.NewRouter(
		logger,
		cfg.DDService,
		userHandler, internalUserHandler, customerHandler,
		vehicleHandler, customerVehicleHandler, companyHandler, employeeHandler,
	)
	router.Setup()

	addr := ":" + cfg.HTTPPort
	sugar.Infof("starting %s on %s (env=%s)", cfg.DDService, addr, cfg.ENV)

	if err := router.Engine().Run(addr); err != nil {
		sugar.Fatalf("server error: %v", err)
	}
}
