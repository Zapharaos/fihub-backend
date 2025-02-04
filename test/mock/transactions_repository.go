package mock

import (
	transactions2 "github.com/Zapharaos/fihub-backend/internal/transactions"
	"github.com/google/uuid"
)

// TransactionsRepository represents a mock transactions.Repository
type TransactionsRepository struct {
	ID           uuid.UUID
	error        error
	found        bool
	transaction  transactions2.Transaction
	transactions []transactions2.Transaction
}

// NewTransactionsRepository creates a new TransactionsRepository of the transactions.Repository interface
func NewTransactionsRepository() transactions2.Repository {
	r := TransactionsRepository{}
	var repo transactions2.Repository = &r
	return repo
}

func (r TransactionsRepository) Create(_ transactions2.TransactionInput) (uuid.UUID, error) {
	return r.ID, r.error
}

func (r TransactionsRepository) Get(_ uuid.UUID) (transactions2.Transaction, bool, error) {
	return r.transaction, r.found, r.error
}

func (r TransactionsRepository) Update(_ transactions2.TransactionInput) error {
	return r.error
}

func (r TransactionsRepository) Delete(_ transactions2.Transaction) error {
	return r.error
}

func (r TransactionsRepository) Exists(_ uuid.UUID, _ uuid.UUID) (bool, error) {
	return r.found, r.error
}

func (r TransactionsRepository) GetAll(_ uuid.UUID) ([]transactions2.Transaction, error) {
	return r.transactions, r.error
}
