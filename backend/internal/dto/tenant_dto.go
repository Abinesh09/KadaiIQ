package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTenantAdminRequest struct {
	Name          string     `json:"name"`
	Slug          string     `json:"slug"`
	BusinessName  string     `json:"business_name"`
	BusinessType  string     `json:"business_type"`
	GSTIN         string     `json:"gstin,omitempty"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone,omitempty"`
	Address       string     `json:"address,omitempty"`
	PlanID        *uuid.UUID `json:"plan_id,omitempty"`
	OwnerName     string     `json:"owner_name"`
	OwnerEmail    string     `json:"owner_email"`
	OwnerPassword string     `json:"owner_password"`
}

type UpdateTenantRequest struct {
	Name         string     `json:"name"`
	BusinessName string     `json:"business_name"`
	BusinessType string     `json:"business_type"`
	GSTIN        string     `json:"gstin,omitempty"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone,omitempty"`
	Address      string     `json:"address,omitempty"`
	Status       string     `json:"status"`
	PlanID       *uuid.UUID `json:"plan_id,omitempty"`
}

type TenantResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	BusinessName string    `json:"business_name"`
	BusinessType string    `json:"business_type"`
	GSTIN        string    `json:"gstin,omitempty"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone,omitempty"`
	Address      string    `json:"address,omitempty"`
	Status       string    `json:"status"`
	PlanName     string    `json:"plan_name,omitempty"`
	MemberCount  int64     `json:"member_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type TenantFeatureToggleRequest struct {
	FeatureKey string `json:"feature_key"`
	Enabled    bool   `json:"enabled"`
}

type InviteUserRequest struct {
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	RoleID uuid.UUID `json:"role_id"`
}

type MemberResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	RoleID    uuid.UUID `json:"role_id"`
	RoleName  string    `json:"role_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
