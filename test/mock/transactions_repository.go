package transactions

import "github.com/google/uuid"

// MockRepository represents a mock Repository
type MockRepository struct {
	ID           uuid.UUID
	error        error
	found        bool
	transaction  Transaction
	transactions []Transaction
}

// NewMockRepository creates a new MockRepository of the Repository interface
func NewMockRepository() Repository {
	r := MockRepository{}
	var repo Repository = &r
	return repo
}

func (r MockRepository) Create(_ TransactionInput) (uuid.UUID, error) {
	return r.ID, r.error
}

func (r MockRepository) Get(_ uuid.UUID) (Transaction, bool, error) {
	return r.transaction, r.found, r.error
}

func (r MockRepository) Update(_ TransactionInput) error {
	return r.error
}

func (r MockRepository) Delete(_ Transaction) error {
	return r.error
}

func (r MockRepository) Exists(_ uuid.UUID, _ uuid.UUID) (bool, error) {
	return r.found, r.error
}

func (r MockRepository) GetAll(_ uuid.UUID) ([]Transaction, error) {
	return r.transactions, r.error
}
