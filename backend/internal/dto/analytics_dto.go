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

type BestSellingMetric struct {
	Label string          `json:"label"` // e.g. "09 Nov 2026" or "November"
	Sales decimal.Decimal `json:"sales"`
	Units int64           `json:"units"`
}

type SalesPerformanceSummary struct {
	PeriodType      string             `json:"period_type"`
	StartDate       time.Time          `json:"start_date"`
	EndDate         time.Time          `json:"end_date"`
	GrossSales      decimal.Decimal    `json:"gross_sales"`
	TotalDiscounts  decimal.Decimal    `json:"total_discounts"`
	TotalReturns    decimal.Decimal    `json:"total_returns"`
	TotalTax        decimal.Decimal    `json:"total_tax"`
	NetSales        decimal.Decimal    `json:"net_sales"`        // Gross - Discounts - Returns
	AmountCollected decimal.Decimal    `json:"amount_collected"` // Tenders collected
	CostOfGoodsSold decimal.Decimal    `json:"cogs"`
	GrossProfit     decimal.Decimal    `json:"gross_profit"` // NetSales - COGS
	ProfitMarginPct decimal.Decimal    `json:"profit_margin_pct"`
	TotalBills      int64              `json:"total_bills"`
	TotalUnitsSold  int64              `json:"total_units_sold"`
	AverageBillVal  decimal.Decimal    `json:"average_bill_value"`
	AverageDailyVal decimal.Decimal    `json:"average_daily_sales"`
	BestSellingDay  *BestSellingMetric `json:"best_selling_day,omitempty"`
	BestSellingMo   *BestSellingMetric `json:"best_selling_month,omitempty"`

	// Comparison with previous period
	PrevGrossSales decimal.Decimal `json:"prev_gross_sales"`
	PrevNetSales   decimal.Decimal `json:"prev_net_sales"`
	PrevUnitsSold  int64           `json:"prev_units_sold"`
	PrevTotalBills int64           `json:"prev_total_bills"`
	SalesGrowthPct decimal.Decimal `json:"sales_growth_pct"`
	BillsGrowthPct decimal.Decimal `json:"bills_growth_pct"`
}

type DailyBreakdownItem struct {
	Date       string          `json:"date"`
	Bills      int64           `json:"bills"`
	UnitsSold  int64           `json:"units_sold"`
	GrossSales decimal.Decimal `json:"gross_sales"`
	Discounts  decimal.Decimal `json:"discounts"`
	Returns    decimal.Decimal `json:"returns"`
	TaxAmount  decimal.Decimal `json:"tax_amount"`
	NetSales   decimal.Decimal `json:"net_sales"`
}

type MonthlyBreakdownItem struct {
	Month      int             `json:"month"`
	MonthName  string          `json:"month_name"`
	Year       int             `json:"year"`
	Bills      int64           `json:"bills"`
	UnitsSold  int64           `json:"units_sold"`
	GrossSales decimal.Decimal `json:"gross_sales"`
	Discounts  decimal.Decimal `json:"discounts"`
	Returns    decimal.Decimal `json:"returns"`
	NetSales   decimal.Decimal `json:"net_sales"`
}

type DetailedProductPerformanceItem struct {
	ProductID    *uuid.UUID      `json:"product_id,omitempty"`
	ProductName  string          `json:"product_name"`
	ProductCode  string          `json:"product_code"`
	SKU          string          `json:"sku,omitempty"`
	Barcode      string          `json:"barcode,omitempty"`
	BrandName    string          `json:"brand_name"`
	CategoryName string          `json:"category_name"`
	Size         string          `json:"size,omitempty"`
	Colour       string          `json:"colour,omitempty"`
	UnitsSold    int64           `json:"units_sold"`
	GrossSales   decimal.Decimal `json:"gross_sales"`
	Discounts    decimal.Decimal `json:"discounts"`
	Returns      decimal.Decimal `json:"returns"`
	NetSales     decimal.Decimal `json:"net_sales"`
	CostOfGoods  decimal.Decimal `json:"cost_of_goods"`
	GrossProfit  decimal.Decimal `json:"gross_profit"`
	ProfitMargin decimal.Decimal `json:"profit_margin"`
	StockOnHand  int             `json:"stock_on_hand"`
	Status       string          `json:"status"` // "BEST_SELLER", "ACTIVE", "SLOW_MOVING", "ZERO_SALES"
}

type ProductPerformanceResponse struct {
	TopSelling    []DetailedProductPerformanceItem `json:"top_selling"`
	SlowMoving    []DetailedProductPerformanceItem `json:"slow_moving"`
	ZeroSales     []DetailedProductPerformanceItem `json:"zero_sales"`
	AllProducts   []DetailedProductPerformanceItem `json:"all_products"`
	TotalProducts int64                            `json:"total_products"`
}
