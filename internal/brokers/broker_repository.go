package brokers

import "github.com/google/uuid"

// BrokerRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Broker
type BrokerRepository interface {
	Get(id uuid.UUID) (Broker, bool, error)
	Create(broker Broker) (uuid.UUID, error)
	Update(broker Broker) error
	Delete(id uuid.UUID) error
	Exists(id uuid.UUID) (bool, error)
	ExistsByName(name string) (bool, error)
	GetAll() ([]Broker, error)
	GetAllEnabled() ([]Broker, error)

	SetImage(id uuid.UUID, imageId uuid.UUID) error
	HasImage(id uuid.UUID) (bool, error)
	DeleteImage(id uuid.UUID) error
}
