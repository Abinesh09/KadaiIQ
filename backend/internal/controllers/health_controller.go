package controllers

import (
	"context"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/database"
	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
	version string
	db      *database.DBManager
}

func NewHealthController(version string, db *database.DBManager) *HealthController {
	return &HealthController{
		version: version,
		db:      db,
	}
}

// Check returns overall application status
func (h *HealthController) Check(c *fiber.Ctx) error {
	data := fiber.Map{
		"service":   "KadaiIQ API",
		"status":    "healthy",
		"version":   h.version,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse("Service is healthy", data))
}

// CheckDB tests the active PostgreSQL connection
func (h *HealthController) CheckDB(c *fiber.Ctx) error {
	if h.db == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(
			dto.ErrorResponse("Database connection is not configured", fiber.Map{
				"status": "unavailable",
			}),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	latency, err := h.db.Ping(ctx)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(
			dto.ErrorResponse("Database is unreachable", fiber.Map{
				"status": "unhealthy",
				"error":  err.Error(),
			}),
		)
	}

	stats := h.db.Stats()

	data := fiber.Map{
		"status":           "healthy",
		"latency_ms":       latency.Milliseconds(),
		"open_connections": stats.OpenConnections,
		"in_use":           stats.InUse,
		"idle":             stats.Idle,
		"timestamp":        time.Now().UTC().Format(time.RFC3339),
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse("Database is healthy", data))
}
