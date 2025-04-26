package repositories

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
)

// RoleRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Role
type RoleRepository interface {
	Get(uuid uuid.UUID) (models.Role, bool, error)
	GetByName(name string) (models.Role, bool, error)
	GetWithPermissions(uuid uuid.UUID) (models.RoleWithPermissions, bool, error)
	GetAll() ([]models.Role, error)
	GetAllWithPermissions() (models.RolesWithPermissions, error)
	Create(role models.Role, permissionUUIDs []uuid.UUID) (uuid.UUID, error)
	Update(role models.Role, permissionUUIDs []uuid.UUID) error
	Delete(uuid uuid.UUID) error

	GetRolesByUserId(userUUID uuid.UUID) ([]models.Role, error)
	SetRolePermissions(roleUUID uuid.UUID, permissionUUIDs []uuid.UUID) error
}
