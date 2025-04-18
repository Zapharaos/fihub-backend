package brokers

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
)

// UserRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on BrokerUser
type UserRepository interface {
	Create(userBroker models.BrokerUser) error
	Delete(userBroker models.BrokerUser) error
	Exists(userBroker models.BrokerUser) (bool, error)
	GetAll(userID uuid.UUID) ([]models.BrokerUser, error)
}
