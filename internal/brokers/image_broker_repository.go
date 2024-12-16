package brokers

import "github.com/google/uuid"

// ImageBrokerRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Users
type ImageBrokerRepository interface {
	Create(image BrokerImage) error
	Get(brokerImageID uuid.UUID) (BrokerImage, bool, error)
	Update(image BrokerImage) error
	Delete(brokerImageID uuid.UUID) error
	Exists(brokerID uuid.UUID, brokerImageID uuid.UUID) (bool, error)
}
