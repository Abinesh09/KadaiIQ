package migrations

import (
	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(Migration{
		Version:     "000003_products_and_inventory",
		Description: "Create brands, categories, products, product_variants, and stock_movements tables",
		Up: func(db *gorm.DB) error {
			tables := []interface{}{
				&models.Brand{},
				&models.Category{},
				&models.Product{},
				&models.ProductVariant{},
				&models.StockMovement{},
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
				"stock_movements",
				"product_variants",
				"products",
				"categories",
				"brands",
			)
		},
	})
}
