package migrations

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type SchemaMigration struct {
	Version   string    `gorm:"primaryKey;type:varchar(100)" json:"version"`
	AppliedAt time.Time `gorm:"not null" json:"applied_at"`
}

type Migration struct {
	Version     string
	Description string
	Up          func(db *gorm.DB) error
	Down        func(db *gorm.DB) error
}

var registry = []Migration{}

func RegisterMigration(m Migration) {
	registry = append(registry, m)
}

func EnsureMigrationTable(db *gorm.DB) error {
	return db.AutoMigrate(&SchemaMigration{})
}

// RunMigrations executes all unapplied registered migrations in order
func RunMigrations(db *gorm.DB) error {
	if err := EnsureMigrationTable(db); err != nil {
		return fmt.Errorf("failed to initialize schema_migrations table: %w", err)
	}

	for _, m := range registry {
		var existing SchemaMigration
		err := db.Where("version = ?", m.Version).First(&existing).Error
		if err == nil {
			// Already applied
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to check migration status for %s: %w", m.Version, err)
		}

		log.Printf("[Migration] Applying %s: %s", m.Version, m.Description)
		// Run inside a database transaction
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if err := m.Up(tx); err != nil {
				return err
			}
			return tx.Create(&SchemaMigration{
				Version:   m.Version,
				AppliedAt: time.Now().UTC(),
			}).Error
		})

		if txErr != nil {
			return fmt.Errorf("failed migration %s: %w", m.Version, txErr)
		}
		log.Printf("[Migration] Successfully applied %s", m.Version)
	}

	return nil
}

// RollbackLatest rolls back the last applied migration
func RollbackLatest(db *gorm.DB) error {
	if err := EnsureMigrationTable(db); err != nil {
		return fmt.Errorf("failed to verify schema_migrations table: %w", err)
	}

	var latest SchemaMigration
	if err := db.Order("applied_at desc").First(&latest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("[Migration] No migrations to rollback")
			return nil
		}
		return err
	}

	for i := len(registry) - 1; i >= 0; i-- {
		m := registry[i]
		if m.Version == latest.Version {
			log.Printf("[Migration] Rolling back %s: %s", m.Version, m.Description)
			txErr := db.Transaction(func(tx *gorm.DB) error {
				if m.Down != nil {
					if err := m.Down(tx); err != nil {
						return err
					}
				}
				return tx.Where("version = ?", m.Version).Delete(&SchemaMigration{}).Error
			})

			if txErr != nil {
				return fmt.Errorf("failed rollback for %s: %w", m.Version, txErr)
			}
			log.Printf("[Migration] Successfully rolled back %s", m.Version)
			return nil
		}
	}

	return fmt.Errorf("migration definition for %s not found in registry", latest.Version)
}
