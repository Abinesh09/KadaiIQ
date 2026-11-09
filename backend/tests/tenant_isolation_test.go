package tests

import (
	"testing"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/Abinesh09/KadaiIQ/backend/internal/repositories"
	"github.com/Abinesh09/KadaiIQ/backend/migrations"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestStrictTenantIsolation(t *testing.T) {
	db := setupTestDB(t)
	if err := migrations.RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	now := time.Now().UTC()

	// 1. Create Tenant A (Coimbatore Footwear)
	tenantA := models.Tenant{
		Name:         "Coimbatore Footwear",
		Slug:         "cbe-footwear",
		BusinessName: "Coimbatore Footwear Pvt Ltd",
		Email:        "contact@cbefootwear.in",
		Status:       models.TenantStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := db.Create(&tenantA).Error; err != nil {
		t.Fatalf("Failed to create Tenant A: %v", err)
	}

	// 2. Create Tenant B (Madurai Retail)
	tenantB := models.Tenant{
		Name:         "Madurai Retail",
		Slug:         "madurai-retail",
		BusinessName: "Madurai Retail Organisers",
		Email:        "info@madurairetail.com",
		Status:       models.TenantStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := db.Create(&tenantB).Error; err != nil {
		t.Fatalf("Failed to create Tenant B: %v", err)
	}

	// 3. Create Tenant Features for Tenant A and Tenant B
	featureA := models.TenantFeature{
		TenantID:   tenantA.ID,
		FeatureKey: "pos_zenbiller_sync",
		Enabled:    true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	db.Create(&featureA)

	featureB := models.TenantFeature{
		TenantID:   tenantB.ID,
		FeatureKey: "gst_advanced_reports",
		Enabled:    true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	db.Create(&featureB)

	baseRepo := repositories.NewBaseTenantRepository(db)

	// 4. Test: Tenant A querying its features
	var featuresA []models.TenantFeature
	if err := baseRepo.FindAllByTenant(tenantA.ID, &featuresA); err != nil {
		t.Fatalf("Tenant A query failed: %v", err)
	}
	if len(featuresA) != 1 || featuresA[0].FeatureKey != "pos_zenbiller_sync" {
		t.Errorf("Expected Tenant A to only find pos_zenbiller_sync, got %v", featuresA)
	}

	// 5. Test: Tenant A cannot access Tenant B's feature record
	var crossFeature models.TenantFeature
	err := baseRepo.FindByIDAndTenant(tenantA.ID, featureB.ID, &crossFeature)
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected ErrRecordNotFound when Tenant A attempts to read Tenant B's record, got %v", err)
	}

	// 6. Test: Tenant A cannot delete Tenant B's feature record
	deleteErr := baseRepo.DeleteByTenant(tenantA.ID, &models.TenantFeature{}, featureB.ID)
	if deleteErr == nil {
		t.Errorf("Expected access denied error when Tenant A tries to delete Tenant B record, but succeeded")
	}

	// 7. Verify Tenant B's record is still intact
	var checkFeatureB models.TenantFeature
	if err := baseRepo.FindByIDAndTenant(tenantB.ID, featureB.ID, &checkFeatureB); err != nil {
		t.Errorf("Tenant B record was corrupted or missing: %v", err)
	}

	// 8. Test: Query with Nil/Empty Tenant ID must fail immediately
	_, nilScopeErr := baseRepo.ScopedDB(uuid.Nil)
	if nilScopeErr != repositories.ErrTenantIDRequired {
		t.Errorf("Expected ErrTenantIDRequired for Nil tenant ID, got %v", nilScopeErr)
	}
}
