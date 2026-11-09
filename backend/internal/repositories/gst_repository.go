package repositories

import (
	"errors"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type GSTRepository struct {
	db *gorm.DB
}

func NewGSTRepository(db *gorm.DB) *GSTRepository {
	return &GSTRepository{db: db}
}

func (r *GSTRepository) GetSettings(tenantID uuid.UUID) (*models.TenantGSTSetting, error) {
	var settings models.TenantGSTSetting
	err := r.db.Where("tenant_id = ?", tenantID).First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Auto initialize default GST setting for the tenant
			now := time.Now().UTC()
			defaultSetting := models.TenantGSTSetting{
				ID:                 uuid.New(),
				TenantID:           tenantID,
				GSTIN:              "33AAAAA0000A1Z5",
				LegalName:          "Retail Enterprise",
				TradeName:          "Retail Enterprise",
				StateCode:          "33",
				StateName:          "Tamil Nadu",
				RegistrationType:   models.GSTRegistrationRegular,
				DefaultFootwearHSN: "6403",
				DefaultGSTRate:     decimal.NewFromInt(12),
				EnableComposition:  false,
				CreatedAt:          now,
				UpdatedAt:          now,
			}
			if createErr := r.db.Create(&defaultSetting).Error; createErr != nil {
				return nil, createErr
			}
			return &defaultSetting, nil
		}
		return nil, err
	}
	return &settings, nil
}

func (r *GSTRepository) SaveSettings(settings *models.TenantGSTSetting) error {
	settings.UpdatedAt = time.Now().UTC()
	return r.db.Save(settings).Error
}

type RawTaxSlab struct {
	TaxRate      decimal.Decimal `gorm:"column:tax_rate"`
	TaxableValue decimal.Decimal `gorm:"column:taxable_value"`
	TaxAmount    decimal.Decimal `gorm:"column:tax_amount"`
	Invoices     int64           `gorm:"column:invoices"`
	TotalQty     int64           `gorm:"column:total_qty"`
}

func (r *GSTRepository) GetTaxSlabs(tenantID uuid.UUID, start, end time.Time) ([]RawTaxSlab, error) {
	var slabs []RawTaxSlab
	err := r.db.Table("sale_items").
		Joins("JOIN sales ON sale_items.sale_id = sales.id").
		Select(`
			sale_items.tax_rate,
			COALESCE(SUM((sale_items.unit_price * sale_items.quantity) - sale_items.discount), 0) as taxable_value,
			COALESCE(SUM(sale_items.tax_amount), 0) as tax_amount,
			COUNT(DISTINCT sales.id) as invoices,
			COALESCE(SUM(sale_items.quantity), 0) as total_qty
		`).
		Where("sale_items.tenant_id = ? AND sales.sale_date >= ? AND sales.sale_date <= ? AND sales.status != 'CANCELLED'", tenantID, start, end).
		Group("sale_items.tax_rate").
		Order("sale_items.tax_rate ASC").
		Scan(&slabs).Error

	return slabs, err
}

type RawHSNAggregate struct {
	HSNCode      string          `gorm:"column:hsn_code"`
	Description  string          `gorm:"column:description"`
	TotalQty     int64           `gorm:"column:total_qty"`
	TaxableValue decimal.Decimal `gorm:"column:taxable_value"`
	TaxAmount    decimal.Decimal `gorm:"column:tax_amount"`
}

func (r *GSTRepository) GetHSNAggregates(tenantID uuid.UUID, start, end time.Time, defaultHSN string) ([]RawHSNAggregate, error) {
	var hsns []RawHSNAggregate
	err := r.db.Table("sale_items").
		Joins("JOIN sales ON sale_items.sale_id = sales.id").
		Select(`
			? as hsn_code,
			'Footwear with outer soles of rubber, plastics or leather' as description,
			COALESCE(SUM(sale_items.quantity), 0) as total_qty,
			COALESCE(SUM((sale_items.unit_price * sale_items.quantity) - sale_items.discount), 0) as taxable_value,
			COALESCE(SUM(sale_items.tax_amount), 0) as tax_amount
		`, defaultHSN).
		Where("sale_items.tenant_id = ? AND sales.sale_date >= ? AND sales.sale_date <= ? AND sales.status != 'CANCELLED'", tenantID, start, end).
		Scan(&hsns).Error

	return hsns, err
}
