package dto

import "github.com/google/uuid"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TenantMembershipInfo struct {
	TenantID     uuid.UUID `json:"tenant_id"`
	TenantName   string    `json:"tenant_name"`
	TenantSlug   string    `json:"tenant_slug"`
	RoleID       uuid.UUID `json:"role_id"`
	RoleName     string    `json:"role_name"`
	Status       string    `json:"status"`
	IsSuperAdmin bool      `json:"is_super_admin"`
}

type UserProfileResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type AuthTokenResponse struct {
	AccessToken      string                 `json:"access_token"`
	RefreshToken     string                 `json:"refresh_token"`
	ExpiresIn        int64                  `json:"expires_in"` // seconds
	User             UserProfileResponse    `json:"user"`
	ActiveTenant     *TenantMembershipInfo  `json:"active_tenant,omitempty"`
	AvailableTenants []TenantMembershipInfo `json:"available_tenants"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type SwitchTenantRequest struct {
	TenantID uuid.UUID `json:"tenant_id"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type RegisterTenantRequest struct {
	OrganizationName string `json:"organization_name"`
	Slug             string `json:"slug"`
	OwnerName        string `json:"owner_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Phone            string `json:"phone,omitempty"`
}
