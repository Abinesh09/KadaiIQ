package repositories

import (
	"errors"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) FindUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) UpdateUserPassword(userID uuid.UUID, passwordHash string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("password_hash", passwordHash).Error
}

func (r *AuthRepository) GetUserMemberships(userID uuid.UUID) ([]models.TenantMembership, error) {
	var memberships []models.TenantMembership
	err := r.db.Preload("Tenant").Preload("Role").Where("user_id = ?", userID).Find(&memberships).Error
	return memberships, err
}

func (r *AuthRepository) CheckUserIsSuperAdmin(userID uuid.UUID) (bool, uuid.UUID, error) {
	var role models.Role
	err := r.db.Where("name = ? AND scope = ?", "Super Admin", models.RoleScopePlatform).First(&role).Error
	if err != nil {
		return false, uuid.Nil, err
	}

	var count int64
	r.db.Model(&models.TenantMembership{}).Where("user_id = ? AND role_id = ?", userID, role.ID).Count(&count)
	if count > 0 {
		return true, role.ID, nil
	}

	// Also check if user has email matching seeded superadmin
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err == nil {
		if user.Email == "superadmin@kadaiiq.com" {
			return true, role.ID, nil
		}
	}

	return false, uuid.Nil, nil
}

func (r *AuthRepository) GetRolePermissions(roleID uuid.UUID) ([]string, error) {
	var rolePerms []models.RolePermission
	err := r.db.Preload("Permission").Where("role_id = ?", roleID).Find(&rolePerms).Error
	if err != nil {
		return nil, err
	}

	perms := make([]string, 0, len(rolePerms))
	for _, rp := range rolePerms {
		if rp.Permission != nil {
			perms = append(perms, rp.Permission.Name)
		}
	}
	return perms, nil
}

// Token operations
func (r *AuthRepository) StoreRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) FindRefreshToken(tokenHash string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash, time.Now().UTC()).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *AuthRepository) RevokeRefreshToken(tokenHash string) error {
	now := time.Now().UTC()
	return r.db.Model(&models.RefreshToken{}).Where("token_hash = ?", tokenHash).Update("revoked_at", &now).Error
}

func (r *AuthRepository) RevokeAllUserTokens(userID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.Model(&models.RefreshToken{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", &now).Error
}

// Password reset operations
func (r *AuthRepository) StorePasswordResetToken(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) FindValidPasswordResetToken(tokenHash string) (*models.PasswordResetToken, error) {
	var prt models.PasswordResetToken
	err := r.db.Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", tokenHash, time.Now().UTC()).First(&prt).Error
	if err != nil {
		return nil, err
	}
	return &prt, nil
}

func (r *AuthRepository) MarkPasswordResetTokenUsed(id uuid.UUID) error {
	now := time.Now().UTC()
	res := r.db.Model(&models.PasswordResetToken{}).Where("id = ? AND used_at IS NULL", id).Update("used_at", &now)
	if res.RowsAffected == 0 {
		return errors.New("token was already used")
	}
	return res.Error
}
