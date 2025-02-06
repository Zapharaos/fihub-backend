package mocks

import (
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/google/uuid"
)

// RolesRepository represents a mocks roles.Repository
type RolesRepository struct {
	ID                   uuid.UUID
	Found                bool
	Error                error
	Role                 roles.Role
	Roles                []roles.Role
	RoleWithPermissions  roles.RoleWithPermissions
	RolesWithPermissions roles.RolesWithPermissions
}

// NewRolesRepository creates a new RolesRepository of the roles.Repository interface
func NewRolesRepository(r RolesRepository) roles.Repository {
	var repo roles.Repository = &r
	return repo
}

func (m RolesRepository) Get(_ uuid.UUID) (roles.Role, bool, error) {
	return m.Role, m.Found, m.Error
}

func (m RolesRepository) GetByName(_ string) (roles.Role, bool, error) {
	return m.Role, m.Found, m.Error
}

func (m RolesRepository) GetWithPermissions(_ uuid.UUID) (roles.RoleWithPermissions, bool, error) {
	return m.RoleWithPermissions, m.Found, m.Error
}

func (m RolesRepository) GetAll() ([]roles.Role, error) {
	return m.Roles, m.Error
}

func (m RolesRepository) GetAllWithPermissions() (roles.RolesWithPermissions, error) {
	return m.RolesWithPermissions, m.Error
}

func (m RolesRepository) Create(_ roles.Role, _ []uuid.UUID) (uuid.UUID, error) {
	return m.ID, m.Error
}

func (m RolesRepository) Update(_ roles.Role, _ []uuid.UUID) error {
	return m.Error
}

func (m RolesRepository) Delete(_ uuid.UUID) error {
	return m.Error
}

func (m RolesRepository) GetRolesByUserId(_ uuid.UUID) ([]roles.Role, error) {
	return m.Roles, m.Error
}

func (m RolesRepository) SetRolePermissions(_ uuid.UUID, _ []uuid.UUID) error {
	return m.Error
}
