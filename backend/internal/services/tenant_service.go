package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/Abinesh09/KadaiIQ/backend/internal/repositories"
	"github.com/Abinesh09/KadaiIQ/backend/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantService struct {
	tenantRepo *repositories.TenantRepository
	authRepo   *repositories.AuthRepository
}

func NewTenantService(
	tenantRepo *repositories.TenantRepository,
	authRepo *repositories.AuthRepository,
) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
		authRepo:   authRepo,
	}
}

func (s *TenantService) ListTenants() ([]dto.TenantResponse, error) {
	tenants, err := s.tenantRepo.GetAllTenants()
	if err != nil {
		return nil, err
	}

	res := make([]dto.TenantResponse, 0, len(tenants))
	for _, t := range tenants {
		planName := "None"
		if t.Plan != nil {
			planName = t.Plan.Name
		}
		res = append(res, dto.TenantResponse{
			ID:           t.ID,
			Name:         t.Name,
			Slug:         t.Slug,
			BusinessName: t.BusinessName,
			BusinessType: t.BusinessType,
			GSTIN:        t.GSTIN,
			Email:        t.Email,
			Phone:        t.Phone,
			Address:      t.Address,
			Status:       string(t.Status),
			PlanName:     planName,
			MemberCount:  int64(len(t.Memberships)),
			CreatedAt:    t.CreatedAt,
		})
	}
	return res, nil
}

func (s *TenantService) CreateTenantBySuperAdmin(
	req dto.CreateTenantAdminRequest,
	actorID uuid.UUID,
) (*dto.TenantResponse, error) {
	hashedPassword, err := utils.HashPassword(req.OwnerPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now().UTC()
	tenant := models.Tenant{
		ID:           uuid.New(),
		Name:         req.Name,
		Slug:         req.Slug,
		BusinessName: req.BusinessName,
		BusinessType: req.BusinessType,
		GSTIN:        req.GSTIN,
		Email:        req.Email,
		Phone:        req.Phone,
		Address:      req.Address,
		PlanID:       req.PlanID,
		Status:       models.TenantStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	owner := models.User{
		ID:           uuid.New(),
		Name:         req.OwnerName,
		Email:        req.OwnerEmail,
		PasswordHash: hashedPassword,
		Status:       models.UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.tenantRepo.OnboardTenantWithTx(&tenant, &owner, actorID); err != nil {
		return nil, err
	}

	return &dto.TenantResponse{
		ID:           tenant.ID,
		Name:         tenant.Name,
		Slug:         tenant.Slug,
		BusinessName: tenant.BusinessName,
		BusinessType: tenant.BusinessType,
		GSTIN:        tenant.GSTIN,
		Email:        tenant.Email,
		Phone:        tenant.Phone,
		Address:      tenant.Address,
		Status:       string(tenant.Status),
		MemberCount:  1,
		CreatedAt:    tenant.CreatedAt,
	}, nil
}

func (s *TenantService) UpdateTenant(id uuid.UUID, req dto.UpdateTenantRequest) (*dto.TenantResponse, error) {
	tenant, err := s.tenantRepo.GetTenantByID(id)
	if err != nil {
		return nil, errors.New("tenant not found")
	}

	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.BusinessName != "" {
		tenant.BusinessName = req.BusinessName
	}
	if req.BusinessType != "" {
		tenant.BusinessType = req.BusinessType
	}
	if req.GSTIN != "" {
		tenant.GSTIN = req.GSTIN
	}
	if req.Email != "" {
		tenant.Email = req.Email
	}
	if req.Phone != "" {
		tenant.Phone = req.Phone
	}
	if req.Address != "" {
		tenant.Address = req.Address
	}
	if req.Status != "" {
		tenant.Status = models.TenantStatus(req.Status)
	}
	if req.PlanID != nil {
		tenant.PlanID = req.PlanID
	}
	tenant.UpdatedAt = time.Now().UTC()

	if err := s.tenantRepo.UpdateTenant(tenant); err != nil {
		return nil, err
	}

	return &dto.TenantResponse{
		ID:           tenant.ID,
		Name:         tenant.Name,
		Slug:         tenant.Slug,
		BusinessName: tenant.BusinessName,
		BusinessType: tenant.BusinessType,
		GSTIN:        tenant.GSTIN,
		Email:        tenant.Email,
		Phone:        tenant.Phone,
		Address:      tenant.Address,
		Status:       string(tenant.Status),
		CreatedAt:    tenant.CreatedAt,
	}, nil
}

func (s *TenantService) SetStatus(id uuid.UUID, status string) error {
	return s.tenantRepo.SetTenantStatus(id, models.TenantStatus(status))
}

func (s *TenantService) ToggleFeature(tenantID uuid.UUID, req dto.TenantFeatureToggleRequest) error {
	return s.tenantRepo.ToggleFeature(tenantID, req.FeatureKey, req.Enabled)
}

func (s *TenantService) GetMembers(tenantID uuid.UUID) ([]dto.MemberResponse, error) {
	members, err := s.tenantRepo.GetTenantMembers(tenantID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.MemberResponse, 0, len(members))
	for _, m := range members {
		name := ""
		email := ""
		roleName := ""
		if m.User != nil {
			name = m.User.Name
			email = m.User.Email
		}
		if m.Role != nil {
			roleName = m.Role.Name
		}
		res = append(res, dto.MemberResponse{
			ID:        m.ID,
			UserID:    m.UserID,
			Name:      name,
			Email:     email,
			RoleID:    m.RoleID,
			RoleName:  roleName,
			Status:    string(m.MembershipStatus),
			CreatedAt: m.CreatedAt,
		})
	}
	return res, nil
}

func (s *TenantService) InviteMember(tenantID uuid.UUID, req dto.InviteUserRequest) (*dto.MemberResponse, error) {
	var user models.User
	err := s.tenantRepo.DB().Where("LOWER(email) = LOWER(?)", req.Email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		// Create invited placeholder user
		tempPassword, _ := utils.GenerateRandomToken(16)
		hashed, _ := utils.HashPassword(tempPassword)
		now := time.Now().UTC()
		user = models.User{
			ID:           uuid.New(),
			Name:         req.Name,
			Email:        req.Email,
			PasswordHash: hashed,
			Status:       models.UserStatusActive,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if err := s.authRepo.CreateUser(&user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else if err != nil {
		return nil, err
	}

	membership, err := s.tenantRepo.AddTenantMember(tenantID, user.ID, req.RoleID)
	if err != nil {
		return nil, err
	}

	var role models.Role
	s.tenantRepo.DB().First(&role, "id = ?", req.RoleID)

	return &dto.MemberResponse{
		ID:        membership.ID,
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		RoleID:    req.RoleID,
		RoleName:  role.Name,
		Status:    string(membership.MembershipStatus),
		CreatedAt: membership.CreatedAt,
	}, nil
}

func (s *TenantService) UpdateMemberRole(tenantID, memberID, roleID uuid.UUID) error {
	return s.tenantRepo.UpdateMemberRole(tenantID, memberID, roleID)
}
