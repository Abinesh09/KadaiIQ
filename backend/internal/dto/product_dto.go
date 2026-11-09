package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateVariantRequest struct {
	Size              string          `json:"size"`
	Colour            string          `json:"colour"`
	SKU               string          `json:"sku"`
	Barcode           string          `json:"barcode"`
	PurchasePrice     decimal.Decimal `json:"purchase_price"`
	SellingPrice      decimal.Decimal `json:"selling_price"`
	GSTRate           decimal.Decimal `json:"gst_rate"`
	OpeningStock      int             `json:"opening_stock"`
	MinStockThreshold int             `json:"min_stock_threshold"`
}

type CreateProductRequest struct {
	Name        string                 `json:"name"`
	ProductCode string                 `json:"product_code"`
	SKU         string                 `json:"sku,omitempty"`
	Barcode     string                 `json:"barcode,omitempty"`
	BrandID     *uuid.UUID             `json:"brand_id,omitempty"`
	CategoryID  *uuid.UUID             `json:"category_id,omitempty"`
	Description string                 `json:"description,omitempty"`
	Gender      string                 `json:"gender"`
	Variants    []CreateVariantRequest `json:"variants"`
}

type UpdateProductRequest struct {
	Name        string     `json:"name,omitempty"`
	BrandID     *uuid.UUID `json:"brand_id,omitempty"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	Description string     `json:"description,omitempty"`
	Gender      string     `json:"gender,omitempty"`
	Status      string     `json:"status,omitempty"`
}

type ProductVariantResponse struct {
	ID                uuid.UUID       `json:"id"`
	ProductID         uuid.UUID       `json:"product_id"`
	Size              string          `json:"size"`
	Colour            string          `json:"colour"`
	SKU               string          `json:"sku"`
	Barcode           string          `json:"barcode"`
	PurchasePrice     decimal.Decimal `json:"purchase_price"`
	SellingPrice      decimal.Decimal `json:"selling_price"`
	GSTRate           decimal.Decimal `json:"gst_rate"`
	CurrentStock      int             `json:"current_stock"`
	MinStockThreshold int             `json:"min_stock_threshold"`
	StockStatus       string          `json:"stock_status"` // "IN_STOCK", "LOW_STOCK", "OUT_OF_STOCK"
}

type ProductResponse struct {
	ID           uuid.UUID                `json:"id"`
	Name         string                   `json:"name"`
	ProductCode  string                   `json:"product_code"`
	SKU          string                   `json:"sku,omitempty"`
	Barcode      string                   `json:"barcode,omitempty"`
	BrandName    string                   `json:"brand_name,omitempty"`
	CategoryName string                   `json:"category_name,omitempty"`
	BrandID      *uuid.UUID               `json:"brand_id,omitempty"`
	CategoryID   *uuid.UUID               `json:"category_id,omitempty"`
	Description  string                   `json:"description,omitempty"`
	Gender       string                   `json:"gender"`
	Status       string                   `json:"status"`
	TotalStock   int                      `json:"total_stock"`
	Variants     []ProductVariantResponse `json:"variants"`
	CreatedAt    time.Time                `json:"created_at"`
}

type ProductListResponse struct {
	Items      []ProductResponse `json:"items"`
	TotalCount int64             `json:"total_count"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
}

type StockAdjustmentRequest struct {
	VariantID       uuid.UUID       `json:"variant_id"`
	MovementType    string          `json:"movement_type"` // "PURCHASE", "SALE", "CUSTOMER_RETURN", "MANUAL_ADJUSTMENT"
	Quantity        int             `json:"quantity"`      // positive integer representing amount moved
	UnitCost        decimal.Decimal `json:"unit_cost,omitempty"`
	ReferenceNumber string          `json:"reference_number,omitempty"`
	Notes           string          `json:"notes,omitempty"`
}

type StockMovementResponse struct {
	ID              uuid.UUID       `json:"id"`
	VariantID       uuid.UUID       `json:"variant_id"`
	ProductName     string          `json:"product_name"`
	VariantDetails  string          `json:"variant_details"` // e.g. "Size 8 / Black"
	Barcode         string          `json:"barcode"`
	MovementType    string          `json:"movement_type"`
	Quantity        int             `json:"quantity"`
	PreviousStock   int             `json:"previous_stock"`
	NewStock        int             `json:"new_stock"`
	UnitCost        decimal.Decimal `json:"unit_cost"`
	ReferenceNumber string          `json:"reference_number,omitempty"`
	Notes           string          `json:"notes,omitempty"`
	CreatedByName   string          `json:"created_by_name"`
	CreatedAt       time.Time       `json:"created_at"`
}

type InventorySummaryResponse struct {
	TotalProducts   int64           `json:"total_products"`
	TotalVariants   int64           `json:"total_variants"`
	TotalStockUnits int64           `json:"total_stock_units"`
	LowStockCount   int64           `json:"low_stock_count"`
	OutOfStockCount int64           `json:"out_of_stock_count"`
	StockValuation  decimal.Decimal `json:"stock_valuation"`
}
