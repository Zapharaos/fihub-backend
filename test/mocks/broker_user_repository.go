package mocks

import (
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/google/uuid"
)

// BrokerUserRepository is a mocks brokers.UserRepository
type BrokerUserRepository struct {
	error       error
	found       bool
	userBrokers []brokers.User
}

// NewBrokerUserRepository creates a new BrokerUserRepository of the brokers.UserRepository
func NewBrokerUserRepository() brokers.UserRepository {
	r := BrokerUserRepository{}
	var repo brokers.UserRepository = &r
	return repo
}

func (m BrokerUserRepository) Create(_ brokers.User) error {
	return m.error
}

func (m BrokerUserRepository) Delete(_ brokers.User) error {
	return m.error
}

func (m BrokerUserRepository) Exists(_ brokers.User) (bool, error) {
	return m.found, m.error
}

func (m BrokerUserRepository) GetAll(_ uuid.UUID) ([]brokers.User, error) {
	return m.userBrokers, m.error
}
