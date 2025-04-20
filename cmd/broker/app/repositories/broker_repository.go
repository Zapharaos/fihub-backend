package repositories

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
)

// BrokerRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Broker
type BrokerRepository interface {
	Get(id uuid.UUID) (models.Broker, bool, error)
	Create(broker models.Broker) (uuid.UUID, error)
	Update(broker models.Broker) error
	Delete(id uuid.UUID) error
	Exists(id uuid.UUID) (bool, error)
	ExistsByName(name string) (bool, error)
	GetAll() ([]models.Broker, error)
	GetAllEnabled() ([]models.Broker, error)

	SetImage(id uuid.UUID, imageId uuid.UUID) error
	HasImage(id uuid.UUID) (bool, error)
	DeleteImage(id uuid.UUID) error
}
