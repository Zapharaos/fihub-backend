package permissions

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"sync"

	"github.com/google/uuid"
)

// Repository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on facts
type Repository interface {
	Get(uuid uuid.UUID) (models.Permission, bool, error)
	Create(permission models.Permission) (uuid.UUID, error)
	Update(permission models.Permission) error
	Delete(uuid uuid.UUID) error
	GetAll() ([]models.Permission, error)

	GetAllByRoleId(roleUUID uuid.UUID) ([]models.Permission, error)
	GetAllForUser(userUUID uuid.UUID) ([]models.Permission, error)
}

var (
	_globalRepositoryMu sync.RWMutex
	_globalRepository   Repository
)

// R is used to access the global repository singleton
func R() Repository {
	_globalRepositoryMu.RLock()
	defer _globalRepositoryMu.RUnlock()

	repository := _globalRepository
	return repository
}

// ReplaceGlobals affect a new repository to the global repository singleton
func ReplaceGlobals(repository Repository) func() {
	_globalRepositoryMu.Lock()
	defer _globalRepositoryMu.Unlock()

	prev := _globalRepository
	_globalRepository = repository
	return func() { ReplaceGlobals(prev) }
}
