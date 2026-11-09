package routes

import (
	"github.com/Abinesh09/KadaiIQ/backend/config"
	"github.com/Abinesh09/KadaiIQ/backend/database"
	"github.com/Abinesh09/KadaiIQ/backend/internal/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, dbManager *database.DBManager) {
	// Global Middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:7070,http://127.0.0.1:7070,http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Tenant-ID",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// API v1 group
	api := app.Group("/api/v1")

	// Health endpoints
	healthController := controllers.NewHealthController(cfg.AppVersion, dbManager)
	api.Get("/health", healthController.Check)
	api.Get("/health/db", healthController.CheckDB)
}
