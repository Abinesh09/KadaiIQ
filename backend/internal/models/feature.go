package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantFeature struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID   uuid.UUID `gorm:"type:uuid;not null;index;index:idx_tenant_feature,unique" json:"tenant_id"`
	FeatureKey string    `gorm:"type:varchar(100);not null;index:idx_tenant_feature,unique" json:"feature_key"`
	Enabled    bool      `gorm:"default:true;not null" json:"enabled"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null" json:"updated_at"`

	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

func (tf *TenantFeature) BeforeCreate(tx *gorm.DB) error {
	if tf.ID == uuid.Nil {
		tf.ID = uuid.New()
	}
	return nil
}
