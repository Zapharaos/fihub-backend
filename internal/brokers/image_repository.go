package brokers

import "github.com/google/uuid"

// ImageRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Image
type ImageRepository interface {
	Create(image Image) error
	Get(brokerImageID uuid.UUID) (Image, bool, error)
	Update(image Image) error
	Delete(brokerImageID uuid.UUID) error
	Exists(brokerID uuid.UUID, brokerImageID uuid.UUID) (bool, error)
}
