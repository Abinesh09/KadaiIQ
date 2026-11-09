package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DatabaseURL       string
	Environment       string
	AppVersion        string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration
	DBTimeout         time.Duration
}

func LoadConfig() *Config {
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

	maxOpen := getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	maxIdle := getEnvAsInt("DB_MAX_IDLE_CONNS", 10)
	lifetimeMinutes := getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 30)
	idleMinutes := getEnvAsInt("DB_CONN_MAX_IDLE_MINUTES", 5)
	timeoutSeconds := getEnvAsInt("DB_TIMEOUT_SECONDS", 5)

	return &Config{
		Port:              port,
		DatabaseURL:       dbURL,
		Environment:       env,
		AppVersion:        "1.0.0",
		DBMaxOpenConns:    maxOpen,
		DBMaxIdleConns:    maxIdle,
		DBConnMaxLifetime: time.Duration(lifetimeMinutes) * time.Minute,
		DBConnMaxIdleTime: time.Duration(idleMinutes) * time.Minute,
		DBTimeout:         time.Duration(timeoutSeconds) * time.Second,
	}
}

func getEnvAsInt(name string, defaultValue int) int {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}
