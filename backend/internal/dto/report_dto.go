package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type SalesSummaryRow struct {
	BillNumber    string          `json:"bill_number"`
	SaleDate      time.Time       `json:"sale_date"`
	PaymentMethod string          `json:"payment_method"`
	ItemCount     int             `json:"item_count"`
	Subtotal      decimal.Decimal `json:"subtotal"`
	Discount      decimal.Decimal `json:"discount"`
	TaxAmount     decimal.Decimal `json:"tax_amount"`
	GrandTotal    decimal.Decimal `json:"grand_total"`
	Status        string          `json:"status"`
}

type StockValuationRow struct {
	ProductName     string          `json:"product_name"`
	ProductCode     string          `json:"product_code"`
	Size            string          `json:"size"`
	Colour          string          `json:"colour"`
	SKU             string          `json:"sku"`
	Barcode         string          `json:"barcode"`
	CurrentStock    int             `json:"current_stock"`
	PurchasePrice   decimal.Decimal `json:"purchase_price"`
	SellingPrice    decimal.Decimal `json:"selling_price"`
	CostValuation   decimal.Decimal `json:"cost_valuation"`
	RetailValuation decimal.Decimal `json:"retail_valuation"`
	StockStatus     string          `json:"stock_status"`
}

type StockMovementReportRow struct {
	Date            time.Time       `json:"date"`
	ProductName     string          `json:"product_name"`
	Barcode         string          `json:"barcode"`
	VariantDetails  string          `json:"variant_details"`
	MovementType    string          `json:"movement_type"`
	Quantity        int             `json:"quantity"`
	PreviousStock   int             `json:"previous_stock"`
	NewStock        int             `json:"new_stock"`
	UnitCost        decimal.Decimal `json:"unit_cost"`
	ReferenceNumber string          `json:"reference_number,omitempty"`
	Notes           string          `json:"notes,omitempty"`
	CreatedBy       string          `json:"created_by"`
}

type ReportResponse[T any] struct {
	ReportName   string          `json:"report_name"`
	GeneratedAt  time.Time       `json:"generated_at"`
	DateFrom     *time.Time      `json:"date_from,omitempty"`
	DateTo       *time.Time      `json:"date_to,omitempty"`
	TotalRecords int             `json:"total_records"`
	SummaryValue decimal.Decimal `json:"summary_value,omitempty"`
	Data         []T             `json:"data"`
}
