package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Module      string    `gorm:"type:varchar(50);not null;index" json:"module"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type RolePermission struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	RoleID       uuid.UUID   `gorm:"type:uuid;not null;index;index:idx_role_perm,unique" json:"role_id"`
	PermissionID uuid.UUID   `gorm:"type:uuid;not null;index;index:idx_role_perm,unique" json:"permission_id"`
	Role         *Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission   *Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

func (rp *RolePermission) BeforeCreate(tx *gorm.DB) error {
	if rp.ID == uuid.Nil {
		rp.ID = uuid.New()
	}
	return nil
}
