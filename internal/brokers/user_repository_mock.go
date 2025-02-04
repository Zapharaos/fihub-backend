package brokers

import "github.com/google/uuid"

// MockUserBrokerRepository is a mock implementation of the UserRepository interface
type MockUserBrokerRepository struct {
	error       error
	found       bool
	userBrokers []User
}

// NewMockUserBrokerRepository creates a new instance of the MockUserBrokerRepository
func NewMockUserBrokerRepository() UserRepository {
	r := MockUserBrokerRepository{}
	var repo UserRepository = &r
	return repo
}

func (m MockUserBrokerRepository) Create(_ User) error {
	return m.error
}

func (m MockUserBrokerRepository) Delete(_ User) error {
	return m.error
}

func (m MockUserBrokerRepository) Exists(_ User) (bool, error) {
	return m.found, m.error
}

func (m MockUserBrokerRepository) GetAll(_ uuid.UUID) ([]User, error) {
	return m.userBrokers, m.error
}
