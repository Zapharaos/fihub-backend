package repositories

// Repository is a struct that contains all the repositories
type Repository struct {
	user       UserRepository
	role       RoleRepository
	permission PermissionRepository
}

// NewRepository returns a new instance of Repository
func NewRepository(user UserRepository, role RoleRepository, permission PermissionRepository) Repository {
	return Repository{
		user:       user,
		role:       role,
		permission: permission,
	}
}

// U is used to access the UserRepository singleton
func (r Repository) U() UserRepository {
	return r.user
}

// R is used to access the RoleRepository singleton
func (r Repository) R() RoleRepository {
	return r.role
}

// P is used to access the PermissionRepository singleton
func (r Repository) P() PermissionRepository {
	return r.permission
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
