package migrations

import (
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(Migration{
		Version:     "000001_initial_schema",
		Description: "Create tenants, users, roles, permissions, memberships, features, subscriptions, and audit logs tables",
		Up: func(db *gorm.DB) error {
			// Create Tables
			tables := []interface{}{
				&models.Tenant{},
				&models.User{},
				&models.Role{},
				&models.Permission{},
				&models.TenantMembership{},
				&models.RolePermission{},
				&models.TenantFeature{},
				&models.SubscriptionPlan{},
				&models.AuditLog{},
			}

			for _, t := range tables {
				if err := db.AutoMigrate(t); err != nil {
					return err
				}
			}

			// Seed Initial System Roles
			now := time.Now().UTC()
			initialRoles := []models.Role{
				{
					ID:          uuid.New(),
					Name:        "Super Admin",
					Description: "Platform-level administrator with cross-tenant management privileges",
					Scope:       models.RoleScopePlatform,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          uuid.New(),
					Name:        "Tenant Owner",
					Description: "Full administrative and financial authority over organisation workspace",
					Scope:       models.RoleScopeTenant,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          uuid.New(),
					Name:        "Tenant Admin",
					Description: "Organisation administrator for inventory, products, sales, and staff",
					Scope:       models.RoleScopeTenant,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          uuid.New(),
					Name:        "Manager",
					Description: "Store manager with access to sales, inventory operations, and reports",
					Scope:       models.RoleScopeTenant,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          uuid.New(),
					Name:        "Staff",
					Description: "Cashier and floor staff with basic POS entry and product lookup",
					Scope:       models.RoleScopeTenant,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			}

			for _, r := range initialRoles {
				var count int64
				db.Model(&models.Role{}).Where("name = ?", r.Name).Count(&count)
				if count == 0 {
					if err := db.Create(&r).Error; err != nil {
						return err
					}
				}
			}

			// Seed Initial Core Permissions
			initialPerms := []models.Permission{
				{ID: uuid.New(), Name: "products:view", Module: "products", Description: "View product catalog and variants", CreatedAt: now},
				{ID: uuid.New(), Name: "products:manage", Module: "products", Description: "Create, edit, and deactivate products", CreatedAt: now},
				{ID: uuid.New(), Name: "inventory:view", Module: "inventory", Description: "View stock levels and movement history", CreatedAt: now},
				{ID: uuid.New(), Name: "inventory:adjust", Module: "inventory", Description: "Perform manual stock adjustments and restocks", CreatedAt: now},
				{ID: uuid.New(), Name: "sales:view", Module: "sales", Description: "View sales bills and transactions", CreatedAt: now},
				{ID: uuid.New(), Name: "sales:import", Module: "sales", Description: "Import POS records and bills", CreatedAt: now},
				{ID: uuid.New(), Name: "analytics:view", Module: "analytics", Description: "Access analytics KPIs, sales trends and profit reports", CreatedAt: now},
				{ID: uuid.New(), Name: "reports:download", Module: "reports", Description: "Generate and export Excel, CSV and PDF reports", CreatedAt: now},
				{ID: uuid.New(), Name: "gst:manage", Module: "gst", Description: "Configure tax rates, slabs, and GST reports", CreatedAt: now},
				{ID: uuid.New(), Name: "users:manage", Module: "users", Description: "Manage store staff memberships and roles", CreatedAt: now},
			}

			for _, p := range initialPerms {
				var count int64
				db.Model(&models.Permission{}).Where("name = ?", p.Name).Count(&count)
				if count == 0 {
					if err := db.Create(&p).Error; err != nil {
						return err
					}
				}
			}

			// Seed Initial Subscription Plans
			initialPlans := []models.SubscriptionPlan{
				{
					ID:          uuid.New(),
					Name:        "Starter",
					Description: "Ideal for single retail stores managing inventory and basic analytics",
					Status:      models.SubscriptionPlanStatusActive,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          uuid.New(),
					Name:        "Professional",
					Description: "Advanced analytics, multiple POS integrations, GST summaries, and export capabilities",
					Status:      models.SubscriptionPlanStatusActive,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          uuid.New(),
					Name:        "Enterprise",
					Description: "High-volume retail chains with multi-outlet analytics and priority support",
					Status:      models.SubscriptionPlanStatusActive,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			}

			for _, sp := range initialPlans {
				var count int64
				db.Model(&models.SubscriptionPlan{}).Where("name = ?", sp.Name).Count(&count)
				if count == 0 {
					if err := db.Create(&sp).Error; err != nil {
						return err
					}
				}
			}

			return nil
		},
		Down: func(db *gorm.DB) error {
			// Drop tables in reverse order
			tables := []string{
				"audit_logs",
				"tenant_features",
				"role_permissions",
				"tenant_memberships",
				"permissions",
				"roles",
				"subscription_plans",
				"users",
				"tenants",
			}
			for _, tbl := range tables {
				if err := db.Migrator().DropTable(tbl); err != nil {
					return err
				}
			}
			return nil
		},
	})
}
