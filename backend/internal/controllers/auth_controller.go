package controllers

import (
	"net/http"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/middleware"
	"github.com/Abinesh09/KadaiIQ/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid request body", err.Error()))
	}
	if req.Email == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "Email and password are required"))
	}

	res, err := ctrl.authService.Login(req)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Authentication failed", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Login successful", res))
}

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req dto.RegisterTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid request body", err.Error()))
	}
	if req.OrganizationName == "" || req.Slug == "" || req.OwnerName == "" || req.Email == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "Organization name, slug, owner name, email, and password are required"))
	}

	res, err := ctrl.authService.RegisterTenant(req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Registration failed", err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(dto.SuccessResponse("Organisation onboarded successfully", res))
}

func (ctrl *AuthController) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid request body", err.Error()))
	}
	if req.RefreshToken == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "refresh_token is required"))
	}

	res, err := ctrl.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Refresh failed", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Token refreshed successfully", res))
}

func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	_ = c.BodyParser(&req)

	if err := ctrl.authService.Logout(req.RefreshToken); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Logout error", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Logged out successfully", nil))
}

func (ctrl *AuthController) SwitchTenant(c *fiber.Ctx) error {
	tc, ok := c.Locals(middleware.TenantContextKey).(*middleware.TenantContext)
	if !ok || tc == nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", "Session expired"))
	}

	var req dto.SwitchTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid request body", err.Error()))
	}

	res, err := ctrl.authService.SwitchTenant(tc.UserID, req.TenantID)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(dto.ErrorResponse("Tenant switch failed", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Tenant context switched", res))
}

func (ctrl *AuthController) ForgotPassword(c *fiber.Ctx) error {
	var req dto.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid request body", err.Error()))
	}
	if req.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "Email is required"))
	}

	token, err := ctrl.authService.ForgotPassword(req.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Error initiating reset", err.Error()))
	}

	// In development, return token in payload for testability
	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("If that account exists, a password reset token was generated", fiber.Map{
		"reset_token": token,
	}))
}

func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid request body", err.Error()))
	}
	if req.Token == "" || req.NewPassword == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "Token and new_password are required"))
	}

	if err := ctrl.authService.ResetPassword(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Password reset failed", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Password reset successfully. Please log in with your new password.", nil))
}

func (ctrl *AuthController) Me(c *fiber.Ctx) error {
	tc, ok := c.Locals(middleware.TenantContextKey).(*middleware.TenantContext)
	if !ok || tc == nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", "Session expired"))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Current user context", tc))
}
