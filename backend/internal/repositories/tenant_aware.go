package repositories

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrTenantIDRequired = errors.New("operation requires a valid tenant ID to prevent cross-tenant data leakage")
)

// TenantScope returns a GORM scope that strictly filters queries by tenant_id
func TenantScope(tenantID uuid.UUID) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if tenantID == uuid.Nil {
			_ = db.AddError(ErrTenantIDRequired)
			return db
		}
		return db.Where("tenant_id = ?", tenantID)
	}
}

// BaseTenantRepository provides reusable tenant-scoped database methods
type BaseTenantRepository struct {
	db *gorm.DB
}

func NewBaseTenantRepository(db *gorm.DB) *BaseTenantRepository {
	return &BaseTenantRepository{db: db}
}

// ScopedDB returns a GORM database session strictly restricted to the specified tenant
func (r *BaseTenantRepository) ScopedDB(tenantID uuid.UUID) (*gorm.DB, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	return r.db.Scopes(TenantScope(tenantID)), nil
}

// FindAllByTenant queries records strictly belonging to a tenant
func (r *BaseTenantRepository) FindAllByTenant(tenantID uuid.UUID, out interface{}, conditions ...interface{}) error {
	scoped, err := r.ScopedDB(tenantID)
	if err != nil {
		return err
	}
	if len(conditions) > 0 {
		return scoped.Find(out, conditions...).Error
	}
	return scoped.Find(out).Error
}

// FindByIDAndTenant queries a single record by ID while strictly verifying tenant ownership
func (r *BaseTenantRepository) FindByIDAndTenant(tenantID uuid.UUID, id interface{}, out interface{}) error {
	scoped, err := r.ScopedDB(tenantID)
	if err != nil {
		return err
	}
	return scoped.First(out, "id = ?", id).Error
}

// DeleteByTenant deletes a record while ensuring it strictly belongs to the specified tenant
func (r *BaseTenantRepository) DeleteByTenant(tenantID uuid.UUID, model interface{}, id interface{}) error {
	scoped, err := r.ScopedDB(tenantID)
	if err != nil {
		return err
	}
	res := scoped.Where("id = ?", id).Delete(model)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("record not found or access denied for tenant %s", tenantID)
	}
	return nil
}
