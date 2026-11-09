package controllers

import (
	"net/http"
	"strconv"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/middleware"
	"github.com/Abinesh09/KadaiIQ/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsController struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsController(analyticsService *services.AnalyticsService) *AnalyticsController {
	return &AnalyticsController{analyticsService: analyticsService}
}

func (ctrl *AnalyticsController) GetOverview(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	rangeKey := c.Query("range", "30d")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	overview, err := ctrl.analyticsService.GetOverview(tc.TenantID, rangeKey, fromStr, toStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch analytics overview", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Analytics overview retrieved", overview))
}

func (ctrl *AnalyticsController) GetTrends(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	rangeKey := c.Query("range", "30d")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	trends, err := ctrl.analyticsService.GetTrends(tc.TenantID, rangeKey, fromStr, toStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch sales trends", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Sales trends retrieved", trends))
}

func (ctrl *AnalyticsController) GetPaymentMethods(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	rangeKey := c.Query("range", "30d")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	methods, err := ctrl.analyticsService.GetPaymentMethods(tc.TenantID, rangeKey, fromStr, toStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch payment method distribution", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Payment methods distribution retrieved", methods))
}

func (ctrl *AnalyticsController) GetTopProducts(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	rangeKey := c.Query("range", "30d")
	fromStr := c.Query("from")
	toStr := c.Query("to")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	products, err := ctrl.analyticsService.GetTopProducts(tc.TenantID, rangeKey, fromStr, toStr, limit)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch top products", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Top products retrieved", products))
}

func (ctrl *AnalyticsController) GetCategorySales(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	rangeKey := c.Query("range", "30d")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	categories, err := ctrl.analyticsService.GetCategorySales(tc.TenantID, rangeKey, fromStr, toStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch category sales", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Category sales breakdown retrieved", categories))
}

func (ctrl *AnalyticsController) GetFootwearIntelligence(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	rangeKey := c.Query("range", "30d")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	intel, err := ctrl.analyticsService.GetFootwearIntelligence(tc.TenantID, rangeKey, fromStr, toStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch footwear intelligence", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Footwear intelligence retrieved", intel))
}

func (ctrl *AnalyticsController) GetPerformance(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	rangeKey := c.Query("range", "today")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	perf, err := ctrl.analyticsService.GetSalesPerformance(tc.TenantID, rangeKey, fromStr, toStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch sales performance", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Sales performance retrieved", perf))
}

func (ctrl *AnalyticsController) GetDailyBreakdown(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	rangeKey := c.Query("range", "daily")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	daily, err := ctrl.analyticsService.GetDetailedDailyBreakdown(tc.TenantID, rangeKey, fromStr, toStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch daily breakdown", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Daily breakdown retrieved", daily))
}

func (ctrl *AnalyticsController) GetMonthlyBreakdown(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	rangeKey := c.Query("range", "monthly")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	monthly, err := ctrl.analyticsService.GetMonthlyBreakdown(tc.TenantID, rangeKey, fromStr, toStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch monthly breakdown", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Monthly breakdown retrieved", monthly))
}

func (ctrl *AnalyticsController) GetDetailedProducts(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	rangeKey := c.Query("range", "30d")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	products, err := ctrl.analyticsService.GetDetailedProductPerformance(tc.TenantID, rangeKey, fromStr, toStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch product performance", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Product performance retrieved", products))
}

func (ctrl *AnalyticsController) ExportExcel(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	rangeKey := c.Query("range", "today")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	data, filename, err := ctrl.analyticsService.ExportComprehensiveExcel(tc.TenantID, rangeKey, fromStr, toStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to generate Excel analytics export", err.Error()))
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	return c.Send(data)
}
