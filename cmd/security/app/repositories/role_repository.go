package repositories

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
)

// RoleRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Role
type RoleRepository interface {
	Create(role models.Role, permissionUUIDs []uuid.UUID) (uuid.UUID, error)
	Get(uuid uuid.UUID) (models.Role, bool, error)
	GetByName(name string) (models.Role, bool, error)
	GetWithPermissions(uuid uuid.UUID) (models.RoleWithPermissions, bool, error)
	Update(role models.Role, permissionUUIDs []uuid.UUID) error
	Delete(uuid uuid.UUID) error
	List() (models.Roles, error)
	ListByUserId(userUUID uuid.UUID) (models.Roles, error)
	ListWithPermissions() (models.RolesWithPermissions, error)
	ListWithPermissionsByUserId(userUUID uuid.UUID) (models.RolesWithPermissions, error)

	SetForUser(userUUID uuid.UUID, roleUUIDs []uuid.UUID) error
	AddToUsers(userUUIDs []uuid.UUID, id uuid.UUID) error
	RemoveFromUsers(userUUIDs []uuid.UUID, roleUUID uuid.UUID) error

	ListUsersByRoleId(roleUUID uuid.UUID) ([]string, error)
	ListUsers() ([]string, error)

	SetPermissionsByRoleId(roleUUID uuid.UUID, permissionUUIDs []uuid.UUID) error
	ListPermissionsByRoleId(roleUUID uuid.UUID) (models.Permissions, error)
	ListPermissionsByUserId(userUUID uuid.UUID) (models.Permissions, error)
}
