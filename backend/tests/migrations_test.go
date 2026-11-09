package tests

import (
	"testing"

	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/Abinesh09/KadaiIQ/backend/migrations"
)

func TestMigrationsExecutionAndIdempotency(t *testing.T) {
	db := setupTestDB(t)

	// Step 1: Run migrations
	if err := migrations.RunMigrations(db); err != nil {
		t.Fatalf("RunMigrations failed: %v", err)
	}

	// Step 2: Verify schema_migrations record
	var sm migrations.SchemaMigration
	if err := db.Where("version = ?", "000001_initial_schema").First(&sm).Error; err != nil {
		t.Fatalf("Expected version 000001_initial_schema to be recorded in schema_migrations, got: %v", err)
	}

	// Step 3: Verify seeded roles
	var roleCount int64
	if err := db.Model(&models.Role{}).Count(&roleCount).Error; err != nil {
		t.Fatalf("Failed to query roles count: %v", err)
	}
	if roleCount < 5 {
		t.Errorf("Expected at least 5 initial roles, got %d", roleCount)
	}

	// Step 4: Verify seeded subscription plans
	var planCount int64
	if err := db.Model(&models.SubscriptionPlan{}).Count(&planCount).Error; err != nil {
		t.Fatalf("Failed to query subscription plans: %v", err)
	}
	if planCount < 3 {
		t.Errorf("Expected at least 3 initial subscription plans, got %d", planCount)
	}

	// Step 5: Test idempotency (running again must not fail or duplicate data)
	if err := migrations.RunMigrations(db); err != nil {
		t.Fatalf("Second RunMigrations execution failed: %v", err)
	}

	var secondRoleCount int64
	db.Model(&models.Role{}).Count(&secondRoleCount)
	if secondRoleCount != roleCount {
		t.Errorf("Roles were duplicated after second migration run: %d vs %d", secondRoleCount, roleCount)
	}
}
