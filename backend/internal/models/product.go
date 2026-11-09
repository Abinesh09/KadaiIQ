package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusArchived ProductStatus = "archived"
	ProductStatusInactive ProductStatus = "inactive"
)

type Product struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID    uuid.UUID     `gorm:"type:uuid;not null;index:idx_tenant_code,unique" json:"tenant_id"`
	Name        string        `gorm:"type:varchar(255);not null;index" json:"name"`
	ProductCode string        `gorm:"type:varchar(100);not null;index:idx_tenant_code,unique" json:"product_code"`
	SKU         string        `gorm:"type:varchar(100);index" json:"sku,omitempty"`
	Barcode     string        `gorm:"type:varchar(100);index" json:"barcode,omitempty"`
	BrandID     *uuid.UUID    `gorm:"type:uuid;index" json:"brand_id,omitempty"`
	CategoryID  *uuid.UUID    `gorm:"type:uuid;index" json:"category_id,omitempty"`
	Description string        `gorm:"type:text" json:"description,omitempty"`
	Gender      string        `gorm:"type:varchar(30);default:'Unisex';not null;index" json:"gender"`
	Status      ProductStatus `gorm:"type:varchar(30);default:'active';not null;index" json:"status"`
	CreatedAt   time.Time     `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"not null" json:"updated_at"`

	// Relationships
	Brand    *Brand           `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	Category *Category        `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"variants,omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type ProductVariant struct {
	ID                uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID          uuid.UUID       `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ProductID         uuid.UUID       `gorm:"type:uuid;not null;index" json:"product_id"`
	Size              string          `gorm:"type:varchar(50);not null" json:"size"`
	Colour            string          `gorm:"type:varchar(50);not null" json:"colour"`
	SKU               string          `gorm:"type:varchar(100);not null;index:idx_tenant_var_sku,unique" json:"sku"`
	Barcode           string          `gorm:"type:varchar(100);not null;index:idx_tenant_var_barcode,unique" json:"barcode"`
	PurchasePrice     decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"purchase_price"`
	SellingPrice      decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"selling_price"`
	GSTRate           decimal.Decimal `gorm:"type:numeric(5,2);not null;default:12.00" json:"gst_rate"`
	CurrentStock      int             `gorm:"type:int;not null;default:0;index" json:"current_stock"`
	MinStockThreshold int             `gorm:"type:int;not null;default:5" json:"min_stock_threshold"`
	CreatedAt         time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"not null" json:"updated_at"`

	// Relationships
	Product   *Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Movements []StockMovement `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE" json:"movements,omitempty"`
}

func (pv *ProductVariant) BeforeCreate(tx *gorm.DB) error {
	if pv.ID == uuid.Nil {
		pv.ID = uuid.New()
	}
	return nil
}

type MovementType string

const (
	MovementOpening          MovementType = "OPENING"
	MovementPurchase         MovementType = "PURCHASE"
	MovementSale             MovementType = "SALE"
	MovementCustomerReturn   MovementType = "CUSTOMER_RETURN"
	MovementManualAdjustment MovementType = "MANUAL_ADJUSTMENT"
	MovementTransfer         MovementType = "TRANSFER"
)

type StockMovement struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"tenant_id"`
	VariantID       uuid.UUID       `gorm:"type:uuid;not null;index" json:"variant_id"`
	MovementType    MovementType    `gorm:"type:varchar(30);not null;index" json:"movement_type"`
	Quantity        int             `gorm:"type:int;not null" json:"quantity"` // + for in, - for out
	PreviousStock   int             `gorm:"type:int;not null" json:"previous_stock"`
	NewStock        int             `gorm:"type:int;not null" json:"new_stock"`
	UnitCost        decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"unit_cost"`
	ReferenceNumber string          `gorm:"type:varchar(100);index" json:"reference_number,omitempty"`
	Notes           string          `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy       uuid.UUID       `gorm:"type:uuid;not null;index" json:"created_by"`
	CreatedAt       time.Time       `gorm:"not null;index" json:"created_at"`

	// Relationships
	Variant *ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
	User    *User           `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
}

func (sm *StockMovement) BeforeCreate(tx *gorm.DB) error {
	if sm.ID == uuid.Nil {
		sm.ID = uuid.New()
	}
	return nil
}
