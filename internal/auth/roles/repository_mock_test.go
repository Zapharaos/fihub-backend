package roles

import "github.com/google/uuid"

// MockRepository represents a mock Repository
type MockRepository struct {
	ID                   uuid.UUID
	Found                bool
	Error                error
	Role                 Role
	Roles                []Role
	RoleWithPermissions  RoleWithPermissions
	RolesWithPermissions RolesWithPermissions
}

// NewMockRepository creates a new MockRepository of the Repository interface
func NewMockRepository() Repository {
	r := MockRepository{}
	var repo Repository = &r
	return repo
}

func (m MockRepository) Get(_ uuid.UUID) (Role, bool, error) {
	return m.Role, m.Found, m.Error
}

func (m MockRepository) GetByName(_ string) (Role, bool, error) {
	return m.Role, m.Found, m.Error
}

func (m MockRepository) GetWithPermissions(_ uuid.UUID) (RoleWithPermissions, bool, error) {
	return m.RoleWithPermissions, m.Found, m.Error
}

func (m MockRepository) GetAll() ([]Role, error) {
	return m.Roles, m.Error
}

func (m MockRepository) GetAllWithPermissions() (RolesWithPermissions, error) {
	return m.RolesWithPermissions, m.Error
}

func (m MockRepository) Create(_ Role, _ []uuid.UUID) (uuid.UUID, error) {
	return m.ID, m.Error
}

func (m MockRepository) Update(_ Role, _ []uuid.UUID) error {
	return m.Error
}

func (m MockRepository) Delete(_ uuid.UUID) error {
	return m.Error
}

func (m MockRepository) GetRolesByUserId(_ uuid.UUID) ([]Role, error) {
	return m.Roles, m.Error
}

func (m MockRepository) SetRolePermissions(_ uuid.UUID, _ []uuid.UUID) error {
	return m.Error
}
