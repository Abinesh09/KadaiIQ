package middleware

import (
	"errors"
	"net/http"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	TenantContextKey = "tenant_context"
	TenantIDHeader   = "X-Tenant-ID"
)

var (
	ErrMissingTenantContext = errors.New("missing tenant context in request")
	ErrInvalidTenantID      = errors.New("invalid or empty tenant ID")
	ErrCrossTenantForbidden = errors.New("cross-tenant access is strictly prohibited")
)

type TenantContext struct {
	TenantID     uuid.UUID `json:"tenant_id"`
	UserID       uuid.UUID `json:"user_id"`
	RoleID       uuid.UUID `json:"role_id"`
	RoleName     string    `json:"role_name"`
	IsSuperAdmin bool      `json:"is_super_admin"`
}

// RequireTenantContext middleware validates that a valid tenant context is attached to the request
func RequireTenantContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tc, ok := c.Locals(TenantContextKey).(*TenantContext)
		if !ok || tc == nil || tc.TenantID == uuid.Nil {
			return c.Status(http.StatusUnauthorized).JSON(
				dto.ErrorResponse("Tenant context required", "Authentication and tenant membership are required to access this resource"),
			)
		}
		return c.Next()
	}
}

// SetTenantContext stores the tenant context in Fiber locals
func SetTenantContext(c *fiber.Ctx, tc *TenantContext) {
	c.Locals(TenantContextKey, tc)
}

// GetTenantContext retrieves the tenant context from Fiber locals
func GetTenantContext(c *fiber.Ctx) (*TenantContext, error) {
	tc, ok := c.Locals(TenantContextKey).(*TenantContext)
	if !ok || tc == nil || tc.TenantID == uuid.Nil {
		return nil, ErrMissingTenantContext
	}
	return tc, nil
}
