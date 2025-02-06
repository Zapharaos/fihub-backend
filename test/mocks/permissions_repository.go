package mocks

import (
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/google/uuid"
)

// PermissionsRepository represents a mocks permissions.Repository
type PermissionsRepository struct {
	ID    uuid.UUID
	Found bool
	Error error
	Perm  permissions.Permission
	Perms []permissions.Permission
}

// NewPermissionsRepository creates a new PermissionsRepository of the permissions.Repository interface
func NewPermissionsRepository(r PermissionsRepository) permissions.Repository {
	var repo permissions.Repository = &r
	return repo
}

func (m PermissionsRepository) Get(_ uuid.UUID) (permissions.Permission, bool, error) {
	return m.Perm, m.Found, m.Error
}

func (m PermissionsRepository) Create(_ permissions.Permission) (uuid.UUID, error) {
	return m.ID, m.Error
}

func (m PermissionsRepository) Update(_ permissions.Permission) error {
	return m.Error
}

func (m PermissionsRepository) Delete(_ uuid.UUID) error {
	return m.Error
}

func (m PermissionsRepository) GetAll() ([]permissions.Permission, error) {
	return m.Perms, m.Error
}

func (m PermissionsRepository) GetAllByRoleId(_ uuid.UUID) ([]permissions.Permission, error) {
	return m.Perms, m.Error
}

func (m PermissionsRepository) GetAllByRoleIds(_ []uuid.UUID) ([]permissions.Permission, error) {
	return m.Perms, m.Error
}

func (m PermissionsRepository) GetAllForUser(_ uuid.UUID) ([]permissions.Permission, error) {
	return m.Perms, m.Error
}
