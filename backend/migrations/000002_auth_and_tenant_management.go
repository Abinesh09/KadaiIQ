package migrations

import (
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(Migration{
		Version:     "000002_auth_and_tenant_management",
		Description: "Create refresh_tokens, password_reset_tokens and seed platform Super Admin user",
		Up: func(db *gorm.DB) error {
			// Migrate tokens and updated tenants
			tables := []interface{}{
				&models.Tenant{},
				&models.RefreshToken{},
				&models.PasswordResetToken{},
			}
			for _, t := range tables {
				if err := db.AutoMigrate(t); err != nil {
					return err
				}
			}

			// Ensure Super Admin Role exists
			var superAdminRole models.Role
			if err := db.Where("name = ?", "Super Admin").First(&superAdminRole).Error; err != nil {
				return err
			}

			// Seed initial Super Admin User if none exists
			var adminCount int64
			db.Model(&models.User{}).Where("email = ?", "superadmin@kadaiiq.com").Count(&adminCount)
			if adminCount == 0 {
				hash, err := bcrypt.GenerateFromPassword([]byte("Admin@KadaiIQ2026"), bcrypt.DefaultCost)
				if err != nil {
					return err
				}

				now := time.Now().UTC()
				adminUser := models.User{
					ID:           uuid.New(),
					Name:         "Super Administrator",
					Email:        "superadmin@kadaiiq.com",
					PasswordHash: string(hash),
					Status:       models.UserStatusActive,
					CreatedAt:    now,
					UpdatedAt:    now,
				}
				if err := db.Create(&adminUser).Error; err != nil {
					return err
				}

				// Assign platform membership (empty/Nil tenant or dedicated membership)
				// For Super Admin, platform-level access is indicated by having the Super Admin role.
			}

			return nil
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable("password_reset_tokens", "refresh_tokens")
		},
	})
}
