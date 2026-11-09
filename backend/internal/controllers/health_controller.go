package controllers

import (
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
	version string
}

func NewHealthController(version string) *HealthController {
	return &HealthController{
		version: version,
	}
}

func (h *HealthController) Check(c *fiber.Ctx) error {
	data := fiber.Map{
		"service":   "KadaiIQ API",
		"status":    "healthy",
		"version":   h.version,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse("Service is healthy", data))
}
