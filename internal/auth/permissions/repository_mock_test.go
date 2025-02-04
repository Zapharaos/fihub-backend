package permissions

import "github.com/google/uuid"

// MockRepository represents a mock Repository
type MockRepository struct {
	ID    uuid.UUID
	Found bool
	Error error
	Perm  Permission
	Perms []Permission
}

// NewMockRepository creates a new MockRepository of the Repository interface
func NewMockRepository() Repository {
	r := MockRepository{}
	var repo Repository = &r
	return repo
}

func (m MockRepository) Get(_ uuid.UUID) (Permission, bool, error) {
	return m.Perm, m.Found, m.Error
}

func (m MockRepository) Create(_ Permission) (uuid.UUID, error) {
	return m.ID, m.Error
}

func (m MockRepository) Update(_ Permission) error {
	return m.Error
}

func (m MockRepository) Delete(_ uuid.UUID) error {
	return m.Error
}

func (m MockRepository) GetAll() ([]Permission, error) {
	return m.Perms, m.Error
}

func (m MockRepository) GetAllByRoleId(_ uuid.UUID) ([]Permission, error) {
	return m.Perms, m.Error
}

func (m MockRepository) GetAllByRoleIds(_ []uuid.UUID) ([]Permission, error) {
	return m.Perms, m.Error
}

func (m MockRepository) GetAllForUser(_ uuid.UUID) ([]Permission, error) {
	return m.Perms, m.Error
}
