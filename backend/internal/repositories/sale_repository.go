package repositories

import (
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SaleRepository struct {
	db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

func (r *SaleRepository) DB() *gorm.DB {
	return r.db
}

func (r *SaleRepository) Scoped(tenantID uuid.UUID) *gorm.DB {
	return r.db.Scopes(TenantScope(tenantID))
}

func (r *SaleRepository) FindSales(
	tenantID uuid.UUID,
	search string,
	status string,
	dateFrom *time.Time,
	dateTo *time.Time,
	page int,
	pageSize int,
) ([]models.Sale, int64, error) {
	query := r.Scoped(tenantID).Model(&models.Sale{})

	if search != "" {
		likePattern := "%" + search + "%"
		query = query.Where("(bill_number LIKE ? OR source_bill_id LIKE ?)", likePattern, likePattern)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if dateFrom != nil {
		query = query.Where("sale_date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("sale_date <= ?", *dateTo)
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var sales []models.Sale
	err := query.
		Preload("Items").
		Order("sale_date desc").
		Offset(offset).
		Limit(pageSize).
		Find(&sales).Error

	return sales, totalCount, err
}

func (r *SaleRepository) GetSaleByID(tenantID, saleID uuid.UUID) (*models.Sale, error) {
	var sale models.Sale
	err := r.Scoped(tenantID).
		Preload("Items").
		Preload("Items.Variant").
		First(&sale, "id = ?", saleID).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *SaleRepository) FindSaleBySource(tenantID uuid.UUID, sourceSystem, sourceBillID string) (*models.Sale, error) {
	var sale models.Sale
	err := r.Scoped(tenantID).
		Where("source_system = ? AND source_bill_id = ?", sourceSystem, sourceBillID).
		First(&sale).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *SaleRepository) CreateImportJob(job *models.ImportJob) error {
	return r.db.Create(job).Error
}

func (r *SaleRepository) UpdateImportJob(job *models.ImportJob) error {
	return r.db.Save(job).Error
}

func (r *SaleRepository) GetImportJobs(tenantID uuid.UUID, limit int) ([]models.ImportJob, error) {
	var jobs []models.ImportJob
	err := r.Scoped(tenantID).
		Order("created_at desc").
		Limit(limit).
		Find(&jobs).Error
	return jobs, err
}
