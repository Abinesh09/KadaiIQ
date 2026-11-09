package migrations

import (
	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(Migration{
		Version:     "000004_sales_and_import_jobs",
		Description: "Create sales, sale_items, and import_jobs tables with deduplication indexes",
		Up: func(db *gorm.DB) error {
			tables := []interface{}{
				&models.Sale{},
				&models.SaleItem{},
				&models.ImportJob{},
			}
			for _, t := range tables {
				if err := db.AutoMigrate(t); err != nil {
					return err
				}
			}
			return nil
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(
				"sale_items",
				"sales",
				"import_jobs",
			)
		},
	})
}
