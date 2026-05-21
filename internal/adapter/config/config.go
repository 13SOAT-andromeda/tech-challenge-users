package config

import (
	"fmt"
	"net"
	"os"
)

type DogStatsDConfig struct {
	Addr     string
	Disabled bool
}

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBTimezone string

	HTTPPort           string
	HTTPAllowedOrigins string

	AdminEmail    string
	AdminPassword string
	AdminDocument string

	DDService  string
	DDSite     string
	DDAPIKey   string
	APIVersion string
	ENV        string

	DogStatsD *DogStatsDConfig
}

func Load() (*Config, error) {
	required := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
	}
	for key, val := range required {
		if val == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	// DogStatsD na porta 8125 do mesmo host do Agent. Sem host, métricas desligadas.
	agentHost := os.Getenv("DD_AGENT_HOST")
	dogstatsdAddr := ""
	dogstatsdDisabled := true
	if agentHost != "" {
		dogstatsdAddr = net.JoinHostPort(agentHost, "8125")
		dogstatsdDisabled = false
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     getEnvOrDefault("DB_PORT", "5432"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  getEnvOrDefault("DB_SSLMODE", "require"),
		DBTimezone: getEnvOrDefault("DB_TIMEZONE", "America/Sao_Paulo"),

		HTTPPort:           getEnvOrDefault("HTTP_PORT", "8081"),
		HTTPAllowedOrigins: getEnvOrDefault("HTTP_ALLOWED_ORIGINS", "*"),

		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		AdminDocument: os.Getenv("ADMIN_DOCUMENT"),

		DDService:  getEnvOrDefault("DD_SERVICE", "tech-challenge-users"),
		DDSite:     getEnvOrDefault("DD_SITE", "datadoghq.com"),
		DDAPIKey:   os.Getenv("DD_API_KEY"),
		APIVersion: getEnvOrDefault("API_VERSION", "v1"),
		ENV:        getEnvOrDefault("ENV", "development"),

		DogStatsD: &DogStatsDConfig{
			Addr:     dogstatsdAddr,
			Disabled: dogstatsdDisabled,
		},
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
