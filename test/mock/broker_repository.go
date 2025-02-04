package brokers

import "github.com/google/uuid"

// MockBrokerRepository represents a mock BrokerRepository
type MockBrokerRepository struct {
	ID      uuid.UUID
	error   error
	found   bool
	Broker  Broker
	Brokers []Broker
}

// NewMockBrokerRepository creates a new MockBrokerRepository of the BrokerRepository interface
func NewMockBrokerRepository() BrokerRepository {
	r := MockBrokerRepository{}
	var repo BrokerRepository = &r
	return repo
}

func (m MockBrokerRepository) Get(_ uuid.UUID) (Broker, bool, error) {
	return m.Broker, m.found, m.error
}

func (m MockBrokerRepository) Create(_ Broker) (uuid.UUID, error) {
	return m.ID, m.error
}

func (m MockBrokerRepository) Update(_ Broker) error {
	return m.error
}

func (m MockBrokerRepository) Delete(_ uuid.UUID) error {
	return m.error
}

func (m MockBrokerRepository) Exists(_ uuid.UUID) (bool, error) {
	return m.found, m.error
}

func (m MockBrokerRepository) ExistsByName(_ string) (bool, error) {
	return m.found, m.error
}

func (m MockBrokerRepository) GetAll() ([]Broker, error) {
	return m.Brokers, m.error
}

func (m MockBrokerRepository) GetAllEnabled() ([]Broker, error) {
	return m.Brokers, m.error
}

func (m MockBrokerRepository) SetImage(_ uuid.UUID, _ uuid.UUID) error {
	return m.error
}

func (m MockBrokerRepository) HasImage(_ uuid.UUID) (bool, error) {
	return m.found, m.error
}

func (m MockBrokerRepository) DeleteImage(_ uuid.UUID) error {
	return m.error
}
