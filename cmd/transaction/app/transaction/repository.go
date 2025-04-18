package transaction

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"sync"
)

// Repository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Transaction
type Repository interface {
	Create(transactionInput models.TransactionInput) (uuid.UUID, error)
	Get(transactionID uuid.UUID) (models.Transaction, bool, error)
	Update(transactionInput models.TransactionInput) error
	Delete(transaction models.Transaction) error
	Exists(transactionID uuid.UUID, userID uuid.UUID) (bool, error)
	GetAll(userID uuid.UUID) ([]models.Transaction, error)
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
