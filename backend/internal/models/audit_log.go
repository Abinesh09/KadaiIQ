package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID   *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"` // Nullable for platform/super-admin actions
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Action     string     `gorm:"type:varchar(100);not null;index" json:"action"`
	EntityType string     `gorm:"type:varchar(100);not null;index" json:"entity_type"`
	EntityID   string     `gorm:"type:varchar(100);not null" json:"entity_id"`
	Details    string     `gorm:"type:text" json:"details,omitempty"`
	CreatedAt  time.Time  `gorm:"not null;index" json:"created_at"`

	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (al *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if al.ID == uuid.Nil {
		al.ID = uuid.New()
	}
	return nil
}
