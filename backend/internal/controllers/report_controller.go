package controllers

import (
	"fmt"
	"net/http"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/middleware"
	"github.com/Abinesh09/KadaiIQ/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ReportController struct {
	reportService *services.ReportService
}

func NewReportController(reportService *services.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

func (ctrl *ReportController) GetSalesSummary(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")

	report, err := ctrl.reportService.GetSalesSummary(tc.TenantID, fromStr, toStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to generate sales report", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Sales summary report", report))
}

func (ctrl *ReportController) GetStockValuation(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	report, err := ctrl.reportService.GetStockValuation(tc.TenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to generate stock valuation report", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Stock valuation report", report))
}

func (ctrl *ReportController) GetStockMovements(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")

	report, err := ctrl.reportService.GetStockMovements(tc.TenantID, fromStr, toStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to generate stock movements report", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Stock movements report", report))
}

func (ctrl *ReportController) DownloadReport(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	reportType := c.Query("type", "sales")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var data []byte
	var fileName string

	switch reportType {
	case "inventory":
		data, fileName, err = ctrl.reportService.ExportStockValuationCSV(tc.TenantID)
	case "movements":
		data, fileName, err = ctrl.reportService.ExportStockMovementCSV(tc.TenantID, fromStr, toStr)
	default: // "sales"
		data, fileName, err = ctrl.reportService.ExportSalesCSV(tc.TenantID, fromStr, toStr)
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to export report", err.Error()))
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	return c.Send(data)
}
