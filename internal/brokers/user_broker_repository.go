package brokers

import "github.com/google/uuid"

// UserBrokerRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Users
type UserBrokerRepository interface {
	Create(userBroker UserBroker) error
	Delete(userBroker UserBroker) error
	Exists(userBroker UserBroker) (bool, error)
	GetAll(userID uuid.UUID) ([]UserBroker, error)
}
