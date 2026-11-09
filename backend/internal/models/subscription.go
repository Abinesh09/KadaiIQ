package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionPlanStatus string

const (
	SubscriptionPlanStatusActive   SubscriptionPlanStatus = "active"
	SubscriptionPlanStatusInactive SubscriptionPlanStatus = "inactive"
)

type SubscriptionPlan struct {
	ID          uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string                 `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string                 `gorm:"type:text" json:"description"`
	Status      SubscriptionPlanStatus `gorm:"type:varchar(30);default:'active';not null;index" json:"status"`
	CreatedAt   time.Time              `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time              `gorm:"not null" json:"updated_at"`
}

func (sp *SubscriptionPlan) BeforeCreate(tx *gorm.DB) error {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}
