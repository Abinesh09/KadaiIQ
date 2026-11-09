package tests

import (
	"testing"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/Abinesh09/KadaiIQ/backend/migrations"
	"github.com/google/uuid"
)

func TestModelCreationAndUniqueConstraints(t *testing.T) {
	db := setupTestDB(t)
	if err := migrations.RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// 1. Create a Tenant and verify UUID auto-generation
	tenant1 := models.Tenant{
		Name:         "Chennai Retail Org",
		Slug:         "chennai-retail",
		BusinessName: "Chennai Retail Pvt Ltd",
		Email:        "admin@chennairetail.com",
		Status:       models.TenantStatusActive,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := db.Create(&tenant1).Error; err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}
	if tenant1.ID == uuid.Nil {
		t.Errorf("Expected tenant ID to be generated, got Nil UUID")
	}

	// 2. Unique Tenant Slug constraint
	tenantDuplicateSlug := models.Tenant{
		Name:         "Another Retail",
		Slug:         "chennai-retail", // duplicate slug
		BusinessName: "Another Pvt Ltd",
		Email:        "other@retail.com",
		Status:       models.TenantStatusActive,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := db.Create(&tenantDuplicateSlug).Error; err == nil {
		t.Errorf("Expected error when inserting duplicate tenant slug, but got nil")
	}

	// 3. Create a User and verify Unique Email constraint
	user1 := models.User{
		Name:         "Rajesh Kumar",
		Email:        "rajesh@chennairetail.com",
		PasswordHash: "$2a$12$e8Yx9r8x...",
		Status:       models.UserStatusActive,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := db.Create(&user1).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userDuplicateEmail := models.User{
		Name:         "Rajesh Imposter",
		Email:        "rajesh@chennairetail.com", // duplicate email
		PasswordHash: "$2a$12$different...",
		Status:       models.UserStatusActive,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := db.Create(&userDuplicateEmail).Error; err == nil {
		t.Errorf("Expected error when inserting duplicate user email, but got nil")
	}

	// 4. Test Tenant Membership and unique (tenant_id, user_id) constraint
	var role models.Role
	if err := db.Where("name = ?", "Tenant Admin").First(&role).Error; err != nil {
		t.Fatalf("Failed to find Tenant Admin role: %v", err)
	}

	membership := models.TenantMembership{
		TenantID:         tenant1.ID,
		UserID:           user1.ID,
		RoleID:           role.ID,
		MembershipStatus: models.MembershipStatusActive,
		CreatedAt:        time.Now().UTC(),
	}
	if err := db.Create(&membership).Error; err != nil {
		t.Fatalf("Failed to create tenant membership: %v", err)
	}

	duplicateMembership := models.TenantMembership{
		TenantID:         tenant1.ID,
		UserID:           user1.ID, // same tenant and user
		RoleID:           role.ID,
		MembershipStatus: models.MembershipStatusActive,
		CreatedAt:        time.Now().UTC(),
	}
	if err := db.Create(&duplicateMembership).Error; err == nil {
		t.Errorf("Expected error when adding duplicate tenant membership for same tenant and user")
	}
}
