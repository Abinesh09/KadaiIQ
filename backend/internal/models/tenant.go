package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusInactive  TenantStatus = "inactive"
)

type Tenant struct {
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string       `gorm:"type:varchar(255);not null" json:"name"`
	Slug         string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	BusinessName string       `gorm:"type:varchar(255);not null" json:"business_name"`
	BusinessType string       `gorm:"type:varchar(100);default:'retail';not null" json:"business_type"`
	GSTIN        string       `gorm:"type:varchar(20)" json:"gstin,omitempty"`
	Email        string       `gorm:"type:varchar(255);not null" json:"email"`
	Phone        string       `gorm:"type:varchar(50)" json:"phone"`
	Address      string       `gorm:"type:text" json:"address"`
	Status       TenantStatus `gorm:"type:varchar(30);default:'active';not null;index" json:"status"`
	PlanID       *uuid.UUID   `gorm:"type:uuid;index" json:"plan_id,omitempty"`
	CreatedAt    time.Time    `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"not null" json:"updated_at"`

	// Relationships
	Plan        *SubscriptionPlan  `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	Memberships []TenantMembership `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"memberships,omitempty"`
	Features    []TenantFeature    `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"features,omitempty"`
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
