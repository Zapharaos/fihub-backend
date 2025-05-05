package repositories

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
)

// ImageRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on BrokerImage
type ImageRepository interface {
	Create(image models.BrokerImage) error
	Get(brokerImageID uuid.UUID) (models.BrokerImage, bool, error)
	Update(image models.BrokerImage) error
	Delete(brokerImageID uuid.UUID) error
	Exists(brokerID uuid.UUID, brokerImageID uuid.UUID) (bool, error)
}
