package repositories

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
)

// PermissionRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on facts
type PermissionRepository interface {
	Create(permission models.Permission) (uuid.UUID, error)
	Get(uuid uuid.UUID) (models.Permission, bool, error)
	Update(permission models.Permission) error
	Delete(uuid uuid.UUID) error
	List() (models.Permissions, error)
}
