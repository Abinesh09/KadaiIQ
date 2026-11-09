package controllers

import (
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/middleware"
	"github.com/Abinesh09/KadaiIQ/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SaleController struct {
	saleService *services.SaleService
}

func NewSaleController(saleService *services.SaleService) *SaleController {
	return &SaleController{saleService: saleService}
}

func (ctrl *SaleController) ListSales(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	search := c.Query("search")
	status := c.Query("status")

	var dateFrom *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse("2006-01-02", fromStr); err == nil {
			dateFrom = &t
		}
	}

	var dateTo *time.Time
	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse("2006-01-02", toStr); err == nil {
			// Include full day
			tEnd := t.Add(24*time.Hour - time.Nanosecond)
			dateTo = &tEnd
		}
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	res, err := ctrl.saleService.ListSales(tc.TenantID, search, status, dateFrom, dateTo, page, pageSize)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to list sales", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Sales records retrieved", res))
}

func (ctrl *SaleController) GetSale(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	saleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid sale ID", err.Error()))
	}

	sale, err := ctrl.saleService.GetSale(tc.TenantID, saleID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(dto.ErrorResponse("Sale not found", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Sale bill details", sale))
}

func (ctrl *SaleController) PreviewImport(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "File upload is required (key: 'file')"))
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Failed to read file", err.Error()))
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Failed to process file bytes", err.Error()))
	}

	preview, err := ctrl.saleService.PreviewImport(tc.TenantID, content, file.Filename)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("File preview failed", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("File parsed for preview", preview))
}

func (ctrl *SaleController) ExecuteImport(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	var content []byte
	var fileName string
	sourceSystem := c.FormValue("source_system", "ZenBiller")
	deductInventory := c.FormValue("deduct_inventory") == "true"

	// Check if uploaded as multipart file
	file, fileErr := c.FormFile("file")
	if fileErr == nil {
		f, err := file.Open()
		if err == nil {
			defer f.Close()
			content, _ = io.ReadAll(f)
			fileName = file.Filename
		}
	}

	// Or check JSON payload with raw CSV data
	if len(content) == 0 {
		var req dto.ExecuteImportRequest
		if err := c.BodyParser(&req); err == nil && req.CSVData != "" {
			fileName = req.FileName
			if fileName == "" {
				fileName = "sales_import.csv"
			}
			sourceSystem = req.SourceSystem
			deductInventory = req.DeductInventory
			decoded, err := base64.StdEncoding.DecodeString(req.CSVData)
			if err == nil {
				content = decoded
			} else {
				content = []byte(req.CSVData)
			}
		}
	}

	if len(content) == 0 {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "No file or CSV data provided"))
	}

	job, err := ctrl.saleService.ExecuteImport(tc.TenantID, tc.UserID, content, fileName, sourceSystem, deductInventory)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Sales import failed", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Import completed", job))
}

func (ctrl *SaleController) GetImportHistory(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	history, err := ctrl.saleService.GetImportHistory(tc.TenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch import history", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Import history", history))
}
