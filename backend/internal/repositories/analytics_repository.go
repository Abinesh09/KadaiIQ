package repositories

import (
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AnalyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

type PeriodAggregates struct {
	TotalBills     int64           `gorm:"column:total_bills"`
	GrossSales     decimal.Decimal `gorm:"column:gross_sales"`
	TotalDiscounts decimal.Decimal `gorm:"column:total_discounts"`
	TotalTax       decimal.Decimal `gorm:"column:total_tax"`
	NetSales       decimal.Decimal `gorm:"column:net_sales"`
}

type ItemsAggregates struct {
	TotalUnitsSold int64           `gorm:"column:total_units_sold"`
	EstimatedCost  decimal.Decimal `gorm:"column:estimated_cost"`
}

// GetSalesAggregates retrieves gross, discount, tax, and net sales totals for a given date window.
func (r *AnalyticsRepository) GetSalesAggregates(tenantID uuid.UUID, start, end time.Time) (*PeriodAggregates, error) {
	var agg PeriodAggregates
	err := r.db.Table("sales").
		Select(`
			COUNT(id) as total_bills,
			COALESCE(SUM(grand_total), 0) as gross_sales,
			COALESCE(SUM(discount), 0) as total_discounts,
			COALESCE(SUM(tax_amount), 0) as total_tax,
			COALESCE(SUM(subtotal), 0) as net_sales
		`).
		Where("tenant_id = ? AND sale_date >= ? AND sale_date <= ? AND status != 'CANCELLED' AND status != 'RETURNED'", tenantID, start, end).
		Scan(&agg).Error

	if err != nil {
		return nil, err
	}
	return &agg, nil
}

// GetItemsAggregates retrieves units sold and cost of goods sold (COGS) for matched variants.
func (r *AnalyticsRepository) GetItemsAggregates(tenantID uuid.UUID, start, end time.Time) (*ItemsAggregates, error) {
	var agg ItemsAggregates
	err := r.db.Table("sale_items").
		Joins("JOIN sales ON sale_items.sale_id = sales.id").
		Joins("LEFT JOIN product_variants ON sale_items.variant_id = product_variants.id").
		Select(`
			COALESCE(SUM(sale_items.quantity), 0) as total_units_sold,
			COALESCE(SUM(product_variants.purchase_price * sale_items.quantity), 0) as estimated_cost
		`).
		Where("sale_items.tenant_id = ? AND sales.sale_date >= ? AND sales.sale_date <= ? AND sales.status != 'CANCELLED' AND sales.status != 'RETURNED'", tenantID, start, end).
		Scan(&agg).Error

	if err != nil {
		return nil, err
	}
	return &agg, nil
}

type TrendRow struct {
	Period    string          `gorm:"column:period"`
	Sales     decimal.Decimal `gorm:"column:sales"`
	NetSales  decimal.Decimal `gorm:"column:net_sales"`
	Bills     int64           `gorm:"column:bills"`
	UnitsSold int64           `gorm:"column:units_sold"`
}

// GetDailyTrends aggregates metrics grouped by DATE(sale_date) within range.
func (r *AnalyticsRepository) GetDailyTrends(tenantID uuid.UUID, start, end time.Time) ([]TrendRow, error) {
	var rows []TrendRow
	err := r.db.Table("sales").
		Select(`
			DATE(sales.sale_date) as period,
			COALESCE(SUM(sales.grand_total), 0) as sales,
			COALESCE(SUM(sales.subtotal), 0) as net_sales,
			COUNT(sales.id) as bills
		`).
		Where("sales.tenant_id = ? AND sales.sale_date >= ? AND sales.sale_date <= ? AND sales.status != 'CANCELLED'", tenantID, start, end).
		Group("DATE(sales.sale_date)").
		Order("period ASC").
		Scan(&rows).Error

	return rows, err
}

type RawPaymentMethod struct {
	Method     string          `gorm:"column:method"`
	TotalSales decimal.Decimal `gorm:"column:total_sales"`
	BillCount  int64           `gorm:"column:bill_count"`
}

// GetPaymentMethodDistribution calculates sales breakdown per payment mode.
func (r *AnalyticsRepository) GetPaymentMethodDistribution(tenantID uuid.UUID, start, end time.Time) ([]RawPaymentMethod, error) {
	var results []RawPaymentMethod
	err := r.db.Table("sales").
		Select(`
			COALESCE(payment_method, 'Cash') as method,
			COALESCE(SUM(grand_total), 0) as total_sales,
			COUNT(id) as bill_count
		`).
		Where("tenant_id = ? AND sale_date >= ? AND sale_date <= ? AND status != 'CANCELLED'", tenantID, start, end).
		Group("method").
		Order("total_sales DESC").
		Scan(&results).Error

	return results, err
}

type RawProductSales struct {
	ProductName  string          `gorm:"column:product_name"`
	BrandName    string          `gorm:"column:brand_name"`
	CategoryName string          `gorm:"column:category_name"`
	UnitsSold    int64           `gorm:"column:units_sold"`
	Revenue      decimal.Decimal `gorm:"column:revenue"`
	EstCost      decimal.Decimal `gorm:"column:est_cost"`
}

// GetTopProducts returns top performing products ranked by revenue.
func (r *AnalyticsRepository) GetTopProducts(tenantID uuid.UUID, start, end time.Time, limit int) ([]RawProductSales, error) {
	var results []RawProductSales
	err := r.db.Table("sale_items").
		Joins("JOIN sales ON sale_items.sale_id = sales.id").
		Joins("LEFT JOIN product_variants ON sale_items.variant_id = product_variants.id").
		Joins("LEFT JOIN products ON product_variants.product_id = products.id").
		Joins("LEFT JOIN brands ON products.brand_id = brands.id").
		Joins("LEFT JOIN categories ON products.category_id = categories.id").
		Select(`
			sale_items.product_name,
			COALESCE(brands.name, 'Generic') as brand_name,
			COALESCE(categories.name, 'Footwear') as category_name,
			COALESCE(SUM(sale_items.quantity), 0) as units_sold,
			COALESCE(SUM(sale_items.line_total), 0) as revenue,
			COALESCE(SUM(product_variants.purchase_price * sale_items.quantity), 0) as est_cost
		`).
		Where("sale_items.tenant_id = ? AND sales.sale_date >= ? AND sales.sale_date <= ? AND sales.status != 'CANCELLED'", tenantID, start, end).
		Group("sale_items.product_name, brands.name, categories.name").
		Order("revenue DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

type RawCategorySales struct {
	CategoryName string          `gorm:"column:category_name"`
	Revenue      decimal.Decimal `gorm:"column:revenue"`
	UnitsSold    int64           `gorm:"column:units_sold"`
}

// GetCategorySales returns revenue grouped by product category.
func (r *AnalyticsRepository) GetCategorySales(tenantID uuid.UUID, start, end time.Time) ([]RawCategorySales, error) {
	var results []RawCategorySales
	err := r.db.Table("sale_items").
		Joins("JOIN sales ON sale_items.sale_id = sales.id").
		Joins("LEFT JOIN product_variants ON sale_items.variant_id = product_variants.id").
		Joins("LEFT JOIN products ON product_variants.product_id = products.id").
		Joins("LEFT JOIN categories ON products.category_id = categories.id").
		Select(`
			COALESCE(categories.name, 'General Footwear') as category_name,
			COALESCE(SUM(sale_items.line_total), 0) as revenue,
			COALESCE(SUM(sale_items.quantity), 0) as units_sold
		`).
		Where("sale_items.tenant_id = ? AND sales.sale_date >= ? AND sales.sale_date <= ? AND sales.status != 'CANCELLED'", tenantID, start, end).
		Group("categories.name").
		Order("revenue DESC").
		Scan(&results).Error

	return results, err
}

type RawAttributeSales struct {
	AttributeVal string          `gorm:"column:attribute_val"`
	UnitsSold    int64           `gorm:"column:units_sold"`
	Revenue      decimal.Decimal `gorm:"column:revenue"`
}

// GetTopSizes returns units and revenue grouped by footwear size.
func (r *AnalyticsRepository) GetTopSizes(tenantID uuid.UUID, start, end time.Time, limit int) ([]RawAttributeSales, error) {
	var results []RawAttributeSales
	err := r.db.Table("sale_items").
		Joins("JOIN sales ON sale_items.sale_id = sales.id").
		Joins("JOIN product_variants ON sale_items.variant_id = product_variants.id").
		Select(`
			product_variants.size as attribute_val,
			COALESCE(SUM(sale_items.quantity), 0) as units_sold,
			COALESCE(SUM(sale_items.line_total), 0) as revenue
		`).
		Where("sale_items.tenant_id = ? AND sales.sale_date >= ? AND sales.sale_date <= ? AND sales.status != 'CANCELLED' AND product_variants.size != ''", tenantID, start, end).
		Group("product_variants.size").
		Order("units_sold DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// GetTopColours returns units and revenue grouped by colour.
func (r *AnalyticsRepository) GetTopColours(tenantID uuid.UUID, start, end time.Time, limit int) ([]RawAttributeSales, error) {
	var results []RawAttributeSales
	err := r.db.Table("sale_items").
		Joins("JOIN sales ON sale_items.sale_id = sales.id").
		Joins("JOIN product_variants ON sale_items.variant_id = product_variants.id").
		Select(`
			product_variants.colour as attribute_val,
			COALESCE(SUM(sale_items.quantity), 0) as units_sold,
			COALESCE(SUM(sale_items.line_total), 0) as revenue
		`).
		Where("sale_items.tenant_id = ? AND sales.sale_date >= ? AND sales.sale_date <= ? AND sales.status != 'CANCELLED' AND product_variants.colour != ''", tenantID, start, end).
		Group("product_variants.colour").
		Order("units_sold DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// GetTopBrands returns units and revenue grouped by brand.
func (r *AnalyticsRepository) GetTopBrands(tenantID uuid.UUID, start, end time.Time, limit int) ([]RawAttributeSales, error) {
	var results []RawAttributeSales
	err := r.db.Table("sale_items").
		Joins("JOIN sales ON sale_items.sale_id = sales.id").
		Joins("JOIN product_variants ON sale_items.variant_id = product_variants.id").
		Joins("JOIN products ON product_variants.product_id = products.id").
		Joins("LEFT JOIN brands ON products.brand_id = brands.id").
		Select(`
			COALESCE(brands.name, 'Generic') as attribute_val,
			COALESCE(SUM(sale_items.quantity), 0) as units_sold,
			COALESCE(SUM(sale_items.line_total), 0) as revenue
		`).
		Where("sale_items.tenant_id = ? AND sales.sale_date >= ? AND sales.sale_date <= ? AND sales.status != 'CANCELLED'", tenantID, start, end).
		Group("COALESCE(brands.name, 'Generic')").
		Order("revenue DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

type ReturnMetrics struct {
	TotalReturns decimal.Decimal `gorm:"column:total_returns"`
	ReturnBills  int64           `gorm:"column:return_bills"`
}

// GetReturnsMetrics calculates total sales returns
func (r *AnalyticsRepository) GetReturnsMetrics(tenantID uuid.UUID, start, end time.Time) (*ReturnMetrics, error) {
	var m ReturnMetrics
	err := r.db.Table("sales").
		Select(`
			COALESCE(SUM(grand_total), 0) as total_returns,
			COUNT(id) as return_bills
		`).
		Where("tenant_id = ? AND sale_date >= ? AND sale_date <= ? AND status = 'RETURNED'", tenantID, start, end).
		Scan(&m).Error

	return &m, err
}

type DetailedDailyRow struct {
	Date       string          `gorm:"column:date"`
	Bills      int64           `gorm:"column:bills"`
	UnitsSold  int64           `gorm:"column:units_sold"`
	GrossSales decimal.Decimal `gorm:"column:gross_sales"`
	Discounts  decimal.Decimal `gorm:"column:discounts"`
	Returns    decimal.Decimal `gorm:"column:returns"`
	TaxAmount  decimal.Decimal `gorm:"column:tax_amount"`
	NetSales   decimal.Decimal `gorm:"column:net_sales"`
}

// GetDetailedDailyBreakdown returns date-by-date breakdown with discounts and returns
func (r *AnalyticsRepository) GetDetailedDailyBreakdown(tenantID uuid.UUID, start, end time.Time) ([]DetailedDailyRow, error) {
	var rows []DetailedDailyRow
	err := r.db.Table("sales").
		Select(`
			DATE(sales.sale_date) as date,
			COUNT(CASE WHEN UPPER(sales.status) != 'CANCELLED' AND UPPER(sales.status) != 'RETURNED' THEN sales.id END) as bills,
			COALESCE(SUM(CASE WHEN UPPER(sales.status) != 'CANCELLED' AND UPPER(sales.status) != 'RETURNED' THEN sales.grand_total ELSE 0 END), 0) as gross_sales,
			COALESCE(SUM(CASE WHEN UPPER(sales.status) != 'CANCELLED' AND UPPER(sales.status) != 'RETURNED' THEN sales.discount ELSE 0 END), 0) as discounts,
			COALESCE(SUM(CASE WHEN UPPER(sales.status) = 'RETURNED' THEN sales.grand_total ELSE 0 END), 0) as returns,
			COALESCE(SUM(CASE WHEN UPPER(sales.status) != 'CANCELLED' AND UPPER(sales.status) != 'RETURNED' THEN sales.tax_amount ELSE 0 END), 0) as tax_amount,
			COALESCE(SUM(CASE WHEN UPPER(sales.status) != 'CANCELLED' AND UPPER(sales.status) != 'RETURNED' THEN sales.subtotal ELSE 0 END), 0) as net_sales
		`).
		Where("sales.tenant_id = ? AND sales.sale_date >= ? AND sales.sale_date <= ? AND UPPER(sales.status) != 'CANCELLED'", tenantID, start, end).
		Group("DATE(sales.sale_date)").
		Order("date ASC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	for i := range rows {
		var units int64
		_ = r.db.Table("sale_items").
			Joins("JOIN sales ON sale_items.sale_id = sales.id").
			Where("sale_items.tenant_id = ? AND DATE(sales.sale_date) = ? AND UPPER(sales.status) != 'CANCELLED' AND UPPER(sales.status) != 'RETURNED'", tenantID, rows[i].Date).
			Select("COALESCE(SUM(sale_items.quantity), 0)").
			Scan(&units)
		rows[i].UnitsSold = units

		// Net sales strictly reconciles with Gross - Discounts - Returns
		calculatedNet := rows[i].GrossSales.Sub(rows[i].Discounts).Sub(rows[i].Returns)
		if calculatedNet.IsNegative() {
			calculatedNet = decimal.Zero
		}
		rows[i].NetSales = calculatedNet
	}

	return rows, nil
}

type ProductSalesDeepRow struct {
	ProductID     uuid.UUID       `gorm:"column:product_id"`
	VariantID     uuid.UUID       `gorm:"column:variant_id"`
	ProductName   string          `gorm:"column:product_name"`
	ProductCode   string          `gorm:"column:product_code"`
	SKU           string          `gorm:"column:sku"`
	Barcode       string          `gorm:"column:barcode"`
	BrandName     string          `gorm:"column:brand_name"`
	CategoryName  string          `gorm:"column:category_name"`
	Size          string          `gorm:"column:size"`
	Colour        string          `gorm:"column:colour"`
	PurchasePrice decimal.Decimal `gorm:"column:purchase_price"`
	UnitsSold     int64           `gorm:"column:units_sold"`
	GrossSales    decimal.Decimal `gorm:"column:gross_sales"`
	Discounts     decimal.Decimal `gorm:"column:discounts"`
	NetSales      decimal.Decimal `gorm:"column:net_sales"`
	CostOfGoods   decimal.Decimal `gorm:"column:cost_of_goods"`
	CurrentStock  int             `gorm:"column:current_stock"`
}

// GetDeepProductSales returns detailed metrics including catalogue products that had 0 sales
func (r *AnalyticsRepository) GetDeepProductSales(tenantID uuid.UUID, start, end time.Time) ([]ProductSalesDeepRow, error) {
	var rows []ProductSalesDeepRow
	err := r.db.Table("product_variants").
		Joins("JOIN products ON product_variants.product_id = products.id").
		Joins("LEFT JOIN brands ON products.brand_id = brands.id").
		Joins("LEFT JOIN categories ON products.category_id = categories.id").
		Select(`
			products.id as product_id,
			product_variants.id as variant_id,
			products.name as product_name,
			products.product_code,
			product_variants.sku,
			product_variants.barcode,
			COALESCE(brands.name, 'Generic') as brand_name,
			COALESCE(categories.name, 'Footwear') as category_name,
			product_variants.size,
			product_variants.colour,
			product_variants.purchase_price,
			product_variants.current_stock
		`).
		Where("product_variants.tenant_id = ? AND products.status != 'archived'", tenantID).
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	for i := range rows {
		type aggItem struct {
			Units       int64           `gorm:"column:units"`
			Gross       decimal.Decimal `gorm:"column:gross"`
			Discounts   decimal.Decimal `gorm:"column:discounts"`
			LineTotal   decimal.Decimal `gorm:"column:line_total"`
			CostOfGoods decimal.Decimal `gorm:"column:cogs"`
		}
		var agg aggItem
		_ = r.db.Table("sale_items").
			Joins("JOIN sales ON sale_items.sale_id = sales.id").
			Select(`
				COALESCE(SUM(sale_items.quantity), 0) as units,
				COALESCE(SUM(sale_items.unit_price * sale_items.quantity), 0) as gross,
				COALESCE(SUM(sale_items.discount), 0) as discounts,
				COALESCE(SUM(sale_items.line_total), 0) as line_total,
				COALESCE(SUM(sale_items.quantity), 0) as cogs
			`).
			Where("sale_items.tenant_id = ? AND (sale_items.variant_id = ? OR (sale_items.barcode != '' AND sale_items.barcode = ?)) AND sales.sale_date >= ? AND sales.sale_date <= ? AND sales.status != 'CANCELLED'",
				tenantID, rows[i].VariantID, rows[i].Barcode, start, end).
			Scan(&agg)

		rows[i].UnitsSold = agg.Units
		rows[i].GrossSales = agg.Gross
		rows[i].Discounts = agg.Discounts
		rows[i].NetSales = agg.LineTotal
		rows[i].CostOfGoods = rows[i].PurchasePrice.Mul(decimal.NewFromInt(agg.Units))
	}

	return rows, nil
}

// GetRawTransactions returns sales with line items for a date range
func (r *AnalyticsRepository) GetRawTransactions(tenantID uuid.UUID, start, end time.Time) ([]models.Sale, error) {
	var sales []models.Sale
	err := r.db.Where("tenant_id = ? AND sale_date >= ? AND sale_date <= ?", tenantID, start, end).
		Preload("Items").
		Order("sale_date DESC").
		Limit(1000).
		Find(&sales).Error
	return sales, err
}
