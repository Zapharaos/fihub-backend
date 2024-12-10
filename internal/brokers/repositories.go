package brokers

// Repository stores the different repositories for events, parameters, users, roles and items
type Repository struct {
	broker BrokerRepository
}

// NewRepository returns a new instance of Repository
func NewRepository(broker BrokerRepository) Repository {
	return Repository{
		broker: broker,
	}
}

// B is used to access the broker repository singleton
func (r Repository) B() BrokerRepository {
	return r.broker
}

// R is used to access the global repository singleton
var _globalRepository Repository

// R is used to access the global repository singleton
func R() Repository {
	return _globalRepository
}

// ReplaceGlobals affect a new repository to the global repository singleton
func ReplaceGlobals(repository Repository) func() {
	prev := _globalRepository
	_globalRepository = repository
	return func() { ReplaceGlobals(prev) }
}
