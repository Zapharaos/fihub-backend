package repositories

// Repository is a struct that contains all the repositories
type Repository struct {
	broker BrokerRepository
	user   UserRepository
	image  ImageRepository
}

// NewRepository returns a new instance of Repository
func NewRepository(broker BrokerRepository, user UserRepository, image ImageRepository) Repository {
	return Repository{
		broker: broker,
		user:   user,
		image:  image,
	}
}

// B is used to access the BrokerRepository singleton
func (r Repository) B() BrokerRepository {
	return r.broker
}

// U is used to access the UserRepository singleton
func (r Repository) U() UserRepository {
	return r.user
}

// I is used to access the ImageRepository singleton
func (r Repository) I() ImageRepository {
	return r.image
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
