package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type GSTRegistrationType string

const (
	GSTRegistrationRegular      GSTRegistrationType = "REGULAR"
	GSTRegistrationComposition  GSTRegistrationType = "COMPOSITION"
	GSTRegistrationUnregistered GSTRegistrationType = "UNREGISTERED"
)

type TenantGSTSetting struct {
	ID                 uuid.UUID           `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID           uuid.UUID           `gorm:"type:uuid;not null;uniqueIndex" json:"tenant_id"`
	GSTIN              string              `gorm:"type:varchar(20);not null" json:"gstin"`
	LegalName          string              `gorm:"type:varchar(255);not null" json:"legal_name"`
	TradeName          string              `gorm:"type:varchar(255)" json:"trade_name"`
	StateCode          string              `gorm:"type:varchar(5);not null" json:"state_code"` // e.g. "33" for Tamil Nadu
	StateName          string              `gorm:"type:varchar(100);not null" json:"state_name"`
	RegistrationType   GSTRegistrationType `gorm:"type:varchar(30);default:'REGULAR';not null" json:"registration_type"`
	DefaultFootwearHSN string              `gorm:"type:varchar(20);default:'6403';not null" json:"default_footwear_hsn"`
	DefaultGSTRate     decimal.Decimal     `gorm:"type:numeric(5,2);default:12.00;not null" json:"default_gst_rate"`
	EnableComposition  bool                `gorm:"default:false;not null" json:"enable_composition"`
	CreatedAt          time.Time           `gorm:"not null" json:"created_at"`
	UpdatedAt          time.Time           `gorm:"not null" json:"updated_at"`
}

func (s *TenantGSTSetting) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
