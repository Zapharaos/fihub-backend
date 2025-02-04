package brokers

import "github.com/google/uuid"

// MockUserBrokerRepository is a mock implementation of the UserBrokerRepository interface
type MockUserBrokerRepository struct {
	error       error
	found       bool
	userBrokers []UserBroker
}

// NewMockUserBrokerRepository creates a new instance of the MockUserBrokerRepository
func NewMockUserBrokerRepository() UserBrokerRepository {
	r := MockUserBrokerRepository{}
	var repo UserBrokerRepository = &r
	return repo
}

func (m MockUserBrokerRepository) Create(_ UserBroker) error {
	return m.error
}

func (m MockUserBrokerRepository) Delete(_ UserBroker) error {
	return m.error
}

func (m MockUserBrokerRepository) Exists(_ UserBroker) (bool, error) {
	return m.found, m.error
}

func (m MockUserBrokerRepository) GetAll(_ uuid.UUID) ([]UserBroker, error) {
	return m.userBrokers, m.error
}
