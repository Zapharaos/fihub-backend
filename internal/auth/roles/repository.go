package roles

//go:generate mockgen -source=repository.go -destination=../../../test/mocks/roles_repository.go --package=mocks -mock_names=Repository=RolesRepository Repository

import (
	"github.com/google/uuid"
	"sync"
)

// Repository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Role
type Repository interface {
	Get(uuid uuid.UUID) (Role, bool, error)
	GetByName(name string) (Role, bool, error)
	GetWithPermissions(uuid uuid.UUID) (RoleWithPermissions, bool, error)
	GetAll() ([]Role, error)
	GetAllWithPermissions() (RolesWithPermissions, error)
	Create(role Role, permissionUUIDs []uuid.UUID) (uuid.UUID, error)
	Update(role Role, permissionUUIDs []uuid.UUID) error
	Delete(uuid uuid.UUID) error

	GetRolesByUserId(userUUID uuid.UUID) ([]Role, error)
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
