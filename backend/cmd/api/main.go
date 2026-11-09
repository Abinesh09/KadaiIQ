package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Abinesh09/KadaiIQ/backend/config"
	"github.com/Abinesh09/KadaiIQ/backend/database"
	"github.com/Abinesh09/KadaiIQ/backend/internal/routes"
	"github.com/Abinesh09/KadaiIQ/backend/migrations"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize Database Connection
	var dbManager *database.DBManager
	var err error
	dbManager, err = database.NewDatabase(cfg)
	if err != nil {
		log.Printf("[Warning] Database connection failed: %v", err)
		log.Printf("[Warning] Running server without active database connection. DB health endpoint will report unavailable.")
	} else {
		// Run pending versioned migrations
		if err := migrations.RunMigrations(dbManager.DB); err != nil {
			log.Fatalf("[Error] Failed to execute database migrations: %v", err)
		}
		defer func() {
			if err := dbManager.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			}
		}()
	}

	app := fiber.New(fiber.Config{
		AppName:      "KadaiIQ Retail SaaS API",
		ServerHeader: "Fiber",
	})

	// Setup application routes
	routes.SetupRoutes(app, cfg, dbManager)

	// Graceful shutdown channel
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting KadaiIQ API server on port %s in %s mode...", cfg.Port, cfg.Environment)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Printf("Server listen error: %v", err)
		}
	}()

	<-shutdownChan
	log.Println("Shutting down server gracefully...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down Fiber: %v", err)
	}
	log.Println("Server exited successfully.")
}
