package controllers

import (
	"net/http"
	"strconv"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/middleware"
	"github.com/Abinesh09/KadaiIQ/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductController struct {
	productService *services.ProductService
}

func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

func (ctrl *ProductController) ListProducts(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	search := c.Query("search")
	gender := c.Query("gender")
	status := c.Query("status")

	var brandID *uuid.UUID
	if bStr := c.Query("brand_id"); bStr != "" {
		if bUUID, err := uuid.Parse(bStr); err == nil {
			brandID = &bUUID
		}
	}

	var categoryID *uuid.UUID
	if cStr := c.Query("category_id"); cStr != "" {
		if cUUID, err := uuid.Parse(cStr); err == nil {
			categoryID = &cUUID
		}
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	res, err := ctrl.productService.ListProducts(tc.TenantID, search, brandID, categoryID, gender, status, page, pageSize)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to list products", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Products retrieved", res))
}

func (ctrl *ProductController) GetProduct(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	prodID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid product ID", err.Error()))
	}

	res, err := ctrl.productService.GetProduct(tc.TenantID, prodID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(dto.ErrorResponse("Product not found", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Product details", res))
}

func (ctrl *ProductController) CreateProduct(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid request body", err.Error()))
	}

	res, err := ctrl.productService.CreateProduct(tc.TenantID, tc.UserID, req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Product creation failed", err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(dto.SuccessResponse("Product created successfully", res))
}

func (ctrl *ProductController) ArchiveProduct(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	prodID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid product ID", err.Error()))
	}

	if err := ctrl.productService.ArchiveProduct(tc.TenantID, prodID); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Failed to archive product", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Product archived successfully", nil))
}

func (ctrl *ProductController) LookupBarcode(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	barcode := c.Params("barcode")
	if barcode == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "Barcode is required"))
	}

	variant, err := ctrl.productService.LookupBarcode(tc.TenantID, barcode)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(dto.ErrorResponse("Variant not found", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Barcode identified", variant))
}

func (ctrl *ProductController) GetBrands(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	brands, err := ctrl.productService.GetBrands(tc.TenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch brands", err.Error()))
	}
	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Brands retrieved", brands))
}

func (ctrl *ProductController) GetCategories(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	cats, err := ctrl.productService.GetCategories(tc.TenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch categories", err.Error()))
	}
	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Categories retrieved", cats))
}

// Inventory Handlers
func (ctrl *ProductController) GetInventorySummary(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	summary, err := ctrl.productService.GetInventorySummary(tc.TenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch inventory metrics", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Inventory overview", summary))
}

func (ctrl *ProductController) RecordStockAdjustment(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	var req dto.StockAdjustmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid request body", err.Error()))
	}
	if req.VariantID == uuid.Nil || req.MovementType == "" || req.Quantity <= 0 {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "variant_id, movement_type, and positive quantity are required"))
	}

	movement, err := ctrl.productService.RecordStockAdjustment(tc.TenantID, tc.UserID, req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Stock adjustment failed", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Stock adjusted and movement recorded", movement))
}

func (ctrl *ProductController) GetStockMovements(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	var variantID *uuid.UUID
	if vStr := c.Query("variant_id"); vStr != "" {
		if vUUID, err := uuid.Parse(vStr); err == nil {
			variantID = &vUUID
		}
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	items, total, err := ctrl.productService.GetStockMovements(tc.TenantID, variantID, page, pageSize)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch stock movements", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Stock movement ledger", fiber.Map{
		"items":       items,
		"total_count": total,
		"page":        page,
		"page_size":   pageSize,
	}))
}
