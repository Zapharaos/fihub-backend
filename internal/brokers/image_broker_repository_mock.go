package brokers

import "github.com/google/uuid"

// MockImageBrokerRepository represents a mock image broker repository
type MockImageBrokerRepository struct {
	image BrokerImage
	error error
	found bool
}

// NewMockImageBrokerRepository creates a new mock image broker repository
func NewMockImageBrokerRepository() ImageBrokerRepository {
	r := MockImageBrokerRepository{}
	var repo ImageBrokerRepository = &r
	return repo
}

func (m MockImageBrokerRepository) Create(_ BrokerImage) error {
	return m.error
}

func (m MockImageBrokerRepository) Get(_ uuid.UUID) (BrokerImage, bool, error) {
	return m.image, m.found, m.error
}

func (m MockImageBrokerRepository) Update(_ BrokerImage) error {
	return m.error
}

func (m MockImageBrokerRepository) Delete(_ uuid.UUID) error {
	return m.error
}

func (m MockImageBrokerRepository) Exists(_ uuid.UUID, _ uuid.UUID) (bool, error) {
	return m.found, m.error
}
