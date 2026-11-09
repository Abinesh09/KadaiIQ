package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBManager struct {
	DB  *gorm.DB
	sql *sql.DB
	cfg *config.Config
}

func NewDatabase(cfg *config.Config) (*DBManager, error) {
	// Configure GORM Logger based on environment
	logLevel := logger.Warn
	if cfg.Environment == "development" {
		logLevel = logger.Info
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n[GORM] ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open connection with GORM
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pooling
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	// Verify connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DBTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database during initialization: %w", err)
	}

	log.Printf("Connected to PostgreSQL successfully (MaxOpen: %d, MaxIdle: %d)", cfg.DBMaxOpenConns, cfg.DBMaxIdleConns)

	return &DBManager{
		DB:  db,
		sql: sqlDB,
		cfg: cfg,
	}, nil
}

// NewDBManagerFromDB creates a DBManager wrapping an existing *gorm.DB instance
func NewDBManagerFromDB(db *gorm.DB, cfg *config.Config) (*DBManager, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return &DBManager{
		DB:  db,
		sql: sqlDB,
		cfg: cfg,
	}, nil
}

// Ping checks if the database is reachable and returns the latency
func (m *DBManager) Ping(ctx context.Context) (time.Duration, error) {
	if m == nil || m.sql == nil {
		return 0, fmt.Errorf("database connection is not initialized")
	}

	start := time.Now()
	if err := m.sql.PingContext(ctx); err != nil {
		return 0, err
	}
	return time.Since(start), nil
}

// Stats returns the connection pool statistics
func (m *DBManager) Stats() sql.DBStats {
	if m == nil || m.sql == nil {
		return sql.DBStats{}
	}
	return m.sql.Stats()
}

// Close gracefully shuts down the database connection pool
func (m *DBManager) Close() error {
	if m != nil && m.sql != nil {
		log.Println("Closing database connection pool...")
		return m.sql.Close()
	}
	return nil
}
