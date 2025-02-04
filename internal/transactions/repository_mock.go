package transactions

import "github.com/google/uuid"

type MockRepository struct {
	ID           uuid.UUID
	error        error
	found        bool
	transaction  Transaction
	transactions []Transaction
}

func NewMockRepository() Repository {
	r := MockRepository{}
	var repo Repository = &r
	return repo
}

func (r MockRepository) Create(transactionInput TransactionInput) (uuid.UUID, error) {
	return r.ID, r.error
}

func (r MockRepository) Get(transactionID uuid.UUID) (Transaction, bool, error) {
	return r.transaction, r.found, r.error
}

func (r MockRepository) Update(transactionInput TransactionInput) error {
	return r.error
}

func (r MockRepository) Delete(transaction Transaction) error {
	return r.error
}

func (r MockRepository) Exists(transactionID uuid.UUID, userID uuid.UUID) (bool, error) {
	return r.found, r.error
}

func (r MockRepository) GetAll(userID uuid.UUID) ([]Transaction, error) {
	return r.transactions, r.error
}
