package repositories

// Repository is a struct that contains all the repositories
type Repository struct {
	role       RoleRepository
	permission PermissionRepository
}

// NewRepository returns a new instance of Repository
func NewRepository(role RoleRepository, permission PermissionRepository) Repository {
	return Repository{
		role:       role,
		permission: permission,
	}
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
