package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/config"
	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/Abinesh09/KadaiIQ/backend/internal/repositories"
	"github.com/Abinesh09/KadaiIQ/backend/internal/utils"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountInactive    = errors.New("your account is currently inactive")
	ErrTenantAccessDenied = errors.New("you do not have access to the requested organisation")
)

type AuthService struct {
	authRepo   *repositories.AuthRepository
	tenantRepo *repositories.TenantRepository
	cfg        *config.Config
}

func NewAuthService(
	authRepo *repositories.AuthRepository,
	tenantRepo *repositories.TenantRepository,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		authRepo:   authRepo,
		tenantRepo: tenantRepo,
		cfg:        cfg,
	}
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthTokenResponse, error) {
	user, err := s.authRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if user.Status != models.UserStatusActive {
		return nil, ErrAccountInactive
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return s.buildAuthResponseForUser(user, uuid.Nil)
}

func (s *AuthService) RefreshToken(refreshTokenStr string) (*dto.AuthTokenResponse, error) {
	tokenHash := utils.HashToken(refreshTokenStr)
	storedToken, err := s.authRepo.FindRefreshToken(tokenHash)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	user, err := s.authRepo.FindUserByID(storedToken.UserID)
	if err != nil || user.Status != models.UserStatusActive {
		return nil, errors.New("user account is inactive or not found")
	}

	// Token rotation: revoke old token
	_ = s.authRepo.RevokeRefreshToken(tokenHash)

	return s.buildAuthResponseForUser(user, uuid.Nil)
}

func (s *AuthService) SwitchTenant(userID uuid.UUID, targetTenantID uuid.UUID) (*dto.AuthTokenResponse, error) {
	user, err := s.authRepo.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.buildAuthResponseForUser(user, targetTenantID)
}

func (s *AuthService) Logout(refreshTokenStr string) error {
	if refreshTokenStr == "" {
		return nil
	}
	tokenHash := utils.HashToken(refreshTokenStr)
	return s.authRepo.RevokeRefreshToken(tokenHash)
}

func (s *AuthService) ForgotPassword(email string) (string, error) {
	user, err := s.authRepo.FindUserByEmail(email)
	if err != nil {
		// Silent return to avoid email enumeration
		return "", nil
	}

	rawToken, err := utils.GenerateRandomToken(32)
	if err != nil {
		return "", err
	}

	tokenHash := utils.HashToken(rawToken)
	resetToken := models.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(1 * time.Hour), // 1 hour single-use
		CreatedAt: time.Now().UTC(),
	}

	if err := s.authRepo.StorePasswordResetToken(&resetToken); err != nil {
		return "", err
	}

	return rawToken, nil
}

func (s *AuthService) ResetPassword(req dto.ResetPasswordRequest) error {
	tokenHash := utils.HashToken(req.Token)
	resetToken, err := s.authRepo.FindValidPasswordResetToken(tokenHash)
	if err != nil {
		return errors.New("password reset token is invalid, expired, or has already been used")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	if err := s.authRepo.UpdateUserPassword(resetToken.UserID, hashedPassword); err != nil {
		return err
	}

	// Mark token used (enforce single-use)
	if err := s.authRepo.MarkPasswordResetTokenUsed(resetToken.ID); err != nil {
		return err
	}

	// Revoke all existing sessions
	_ = s.authRepo.RevokeAllUserTokens(resetToken.UserID)

	return nil
}

func (s *AuthService) RegisterTenant(req dto.RegisterTenantRequest) (*dto.AuthTokenResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	tenant := models.Tenant{
		ID:           uuid.New(),
		Name:         req.OrganizationName,
		Slug:         req.Slug,
		BusinessName: req.OrganizationName,
		Email:        req.Email,
		Phone:        req.Phone,
		Status:       models.TenantStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	owner := models.User{
		ID:           uuid.New(),
		Name:         req.OwnerName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Status:       models.UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.tenantRepo.OnboardTenantWithTx(&tenant, &owner, uuid.Nil); err != nil {
		return nil, fmt.Errorf("registration failed: %w", err)
	}

	return s.buildAuthResponseForUser(&owner, tenant.ID)
}

func (s *AuthService) buildAuthResponseForUser(user *models.User, requestedTenantID uuid.UUID) (*dto.AuthTokenResponse, error) {
	isSuperAdmin, superAdminRoleID, _ := s.authRepo.CheckUserIsSuperAdmin(user.ID)

	memberships, err := s.authRepo.GetUserMemberships(user.ID)
	if err != nil {
		return nil, err
	}

	availableTenants := make([]dto.TenantMembershipInfo, 0, len(memberships))
	var activeMembership *models.TenantMembership

	for _, m := range memberships {
		info := dto.TenantMembershipInfo{
			TenantID:     m.TenantID,
			RoleID:       m.RoleID,
			Status:       string(m.MembershipStatus),
			IsSuperAdmin: false,
		}
		if m.Tenant != nil {
			info.TenantName = m.Tenant.Name
			info.TenantSlug = m.Tenant.Slug
		}
		if m.Role != nil {
			info.RoleName = m.Role.Name
		}
		availableTenants = append(availableTenants, info)

		if requestedTenantID != uuid.Nil && m.TenantID == requestedTenantID {
			activeCopy := m
			activeMembership = &activeCopy
		}
	}

	// If no requested tenant was provided, pick first active tenant
	if activeMembership == nil && len(memberships) > 0 {
		activeMembership = &memberships[0]
	}

	// Determine active scope
	var activeTenantID uuid.UUID
	var activeRoleID uuid.UUID
	var activeRoleName string
	var permissions []string

	if activeMembership != nil {
		activeTenantID = activeMembership.TenantID
		activeRoleID = activeMembership.RoleID
		if activeMembership.Role != nil {
			activeRoleName = activeMembership.Role.Name
		}
		permissions, _ = s.authRepo.GetRolePermissions(activeRoleID)
	} else if isSuperAdmin {
		activeRoleID = superAdminRoleID
		activeRoleName = "Super Admin"
		permissions = []string{"*"}
	}

	// Generate Access Token (15m)
	secret := s.cfg.DatabaseURL // or JWT_SECRET
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
		activeTenantID,
		activeRoleID,
		activeRoleName,
		isSuperAdmin,
		permissions,
		secret,
		15*time.Minute,
	)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token (7 days)
	rawRefreshToken, err := utils.GenerateRandomToken(32)
	if err != nil {
		return nil, err
	}
	refreshRecord := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: utils.HashToken(rawRefreshToken),
		ExpiresAt: time.Now().UTC().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now().UTC(),
	}
	if err := s.authRepo.StoreRefreshToken(&refreshRecord); err != nil {
		return nil, err
	}

	var activeTenantInfo *dto.TenantMembershipInfo
	if activeMembership != nil {
		for _, info := range availableTenants {
			if info.TenantID == activeMembership.TenantID {
				activeTenantInfo = &info
				break
			}
		}
	} else if isSuperAdmin {
		activeTenantInfo = &dto.TenantMembershipInfo{
			TenantName:   "Platform Administration",
			RoleName:     "Super Admin",
			IsSuperAdmin: true,
		}
	}

	return &dto.AuthTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		ExpiresIn:    900,
		User: dto.UserProfileResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		ActiveTenant:     activeTenantInfo,
		AvailableTenants: availableTenants,
	}, nil
}
