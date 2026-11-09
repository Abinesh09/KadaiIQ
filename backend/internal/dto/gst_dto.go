package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GSTSettingsResponse struct {
	ID                 uuid.UUID       `json:"id"`
	TenantID           uuid.UUID       `json:"tenant_id"`
	GSTIN              string          `json:"gstin"`
	LegalName          string          `json:"legal_name"`
	TradeName          string          `json:"trade_name"`
	StateCode          string          `json:"state_code"`
	StateName          string          `json:"state_name"`
	RegistrationType   string          `json:"registration_type"`
	DefaultFootwearHSN string          `json:"default_footwear_hsn"`
	DefaultGSTRate     decimal.Decimal `json:"default_gst_rate"`
	EnableComposition  bool            `json:"enable_composition"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type UpdateGSTSettingsRequest struct {
	GSTIN              string          `json:"gstin"`
	LegalName          string          `json:"legal_name"`
	TradeName          string          `json:"trade_name"`
	RegistrationType   string          `json:"registration_type"`
	DefaultFootwearHSN string          `json:"default_footwear_hsn"`
	DefaultGSTRate     decimal.Decimal `json:"default_gst_rate"`
	EnableComposition  bool            `json:"enable_composition"`
}

type GSTR1RateSlab struct {
	RatePct      decimal.Decimal `json:"rate_pct"`
	TaxableValue decimal.Decimal `json:"taxable_value"`
	CGST         decimal.Decimal `json:"cgst"`
	SGST         decimal.Decimal `json:"sgst"`
	IGST         decimal.Decimal `json:"igst"`
	TotalTax     decimal.Decimal `json:"total_tax"`
	Invoices     int64           `json:"invoices"`
}

type GSTR1HSNRow struct {
	HSNCode      string          `json:"hsn_code"`
	Description  string          `json:"description"`
	UQC          string          `json:"uqc"` // Unit Quantity Code e.g. "PRS" (Pairs)
	TotalQty     int64           `json:"total_qty"`
	TaxableValue decimal.Decimal `json:"taxable_value"`
	CGST         decimal.Decimal `json:"cgst"`
	SGST         decimal.Decimal `json:"sgst"`
	IGST         decimal.Decimal `json:"igst"`
	TotalTax     decimal.Decimal `json:"total_tax"`
}

type GSTR1SummaryResponse struct {
	ReturnPeriod      string          `json:"return_period"`
	StartDate         time.Time       `json:"start_date"`
	EndDate           time.Time       `json:"end_date"`
	GSTIN             string          `json:"gstin"`
	LegalName         string          `json:"legal_name"`
	StateCode         string          `json:"state_code"`
	TotalTaxableValue decimal.Decimal `json:"total_taxable_value"`
	TotalCGST         decimal.Decimal `json:"total_cgst"`
	TotalSGST         decimal.Decimal `json:"total_sgst"`
	TotalIGST         decimal.Decimal `json:"total_igst"`
	TotalTaxLiability decimal.Decimal `json:"total_tax_liability"`
	B2CSmallCount     int64           `json:"b2c_small_count"`
	RateBreakdown     []GSTR1RateSlab `json:"rate_breakdown"`
	HSNSummary        []GSTR1HSNRow   `json:"hsn_summary"`
}

type GSTR3BTable31Response struct {
	NatureOfSupplies string          `json:"nature_of_supplies"`
	TaxableValue     decimal.Decimal `json:"taxable_value"`
	IntegratedTax    decimal.Decimal `json:"integrated_tax"`
	CentralTax       decimal.Decimal `json:"central_tax"`
	StateUTTax       decimal.Decimal `json:"state_ut_tax"`
	Cess             decimal.Decimal `json:"cess"`
}
