package migrations

import (
	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(Migration{
		Version:     "000005_gst_and_reports",
		Description: "Create tenant_gst_settings table for Indian GST tax compliance",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&models.TenantGSTSetting{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable("tenant_gst_settings")
		},
	})
}
