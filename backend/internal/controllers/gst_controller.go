package controllers

import (
	"net/http"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/middleware"
	"github.com/Abinesh09/KadaiIQ/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type GSTController struct {
	gstService *services.GSTService
}

func NewGSTController(gstService *services.GSTService) *GSTController {
	return &GSTController{gstService: gstService}
}

func (ctrl *GSTController) GetSettings(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	settings, err := ctrl.gstService.GetSettings(tc.TenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch GST settings", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("GST settings retrieved", settings))
}

func (ctrl *GSTController) UpdateSettings(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	var req dto.UpdateGSTSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid request body", err.Error()))
	}

	settings, err := ctrl.gstService.UpdateSettings(tc.TenantID, &req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Failed to update GST settings", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("GST settings updated successfully", settings))
}

func (ctrl *GSTController) GetGSTR1Summary(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	month := c.Query("month")
	year := c.Query("year")

	summary, err := ctrl.gstService.GetGSTR1Summary(tc.TenantID, month, year)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to generate GSTR-1 summary", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("GSTR-1 return summary", summary))
}

func (ctrl *GSTController) GetGSTR3B(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	month := c.Query("month")
	year := c.Query("year")

	t31, err := ctrl.gstService.GetGSTR3BTable31(tc.TenantID, month, year)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to generate GSTR-3B table", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("GSTR-3B Table 3.1 details", t31))
}
