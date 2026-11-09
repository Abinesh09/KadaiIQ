package repositories

import (
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetSalesSummaryReport(tenantID uuid.UUID, start, end time.Time) ([]dto.SalesSummaryRow, error) {
	var rows []dto.SalesSummaryRow
	err := r.db.Table("sales").
		Select(`
			sales.bill_number,
			sales.sale_date,
			COALESCE(sales.payment_method, 'Cash') as payment_method,
			(SELECT COUNT(sale_items.id) FROM sale_items WHERE sale_items.sale_id = sales.id) as item_count,
			sales.subtotal,
			sales.discount,
			sales.tax_amount,
			sales.grand_total,
			sales.status
		`).
		Where("sales.tenant_id = ? AND sales.sale_date >= ? AND sales.sale_date <= ?", tenantID, start, end).
		Order("sales.sale_date DESC").
		Scan(&rows).Error

	return rows, err
}

func (r *ReportRepository) GetStockValuationReport(tenantID uuid.UUID) ([]dto.StockValuationRow, error) {
	var rows []dto.StockValuationRow
	err := r.db.Table("product_variants").
		Joins("JOIN products ON product_variants.product_id = products.id").
		Select(`
			products.name as product_name,
			products.product_code,
			product_variants.size,
			product_variants.colour,
			product_variants.sku,
			product_variants.barcode,
			product_variants.current_stock,
			product_variants.purchase_price,
			product_variants.selling_price,
			(product_variants.current_stock * product_variants.purchase_price) as cost_valuation,
			(product_variants.current_stock * product_variants.selling_price) as retail_valuation,
			CASE
				WHEN product_variants.current_stock <= 0 THEN 'OUT_OF_STOCK'
				WHEN product_variants.current_stock <= product_variants.min_stock_threshold THEN 'LOW_STOCK'
				ELSE 'IN_STOCK'
			END as stock_status
		`).
		Where("product_variants.tenant_id = ? AND products.status != 'archived'", tenantID).
		Order("products.name ASC, product_variants.size ASC").
		Scan(&rows).Error

	return rows, err
}

func (r *ReportRepository) GetStockMovementReport(tenantID uuid.UUID, start, end time.Time) ([]dto.StockMovementReportRow, error) {
	var rows []dto.StockMovementReportRow
	err := r.db.Table("stock_movements").
		Joins("JOIN product_variants ON stock_movements.variant_id = product_variants.id").
		Joins("JOIN products ON product_variants.product_id = products.id").
		Joins("LEFT JOIN users ON stock_movements.created_by = users.id").
		Select(`
			stock_movements.created_at as date,
			products.name as product_name,
			product_variants.barcode,
			(product_variants.size || ' / ' || product_variants.colour) as variant_details,
			stock_movements.movement_type,
			stock_movements.quantity,
			stock_movements.previous_stock,
			stock_movements.new_stock,
			stock_movements.unit_cost,
			stock_movements.reference_number,
			stock_movements.notes,
			COALESCE(users.name, 'System') as created_by
		`).
		Where("stock_movements.tenant_id = ? AND stock_movements.created_at >= ? AND stock_movements.created_at <= ?", tenantID, start, end).
		Order("stock_movements.created_at DESC").
		Scan(&rows).Error

	return rows, err
}
