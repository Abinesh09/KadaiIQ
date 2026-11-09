package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type SaleStatus string

const (
	SaleStatusCompleted SaleStatus = "COMPLETED"
	SaleStatusCancelled SaleStatus = "CANCELLED"
	SaleStatusVoided    SaleStatus = "VOIDED"
	SaleStatusReturned  SaleStatus = "RETURNED"
)

type Sale struct {
	ID                uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID          uuid.UUID       `gorm:"type:uuid;not null;index:idx_tenant_source_bill,unique" json:"tenant_id"`
	SourceSystem      string          `gorm:"type:varchar(50);not null;default:'ZenBiller';index:idx_tenant_source_bill,unique" json:"source_system"`
	SourceBillID      string          `gorm:"type:varchar(100);not null;index:idx_tenant_source_bill,unique" json:"source_bill_id"`
	BillNumber        string          `gorm:"type:varchar(100);not null;index" json:"bill_number"`
	SaleDate          time.Time       `gorm:"not null;index" json:"sale_date"`
	Subtotal          decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"subtotal"`
	Discount          decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"discount"`
	TaxAmount         decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"tax_amount"`
	GrandTotal        decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"grand_total"`
	PaymentMethod     string          `gorm:"type:varchar(50);default:'CASH';not null" json:"payment_method"`
	Status            SaleStatus      `gorm:"type:varchar(30);default:'COMPLETED';not null;index" json:"status"`
	InventoryDeducted bool            `gorm:"default:false;not null" json:"inventory_deducted"`
	ImportJobID       *uuid.UUID      `gorm:"type:uuid;index" json:"import_job_id,omitempty"`
	CreatedAt         time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"not null" json:"updated_at"`

	// Relationships
	Items []SaleItem `gorm:"foreignKey:SaleID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

func (s *Sale) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type SaleItem struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID     uuid.UUID       `gorm:"type:uuid;not null;index" json:"tenant_id"`
	SaleID       uuid.UUID       `gorm:"type:uuid;not null;index" json:"sale_id"`
	SourceLineID string          `gorm:"type:varchar(100)" json:"source_line_id,omitempty"`
	VariantID    *uuid.UUID      `gorm:"type:uuid;index" json:"variant_id,omitempty"`
	ProductName  string          `gorm:"type:varchar(255);not null" json:"product_name"`
	SKU          string          `gorm:"type:varchar(100)" json:"sku,omitempty"`
	Barcode      string          `gorm:"type:varchar(100);index" json:"barcode,omitempty"`
	Quantity     int             `gorm:"type:int;not null;default:1" json:"quantity"`
	UnitPrice    decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"unit_price"`
	Discount     decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"discount"`
	TaxRate      decimal.Decimal `gorm:"type:numeric(5,2);not null;default:0" json:"tax_rate"`
	TaxAmount    decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"tax_amount"`
	LineTotal    decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"line_total"`
	CreatedAt    time.Time       `gorm:"not null" json:"created_at"`

	// Relationships
	Variant *ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
}

func (si *SaleItem) BeforeCreate(tx *gorm.DB) error {
	if si.ID == uuid.Nil {
		si.ID = uuid.New()
	}
	return nil
}

type ImportJobStatus string

const (
	ImportStatusPending    ImportJobStatus = "PENDING"
	ImportStatusProcessing ImportJobStatus = "PROCESSING"
	ImportStatusCompleted  ImportJobStatus = "COMPLETED"
	ImportStatusFailed     ImportJobStatus = "FAILED"
)

type ImportJob struct {
	ID               uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID         uuid.UUID       `gorm:"type:uuid;not null;index" json:"tenant_id"`
	SourceSystem     string          `gorm:"type:varchar(50);not null" json:"source_system"`
	FileName         string          `gorm:"type:varchar(255);not null" json:"file_name"`
	TotalRecords     int             `gorm:"default:0" json:"total_records"`
	ImportedRecords  int             `gorm:"default:0" json:"imported_records"`
	FailedRecords    int             `gorm:"default:0" json:"failed_records"`
	DuplicateRecords int             `gorm:"default:0" json:"duplicate_records"`
	Status           ImportJobStatus `gorm:"type:varchar(30);default:'PENDING';not null;index" json:"status"`
	ErrorDetails     string          `gorm:"type:text" json:"error_details,omitempty"`
	CreatedBy        uuid.UUID       `gorm:"type:uuid;not null;index" json:"created_by"`
	CreatedAt        time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"not null" json:"updated_at"`

	User *User `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
}

func (ij *ImportJob) BeforeCreate(tx *gorm.DB) error {
	if ij.ID == uuid.Nil {
		ij.ID = uuid.New()
	}
	return nil
}
