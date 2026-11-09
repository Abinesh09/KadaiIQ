package controllers

import (
	"net/http"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/middleware"
	"github.com/Abinesh09/KadaiIQ/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TenantController struct {
	tenantService *services.TenantService
}

func NewTenantController(tenantService *services.TenantService) *TenantController {
	return &TenantController{tenantService: tenantService}
}

// Super Admin: List all tenants
func (ctrl *TenantController) ListTenants(c *fiber.Ctx) error {
	tenants, err := ctrl.tenantService.ListTenants()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch organisations", err.Error()))
	}
	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Organisations retrieved", tenants))
}

// Super Admin: Create new tenant + owner
func (ctrl *TenantController) CreateTenant(c *fiber.Ctx) error {
	var req dto.CreateTenantAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid request body", err.Error()))
	}

	tc, _ := middleware.GetTenantContext(c)
	actorID := uuid.Nil
	if tc != nil {
		actorID = tc.UserID
	}

	res, err := ctrl.tenantService.CreateTenantBySuperAdmin(req, actorID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Failed to create organisation", err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(dto.SuccessResponse("Organisation created successfully", res))
}

// Super Admin: Update tenant
func (ctrl *TenantController) UpdateTenant(c *fiber.Ctx) error {
	idStr := c.Params("id")
	tenantID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid tenant ID", err.Error()))
	}

	var req dto.UpdateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid request body", err.Error()))
	}

	res, err := ctrl.tenantService.UpdateTenant(tenantID, req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Update failed", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Organisation updated successfully", res))
}

// Super Admin: Activate / Suspend tenant
func (ctrl *TenantController) SetStatus(c *fiber.Ctx) error {
	idStr := c.Params("id")
	tenantID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid tenant ID", err.Error()))
	}

	type StatusReq struct {
		Status string `json:"status"`
	}
	var req StatusReq
	if err := c.BodyParser(&req); err != nil || req.Status == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "Status is required"))
	}

	if err := ctrl.tenantService.SetStatus(tenantID, req.Status); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Status update failed", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Organisation status updated to "+req.Status, nil))
}

// Super Admin: Toggle tenant feature
func (ctrl *TenantController) ToggleFeature(c *fiber.Ctx) error {
	idStr := c.Params("id")
	tenantID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid tenant ID", err.Error()))
	}

	var req dto.TenantFeatureToggleRequest
	if err := c.BodyParser(&req); err != nil || req.FeatureKey == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "feature_key and enabled status are required"))
	}

	if err := ctrl.tenantService.ToggleFeature(tenantID, req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Feature toggle failed", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Feature updated successfully", nil))
}

// Tenant Admin: List organisation members
func (ctrl *TenantController) GetMembers(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	members, err := ctrl.tenantService.GetMembers(tc.TenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse("Failed to fetch members", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Members retrieved", members))
}

// Tenant Admin: Invite member
func (ctrl *TenantController) InviteMember(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	var req dto.InviteUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid request body", err.Error()))
	}
	if req.Email == "" || req.Name == "" || req.RoleID == uuid.Nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "name, email, and role_id are required"))
	}

	res, err := ctrl.tenantService.InviteMember(tc.TenantID, req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invitation failed", err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(dto.SuccessResponse("Member invited successfully", res))
}

// Tenant Admin: Update member role
func (ctrl *TenantController) UpdateMemberRole(c *fiber.Ctx) error {
	tc, err := middleware.GetTenantContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse("Unauthorized", err.Error()))
	}

	memberIDStr := c.Params("id")
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Invalid member ID", err.Error()))
	}

	type RoleReq struct {
		RoleID uuid.UUID `json:"role_id"`
	}
	var req RoleReq
	if err := c.BodyParser(&req); err != nil || req.RoleID == uuid.Nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Validation error", "role_id is required"))
	}

	if err := ctrl.tenantService.UpdateMemberRole(tc.TenantID, memberID, req.RoleID); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse("Role update failed", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse("Member role updated", nil))
}
