package mock

import (
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/google/uuid"
)

// BrokerRepository represents a mock brokers.Repository
type BrokerRepository struct {
	ID      uuid.UUID
	error   error
	found   bool
	Broker  brokers.Broker
	Brokers []brokers.Broker
}

// NewBrokerRepository creates a new BrokerRepository of the brokers.BrokerRepository interface
func NewBrokerRepository() brokers.BrokerRepository {
	r := BrokerRepository{}
	var repo brokers.BrokerRepository = &r
	return repo
}

func (m BrokerRepository) Get(_ uuid.UUID) (brokers.Broker, bool, error) {
	return m.Broker, m.found, m.error
}

func (m BrokerRepository) Create(_ brokers.Broker) (uuid.UUID, error) {
	return m.ID, m.error
}

func (m BrokerRepository) Update(_ brokers.Broker) error {
	return m.error
}

func (m BrokerRepository) Delete(_ uuid.UUID) error {
	return m.error
}

func (m BrokerRepository) Exists(_ uuid.UUID) (bool, error) {
	return m.found, m.error
}

func (m BrokerRepository) ExistsByName(_ string) (bool, error) {
	return m.found, m.error
}

func (m BrokerRepository) GetAll() ([]brokers.Broker, error) {
	return m.Brokers, m.error
}

func (m BrokerRepository) GetAllEnabled() ([]brokers.Broker, error) {
	return m.Brokers, m.error
}

func (m BrokerRepository) SetImage(_ uuid.UUID, _ uuid.UUID) error {
	return m.error
}

func (m BrokerRepository) HasImage(_ uuid.UUID) (bool, error) {
	return m.found, m.error
}

func (m BrokerRepository) DeleteImage(_ uuid.UUID) error {
	return m.error
}
