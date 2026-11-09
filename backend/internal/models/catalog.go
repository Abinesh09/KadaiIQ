package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Brand struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index:idx_tenant_brand,unique" json:"tenant_id"`
	Name      string    `gorm:"type:varchar(100);not null;index:idx_tenant_brand,unique" json:"name"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	Tenant *Tenant `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"tenant,omitempty"`
}

func (b *Brand) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index:idx_tenant_category,unique" json:"tenant_id"`
	Name      string    `gorm:"type:varchar(100);not null;index:idx_tenant_category,unique" json:"name"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	Tenant *Tenant `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"tenant,omitempty"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
