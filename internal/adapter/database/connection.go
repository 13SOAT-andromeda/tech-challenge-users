package database

import (
	"errors"
	"fmt"
	"log"

	"tech-challenge-users/internal/adapter/config"

	gormtrace "github.com/DataDog/dd-trace-go/contrib/gorm.io/gorm.v1/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode, cfg.DBTimezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	errCheck := gormtrace.WithErrorCheck(func(err error) bool {
		return !errors.Is(err, gorm.ErrRecordNotFound)
	})
	if err := db.Use(gormtrace.NewTracePlugin(errCheck)); err != nil {
		log.Fatalf("failed to register gorm trace plugin: %v", err)
	}

	return db
}
