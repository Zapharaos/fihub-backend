package brokers

// BrokerRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on Users
type BrokerRepository interface {
	GetAll() ([]Broker, error)
}
