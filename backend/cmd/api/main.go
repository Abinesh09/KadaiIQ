package main

import (
	"log"

	"github.com/Abinesh09/KadaiIQ/backend/config"
	"github.com/Abinesh09/KadaiIQ/backend/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()

	app := fiber.New(fiber.Config{
		AppName:      "KadaiIQ Retail SaaS API",
		ServerHeader: "Fiber",
	})

	// Setup application routes
	routes.SetupRoutes(app, cfg)

	log.Printf("Starting KadaiIQ API server on port %s in %s mode...", cfg.Port, cfg.Environment)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
