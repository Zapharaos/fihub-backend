package brokers

import "github.com/google/uuid"

// UserRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on User
type UserRepository interface {
	Create(userBroker User) error
	Delete(userBroker User) error
	Exists(userBroker User) (bool, error)
	GetAll(userID uuid.UUID) ([]User, error)
}
