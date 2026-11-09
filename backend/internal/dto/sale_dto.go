package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type SaleItemResponse struct {
	ID          uuid.UUID       `json:"id"`
	ProductName string          `json:"product_name"`
	SKU         string          `json:"sku,omitempty"`
	Barcode     string          `json:"barcode,omitempty"`
	VariantID   *uuid.UUID      `json:"variant_id,omitempty"`
	Quantity    int             `json:"quantity"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	Discount    decimal.Decimal `json:"discount"`
	TaxRate     decimal.Decimal `json:"tax_rate"`
	TaxAmount   decimal.Decimal `json:"tax_amount"`
	LineTotal   decimal.Decimal `json:"line_total"`
}

type SaleResponse struct {
	ID                uuid.UUID          `json:"id"`
	SourceSystem      string             `json:"source_system"`
	SourceBillID      string             `json:"source_bill_id"`
	BillNumber        string             `json:"bill_number"`
	SaleDate          time.Time          `json:"sale_date"`
	Subtotal          decimal.Decimal    `json:"subtotal"`
	Discount          decimal.Decimal    `json:"discount"`
	TaxAmount         decimal.Decimal    `json:"tax_amount"`
	GrandTotal        decimal.Decimal    `json:"grand_total"`
	PaymentMethod     string             `json:"payment_method"`
	Status            string             `json:"status"`
	InventoryDeducted bool               `json:"inventory_deducted"`
	ItemCount         int                `json:"item_count"`
	Items             []SaleItemResponse `json:"items,omitempty"`
	CreatedAt         time.Time          `json:"created_at"`
}

type SaleListResponse struct {
	Items      []SaleResponse `json:"items"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

type ImportJobResponse struct {
	ID               uuid.UUID `json:"id"`
	SourceSystem     string    `json:"source_system"`
	FileName         string    `json:"file_name"`
	TotalRecords     int       `json:"total_records"`
	ImportedRecords  int       `json:"imported_records"`
	FailedRecords    int       `json:"failed_records"`
	DuplicateRecords int       `json:"duplicate_records"`
	Status           string    `json:"status"`
	ErrorDetails     string    `json:"error_details,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type ImportPreviewItem struct {
	BillNumber  string `json:"bill_number"`
	Date        string `json:"date"`
	Barcode     string `json:"barcode"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	UnitPrice   string `json:"unit_price"`
	Total       string `json:"total"`
	Matched     bool   `json:"matched"`
}

type ImportPreviewResponse struct {
	TotalRows       int                 `json:"total_rows"`
	DetectedColumns []string            `json:"detected_columns"`
	SampleRows      []ImportPreviewItem `json:"sample_rows"`
}

type ExecuteImportRequest struct {
	SourceSystem    string `json:"source_system"`
	FileName        string `json:"file_name"`
	CSVData         string `json:"csv_data"` // Base64 encoded or raw text
	DeductInventory bool   `json:"deduct_inventory"`
}
