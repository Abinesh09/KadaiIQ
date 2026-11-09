package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) DB() *gorm.DB {
	return r.db
}

// Scoped returns GORM query scoped strictly to tenant
func (r *ProductRepository) Scoped(tenantID uuid.UUID) *gorm.DB {
	return r.db.Scopes(TenantScope(tenantID))
}

// FindProducts with search, filters and pagination
func (r *ProductRepository) FindProducts(
	tenantID uuid.UUID,
	search string,
	brandID *uuid.UUID,
	categoryID *uuid.UUID,
	gender string,
	status string,
	page int,
	pageSize int,
) ([]models.Product, int64, error) {
	query := r.Scoped(tenantID).Model(&models.Product{})

	if search != "" {
		likePattern := "%" + search + "%"
		query = query.Where("(LOWER(name) LIKE LOWER(?) OR LOWER(product_code) LIKE LOWER(?) OR LOWER(barcode) LIKE LOWER(?))",
			likePattern, likePattern, likePattern)
	}
	if brandID != nil && *brandID != uuid.Nil {
		query = query.Where("brand_id = ?", *brandID)
	}
	if categoryID != nil && *categoryID != uuid.Nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		// By default do not show archived
		query = query.Where("status != ?", models.ProductStatusArchived)
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var products []models.Product
	err := query.
		Preload("Brand").
		Preload("Category").
		Preload("Variants").
		Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error

	return products, totalCount, err
}

func (r *ProductRepository) GetProductByID(tenantID, productID uuid.UUID) (*models.Product, error) {
	var prod models.Product
	err := r.Scoped(tenantID).
		Preload("Brand").
		Preload("Category").
		Preload("Variants").
		First(&prod, "id = ?", productID).Error
	if err != nil {
		return nil, err
	}
	return &prod, nil
}

func (r *ProductRepository) FindVariantByBarcode(tenantID uuid.UUID, barcode string) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.Scoped(tenantID).
		Preload("Product").
		First(&variant, "barcode = ?", barcode).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *ProductRepository) FindVariantByID(tenantID, variantID uuid.UUID) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.Scoped(tenantID).
		Preload("Product").
		First(&variant, "id = ?", variantID).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

// CreateProductWithVariantsTx handles transactional creation of product, variants, and opening stock movements
func (r *ProductRepository) CreateProductWithVariantsTx(
	product *models.Product,
	variants []models.ProductVariant,
	openingQuantities []int,
	userID uuid.UUID,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Verify Unique Product Code in tenant
		var codeCount int64
		if err := tx.Model(&models.Product{}).
			Where("tenant_id = ? AND product_code = ?", product.TenantID, product.ProductCode).
			Count(&codeCount).Error; err != nil {
			return err
		}
		if codeCount > 0 {
			return fmt.Errorf("product code '%s' already exists in this organisation", product.ProductCode)
		}

		// 2. Insert Product
		if err := tx.Create(product).Error; err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}

		now := time.Now().UTC()

		// 3. Insert Variants & Create Opening Stock Movements
		for i := range variants {
			v := &variants[i]
			v.ProductID = product.ID
			v.TenantID = product.TenantID

			// Check unique barcode in tenant
			var barcodeCount int64
			if err := tx.Model(&models.ProductVariant{}).
				Where("tenant_id = ? AND barcode = ?", v.TenantID, v.Barcode).
				Count(&barcodeCount).Error; err != nil {
				return err
			}
			if barcodeCount > 0 {
				return fmt.Errorf("barcode '%s' already exists for another variant in this organisation", v.Barcode)
			}

			// Check unique SKU in tenant
			var skuCount int64
			if err := tx.Model(&models.ProductVariant{}).
				Where("tenant_id = ? AND sku = ?", v.TenantID, v.SKU).
				Count(&skuCount).Error; err != nil {
				return err
			}
			if skuCount > 0 {
				return fmt.Errorf("SKU '%s' already exists for another variant in this organisation", v.SKU)
			}

			openingQty := 0
			if i < len(openingQuantities) {
				openingQty = openingQuantities[i]
			}
			v.CurrentStock = openingQty

			if err := tx.Create(v).Error; err != nil {
				return fmt.Errorf("failed to create variant: %w", err)
			}

			// If opening stock > 0, record mandatory stock movement
			if openingQty > 0 {
				movement := models.StockMovement{
					ID:              uuid.New(),
					TenantID:        v.TenantID,
					VariantID:       v.ID,
					MovementType:    models.MovementOpening,
					Quantity:        openingQty,
					PreviousStock:   0,
					NewStock:        openingQty,
					UnitCost:        v.PurchasePrice,
					ReferenceNumber: "OPENING_STOCK",
					Notes:           "Initial opening stock on product creation",
					CreatedBy:       userID,
					CreatedAt:       now,
				}
				if err := tx.Create(&movement).Error; err != nil {
					return fmt.Errorf("failed to record opening stock movement: %w", err)
				}
			}
		}

		return nil
	})
}

// RecordStockMovementTx handles transactional stock adjustments and audits
func (r *ProductRepository) RecordStockMovementTx(
	tenantID uuid.UUID,
	variantID uuid.UUID,
	movType models.MovementType,
	qty int,
	unitCost decimal.Decimal,
	reference string,
	notes string,
	userID uuid.UUID,
) (*models.StockMovement, error) {
	if qty <= 0 {
		return nil, errors.New("movement quantity must be greater than zero")
	}

	var movement models.StockMovement

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var variant models.ProductVariant
		// Lock variant row for update to prevent concurrent race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND tenant_id = ?", variantID, tenantID).
			First(&variant).Error; err != nil {
			return fmt.Errorf("variant not found: %w", err)
		}

		prevStock := variant.CurrentStock
		var newStock int

		switch movType {
		case models.MovementPurchase, models.MovementCustomerReturn:
			newStock = prevStock + qty
		case models.MovementSale:
			newStock = prevStock - qty
			if newStock < 0 {
				return fmt.Errorf("insufficient stock: available %d, requested %d", prevStock, qty)
			}
		case models.MovementManualAdjustment:
			// For manual adjustments, notes determine direction or qty signed logic
			newStock = prevStock + qty
		case models.MovementTransfer:
			newStock = prevStock - qty
		default:
			return fmt.Errorf("unsupported movement type: %s", movType)
		}

		// Update variant current_stock
		if err := tx.Model(&models.ProductVariant{}).
			Where("id = ? AND tenant_id = ?", variantID, tenantID).
			Update("current_stock", newStock).Error; err != nil {
			return fmt.Errorf("failed to update current stock: %w", err)
		}

		cost := unitCost
		if cost.IsZero() {
			cost = variant.PurchasePrice
		}

		movement = models.StockMovement{
			ID:              uuid.New(),
			TenantID:        tenantID,
			VariantID:       variantID,
			MovementType:    movType,
			Quantity:        qty,
			PreviousStock:   prevStock,
			NewStock:        newStock,
			UnitCost:        cost,
			ReferenceNumber: reference,
			Notes:           notes,
			CreatedBy:       userID,
			CreatedAt:       time.Now().UTC(),
		}

		if err := tx.Create(&movement).Error; err != nil {
			return fmt.Errorf("failed to write stock movement ledger: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &movement, nil
}

func (r *ProductRepository) GetStockMovements(
	tenantID uuid.UUID,
	variantID *uuid.UUID,
	page int,
	pageSize int,
) ([]models.StockMovement, int64, error) {
	query := r.Scoped(tenantID).Model(&models.StockMovement{})
	if variantID != nil && *variantID != uuid.Nil {
		query = query.Where("variant_id = ?", *variantID)
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var movements []models.StockMovement
	err := query.
		Preload("Variant").
		Preload("Variant.Product").
		Preload("User").
		Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&movements).Error

	return movements, totalCount, err
}

type InventorySummary struct {
	TotalProducts   int64
	TotalVariants   int64
	TotalStockUnits int64
	LowStockCount   int64
	OutOfStockCount int64
	StockValuation  decimal.Decimal
}

func (r *ProductRepository) GetInventorySummary(tenantID uuid.UUID) (*InventorySummary, error) {
	var totalProds int64
	r.Scoped(tenantID).Model(&models.Product{}).Where("status != ?", models.ProductStatusArchived).Count(&totalProds)

	var totalVars int64
	r.Scoped(tenantID).Model(&models.ProductVariant{}).Count(&totalVars)

	type StockAgg struct {
		TotalUnits int64           `gorm:"column:total_units"`
		Valuation  decimal.Decimal `gorm:"column:valuation"`
	}
	var agg StockAgg
	r.Scoped(tenantID).Model(&models.ProductVariant{}).
		Select("COALESCE(SUM(current_stock), 0) as total_units, COALESCE(SUM(current_stock * purchase_price), 0) as valuation").
		Scan(&agg)

	var lowStockCount int64
	r.Scoped(tenantID).Model(&models.ProductVariant{}).
		Where("current_stock > 0 AND current_stock <= min_stock_threshold").
		Count(&lowStockCount)

	var outOfStockCount int64
	r.Scoped(tenantID).Model(&models.ProductVariant{}).
		Where("current_stock <= 0").
		Count(&outOfStockCount)

	return &InventorySummary{
		TotalProducts:   totalProds,
		TotalVariants:   totalVars,
		TotalStockUnits: agg.TotalUnits,
		LowStockCount:   lowStockCount,
		OutOfStockCount: outOfStockCount,
		StockValuation:  agg.Valuation,
	}, nil
}

func (r *ProductRepository) ArchiveProduct(tenantID, productID uuid.UUID) error {
	res := r.Scoped(tenantID).Model(&models.Product{}).
		Where("id = ?", productID).
		Update("status", models.ProductStatusArchived)
	if res.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return res.Error
}

func (r *ProductRepository) GetBrands(tenantID uuid.UUID) ([]models.Brand, error) {
	var brands []models.Brand
	err := r.Scoped(tenantID).Order("name asc").Find(&brands).Error
	return brands, err
}

func (r *ProductRepository) FindOrCreateBrand(tenantID uuid.UUID, name string) (*models.Brand, error) {
	var brand models.Brand
	err := r.Scoped(tenantID).Where("LOWER(name) = LOWER(?)", name).First(&brand).Error
	if err == gorm.ErrRecordNotFound {
		brand = models.Brand{
			ID:        uuid.New(),
			TenantID:  tenantID,
			Name:      name,
			CreatedAt: time.Now().UTC(),
		}
		if err := r.db.Create(&brand).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &brand, nil
}

func (r *ProductRepository) GetCategories(tenantID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.Scoped(tenantID).Order("name asc").Find(&categories).Error
	return categories, err
}

func (r *ProductRepository) FindOrCreateCategory(tenantID uuid.UUID, name string) (*models.Category, error) {
	var category models.Category
	err := r.Scoped(tenantID).Where("LOWER(name) = LOWER(?)", name).First(&category).Error
	if err == gorm.ErrRecordNotFound {
		category = models.Category{
			ID:        uuid.New(),
			TenantID:  tenantID,
			Name:      name,
			CreatedAt: time.Now().UTC(),
		}
		if err := r.db.Create(&category).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &category, nil
}
