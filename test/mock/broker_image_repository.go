package mock

import (
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/google/uuid"
)

// BrokerImageRepository represents a mock brokers.ImageRepository
type BrokerImageRepository struct {
	image brokers.Image
	error error
	found bool
}

// NewBrokerImageRepository creates a new BrokerImageRepository of the brokers.ImageRepository interface
func NewBrokerImageRepository() brokers.ImageRepository {
	r := BrokerImageRepository{}
	var repo brokers.ImageRepository = &r
	return repo
}

func (m BrokerImageRepository) Create(_ brokers.Image) error {
	return m.error
}

func (m BrokerImageRepository) Get(_ uuid.UUID) (brokers.Image, bool, error) {
	return m.image, m.found, m.error
}

func (m BrokerImageRepository) Update(_ brokers.Image) error {
	return m.error
}

func (m BrokerImageRepository) Delete(_ uuid.UUID) error {
	return m.error
}

func (m BrokerImageRepository) Exists(_ uuid.UUID, _ uuid.UUID) (bool, error) {
	return m.found, m.error
}
