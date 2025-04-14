package transaction

import (
	"github.com/google/uuid"
	"sync"
)

// Repository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Transaction
type Repository interface {
	Create(transactionInput TransactionInput) (uuid.UUID, error)
	Get(transactionID uuid.UUID) (Transaction, bool, error)
	Update(transactionInput TransactionInput) error
	Delete(transaction Transaction) error
	Exists(transactionID uuid.UUID, userID uuid.UUID) (bool, error)
	GetAll(userID uuid.UUID) ([]Transaction, error)
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
