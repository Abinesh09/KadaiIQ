package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) DB() *gorm.DB {
	return r.db
}

func (r *TenantRepository) GetAllTenants() ([]models.Tenant, error) {
	var tenants []models.Tenant
	err := r.db.Preload("Plan").Preload("Memberships").Order("created_at desc").Find(&tenants).Error
	return tenants, err
}

func (r *TenantRepository) GetTenantByID(id uuid.UUID) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.Preload("Plan").Preload("Features").First(&tenant, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepository) GetTenantBySlug(slug string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.Preload("Plan").First(&tenant, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepository) UpdateTenant(tenant *models.Tenant) error {
	return r.db.Save(tenant).Error
}

func (r *TenantRepository) SetTenantStatus(id uuid.UUID, status models.TenantStatus) error {
	return r.db.Model(&models.Tenant{}).Where("id = ?", id).Update("status", status).Error
}

// Transactional onboarding: Tenant + Owner User + Membership + Initial Features + Audit Log
func (r *TenantRepository) OnboardTenantWithTx(
	tenant *models.Tenant,
	owner *models.User,
	auditActorID uuid.UUID,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Check if slug exists
		var existingSlugCount int64
		if err := tx.Model(&models.Tenant{}).Where("slug = ?", tenant.Slug).Count(&existingSlugCount).Error; err != nil {
			return err
		}
		if existingSlugCount > 0 {
			return fmt.Errorf("tenant slug '%s' is already in use", tenant.Slug)
		}

		// 2. Create Tenant
		if err := tx.Create(tenant).Error; err != nil {
			return fmt.Errorf("failed to create tenant: %w", err)
		}

		// 3. Find or Create Owner User
		var user models.User
		err := tx.Where("LOWER(email) = LOWER(?)", owner.Email).First(&user).Error
		if err == gorm.ErrRecordNotFound {
			user = *owner
			if err := tx.Create(&user).Error; err != nil {
				return fmt.Errorf("failed to create owner user: %w", err)
			}
		} else if err != nil {
			return err
		}

		// 4. Find Tenant Owner Role
		var ownerRole models.Role
		if err := tx.Where("name = ?", "Tenant Owner").First(&ownerRole).Error; err != nil {
			return fmt.Errorf("failed to locate Tenant Owner role: %w", err)
		}

		// 5. Create Tenant Membership
		membership := models.TenantMembership{
			ID:               uuid.New(),
			TenantID:         tenant.ID,
			UserID:           user.ID,
			RoleID:           ownerRole.ID,
			MembershipStatus: models.MembershipStatusActive,
			CreatedAt:        time.Now().UTC(),
		}
		if err := tx.Create(&membership).Error; err != nil {
			return fmt.Errorf("failed to assign owner membership: %w", err)
		}

		// 6. Seed Default Features
		defaultFeatures := []string{
			"dashboard_access",
			"products_management",
			"inventory_tracking",
			"sales_ledger",
			"pos_import",
			"analytics_basic",
			"reports_export",
			"gst_summary",
		}
		for _, feat := range defaultFeatures {
			tf := models.TenantFeature{
				ID:         uuid.New(),
				TenantID:   tenant.ID,
				FeatureKey: feat,
				Enabled:    true,
				CreatedAt:  time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			}
			if err := tx.Create(&tf).Error; err != nil {
				return fmt.Errorf("failed to initialize feature %s: %w", feat, err)
			}
		}

		// 7. Write Audit Log
		actorID := auditActorID
		if actorID == uuid.Nil {
			actorID = user.ID
		}
		auditLog := models.AuditLog{
			ID:         uuid.New(),
			TenantID:   &tenant.ID,
			UserID:     actorID,
			Action:     "TENANT_ONBOARDED",
			EntityType: "tenant",
			EntityID:   tenant.ID.String(),
			Details:    fmt.Sprintf("Created organisation %s (slug: %s) with owner %s", tenant.Name, tenant.Slug, user.Email),
			CreatedAt:  time.Now().UTC(),
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to record audit log: %w", err)
		}

		return nil
	})
}

func (r *TenantRepository) ToggleFeature(tenantID uuid.UUID, featureKey string, enabled bool) error {
	var tf models.TenantFeature
	err := r.db.Where("tenant_id = ? AND feature_key = ?", tenantID, featureKey).First(&tf).Error
	if err == gorm.ErrRecordNotFound {
		tf = models.TenantFeature{
			ID:         uuid.New(),
			TenantID:   tenantID,
			FeatureKey: featureKey,
			Enabled:    enabled,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}
		return r.db.Create(&tf).Error
	}
	if err != nil {
		return err
	}
	return r.db.Model(&tf).Update("enabled", enabled).Error
}

func (r *TenantRepository) GetTenantMembers(tenantID uuid.UUID) ([]models.TenantMembership, error) {
	var members []models.TenantMembership
	err := r.db.Preload("User").Preload("Role").Where("tenant_id = ?", tenantID).Find(&members).Error
	return members, err
}

func (r *TenantRepository) AddTenantMember(tenantID, userID, roleID uuid.UUID) (*models.TenantMembership, error) {
	var count int64
	r.db.Model(&models.TenantMembership{}).Where("tenant_id = ? AND user_id = ?", tenantID, userID).Count(&count)
	if count > 0 {
		return nil, errors.New("user is already a member of this tenant")
	}

	membership := models.TenantMembership{
		ID:               uuid.New(),
		TenantID:         tenantID,
		UserID:           userID,
		RoleID:           roleID,
		MembershipStatus: models.MembershipStatusActive,
		CreatedAt:        time.Now().UTC(),
	}
	if err := r.db.Create(&membership).Error; err != nil {
		return nil, err
	}
	return &membership, nil
}

func (r *TenantRepository) UpdateMemberRole(tenantID, memberID, roleID uuid.UUID) error {
	res := r.db.Model(&models.TenantMembership{}).Where("id = ? AND tenant_id = ?", memberID, tenantID).Update("role_id", roleID)
	if res.RowsAffected == 0 {
		return errors.New("member not found in this tenant")
	}
	return res.Error
}

func (r *TenantRepository) RecordAuditLog(log *models.AuditLog) error {
	return r.db.Create(log).Error
}
