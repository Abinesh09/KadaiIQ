package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AnalyticsOverviewResponse struct {
	DateRange       string          `json:"date_range"`
	StartDate       time.Time       `json:"start_date"`
	EndDate         time.Time       `json:"end_date"`
	GrossSales      decimal.Decimal `json:"gross_sales"`
	TotalDiscounts  decimal.Decimal `json:"total_discounts"`
	TotalTax        decimal.Decimal `json:"total_tax"`
	NetSales        decimal.Decimal `json:"net_sales"`
	EstimatedCost   decimal.Decimal `json:"estimated_cost"`
	GrossProfit     decimal.Decimal `json:"gross_profit"`
	ProfitMarginPct decimal.Decimal `json:"profit_margin_pct"`
	TotalBills      int64           `json:"total_bills"`
	TotalItemsSold  int64           `json:"total_items_sold"`
	AverageOrderVal decimal.Decimal `json:"average_order_value"`

	// Comparison against previous period
	PrevGrossSales decimal.Decimal `json:"prev_gross_sales"`
	PrevNetSales   decimal.Decimal `json:"prev_net_sales"`
	PrevTotalBills int64           `json:"prev_total_bills"`
	SalesGrowthPct decimal.Decimal `json:"sales_growth_pct"`
	BillsGrowthPct decimal.Decimal `json:"bills_growth_pct"`
}

type TrendPoint struct {
	Period    string          `json:"period"` // e.g. "2026-11-01" or "14:00"
	Timestamp time.Time       `json:"timestamp"`
	Sales     decimal.Decimal `json:"sales"`
	NetSales  decimal.Decimal `json:"net_sales"`
	Profit    decimal.Decimal `json:"profit"`
	Bills     int64           `json:"bills"`
	UnitsSold int64           `json:"units_sold"`
}

type PaymentMethodAnalytics struct {
	Method     string          `json:"method"`
	TotalSales decimal.Decimal `json:"total_sales"`
	BillCount  int64           `json:"bill_count"`
	Percentage decimal.Decimal `json:"percentage"`
}

type ProductSalesAnalytics struct {
	ProductID    *uuid.UUID      `json:"product_id,omitempty"`
	ProductName  string          `json:"product_name"`
	BrandName    string          `json:"brand_name,omitempty"`
	CategoryName string          `json:"category_name,omitempty"`
	UnitsSold    int64           `json:"units_sold"`
	Revenue      decimal.Decimal `json:"revenue"`
	GrossProfit  decimal.Decimal `json:"gross_profit"`
}

type CategorySalesAnalytics struct {
	CategoryName string          `json:"category_name"`
	Revenue      decimal.Decimal `json:"revenue"`
	UnitsSold    int64           `json:"units_sold"`
	Percentage   decimal.Decimal `json:"percentage"`
}

type VariantAttributeAnalytics struct {
	AttributeType string          `json:"attribute_type"` // "size", "colour", or "brand"
	AttributeVal  string          `json:"attribute_value"`
	UnitsSold     int64           `json:"units_sold"`
	Revenue       decimal.Decimal `json:"revenue"`
	Percentage    decimal.Decimal `json:"percentage"`
}

type FootwearIntelligenceResponse struct {
	TopSizes   []VariantAttributeAnalytics `json:"top_sizes"`
	TopColours []VariantAttributeAnalytics `json:"top_colours"`
	TopBrands  []VariantAttributeAnalytics `json:"top_brands"`
}
