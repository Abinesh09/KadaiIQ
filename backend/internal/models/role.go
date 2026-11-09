package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleScope string

const (
	RoleScopePlatform RoleScope = "platform"
	RoleScopeTenant   RoleScope = "tenant"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	Scope       RoleScope `gorm:"type:varchar(30);default:'tenant';not null;index" json:"scope"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`

	// Relationships
	Permissions []RolePermission `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"permissions,omitempty"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type MembershipStatus string

const (
	MembershipStatusActive    MembershipStatus = "active"
	MembershipStatusSuspended MembershipStatus = "suspended"
	MembershipStatusInvited   MembershipStatus = "invited"
)

type TenantMembership struct {
	ID               uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID         uuid.UUID        `gorm:"type:uuid;not null;index;index:idx_tenant_user,unique" json:"tenant_id"`
	UserID           uuid.UUID        `gorm:"type:uuid;not null;index;index:idx_tenant_user,unique" json:"user_id"`
	RoleID           uuid.UUID        `gorm:"type:uuid;not null;index" json:"role_id"`
	MembershipStatus MembershipStatus `gorm:"type:varchar(30);default:'active';not null;index" json:"membership_status"`
	CreatedAt        time.Time        `gorm:"not null" json:"created_at"`

	// Relationships
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role   *Role   `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (m *TenantMembership) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
