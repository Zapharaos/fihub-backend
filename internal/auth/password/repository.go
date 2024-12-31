package password

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

// Repository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Users
type Repository interface {
	Create(request Request) (Request, error)
	GetRequestID(userID uuid.UUID, token string) (uuid.UUID, error)
	GetExpiresAt(userID uuid.UUID) (time.Time, error)
	Delete(requestID uuid.UUID) error
	Valid(userID uuid.UUID, requestID uuid.UUID) (bool, error)
	ValidForUser(userID uuid.UUID) (bool, error)
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
