package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	TenantContextKey = "tenant_context"
	ClaimsContextKey = "jwt_claims"
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
	Permissions  []string  `json:"permissions"`
}

// JWTMiddleware validates the Bearer token and attaches TenantContext to request
func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(
				dto.ErrorResponse("Authentication required", "Missing or malformed Authorization header"),
			)
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateAccessToken(tokenStr, secret)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(
				dto.ErrorResponse("Invalid or expired token", err.Error()),
			)
		}

		tc := &TenantContext{
			TenantID:     claims.TenantID,
			UserID:       claims.UserID,
			RoleID:       claims.RoleID,
			RoleName:     claims.RoleName,
			IsSuperAdmin: claims.IsSuperAdmin,
			Permissions:  claims.Permissions,
		}

		c.Locals(ClaimsContextKey, claims)
		c.Locals(TenantContextKey, tc)

		return c.Next()
	}
}

// RequireTenantContext verifies that the request is bound to a valid tenant
func RequireTenantContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tc, ok := c.Locals(TenantContextKey).(*TenantContext)
		if !ok || tc == nil || tc.TenantID == uuid.Nil {
			return c.Status(http.StatusUnauthorized).JSON(
				dto.ErrorResponse("Tenant context required", "Authentication and active organisation membership are required"),
			)
		}
		return c.Next()
	}
}

// RequireSuperAdmin verifies that the authenticated user is a platform Super Admin
func RequireSuperAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tc, ok := c.Locals(TenantContextKey).(*TenantContext)
		if !ok || tc == nil || !tc.IsSuperAdmin {
			return c.Status(http.StatusForbidden).JSON(
				dto.ErrorResponse("Forbidden", "Access denied: Super Admin privileges required"),
			)
		}
		return c.Next()
	}
}

// RequirePermission verifies that the user possesses a specific permission
func RequirePermission(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tc, ok := c.Locals(TenantContextKey).(*TenantContext)
		if !ok || tc == nil {
			return c.Status(http.StatusUnauthorized).JSON(
				dto.ErrorResponse("Unauthorized", "Tenant context missing"),
			)
		}

		if tc.IsSuperAdmin {
			return c.Next()
		}

		for _, p := range tc.Permissions {
			if p == "*" || p == requiredPermission {
				return c.Next()
			}
		}

		return c.Status(http.StatusForbidden).JSON(
			dto.ErrorResponse("Permission denied", "You lack the required permission: "+requiredPermission),
		)
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
