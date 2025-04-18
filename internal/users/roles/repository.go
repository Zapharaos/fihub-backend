package roles

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"sync"
)

// Repository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Role
type Repository interface {
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

var (
	_globalRepositoryMu sync.RWMutex
	_globalRepository   Repository
)

// R is used to access the global repository singleton
func R() Repository {
	_globalRepositoryMu.RLock()
	repository := _globalRepository
	_globalRepositoryMu.RUnlock()
	return repository
}

// ReplaceGlobals affect a new repository to the global repository singleton
func ReplaceGlobals(repository Repository) func() {
	_globalRepositoryMu.Lock()
	prev := _globalRepository
	_globalRepository = repository
	_globalRepositoryMu.Unlock()
	return func() { ReplaceGlobals(prev) }
}
