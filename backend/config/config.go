package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	Environment string
	AppVersion  string
}

func LoadConfig() *Config {
	// Load .env if present (ignore error if not present in production/container)
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://kadai_user:kadai_password@localhost:5432/kadai_iq?sslmode=disable"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		Environment: env,
		AppVersion:  "1.0.0",
	}
}
