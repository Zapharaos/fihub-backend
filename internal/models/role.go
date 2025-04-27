package models

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
)

var (
	ErrNameRequired = errors.New("name-required")
	ErrNameInvalid  = errors.New("name-invalid")
)

const (
	NameMinLength = 3
)

// Role represents a role in the system
type Role struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// IsValid checks if a Role is valid
func (r Role) IsValid() (bool, error) {
	if r.Name == "" {
		return false, ErrNameRequired
	}
	if len(r.Name) < NameMinLength {
		return false, ErrNameInvalid
	}
	return true, nil
}

// RoleWithPermissions represents a Role with its permissions.Permissions
type RoleWithPermissions struct {
	Role
	Permissions Permissions `json:"permissions"`
}

// RolesWithPermissions represents a list of RoleWithPermissions
type RolesWithPermissions []RoleWithPermissions

// GetUUIDs returns the list of UUIDs for the RolesWithPermissions
func (r RolesWithPermissions) GetUUIDs() []uuid.UUID {
	uuids := make([]uuid.UUID, 0)
	for _, role := range r {
		uuids = append(uuids, role.Id)
	}
	return uuids
}

// ToProtogenRole converts a Role to a protogen.Role
func (r Role) ToProtogenRole() *protogen.Role {
	return &protogen.Role{
		Id:   r.Id.String(),
		Name: r.Name,
	}
}

// ToProtogenRoleWithPermissions converts a RoleWithPermissions to a protogen.RoleWithPermissions
func (r RoleWithPermissions) ToProtogenRoleWithPermissions() *protogen.RoleWithPermissions {
	return &protogen.RoleWithPermissions{
		Role:        r.Role.ToProtogenRole(),
		Permissions: r.Permissions.ToProtogenPermissions(),
	}
}

// ToProtogenRolesWithPermissions converts a slice of RolesWithPermissions to a slice of protogen.RoleWithPermissions
func (r RolesWithPermissions) ToProtogenRolesWithPermissions() []*protogen.RoleWithPermissions {
	if r == nil {
		return nil
	}

	roles := make([]*protogen.RoleWithPermissions, len(r))
	for i, role := range r {
		roles[i] = role.ToProtogenRoleWithPermissions()
	}

	return roles
}

// FromProtogenRole converts a protogen.Role to a Role
func FromProtogenRole(r *protogen.Role) (Role, error) {
	if r == nil {
		return Role{}, errors.New("role is nil")
	}

	id, err := uuid.Parse(r.GetId())
	if err != nil {
		return Role{}, err
	}

	return Role{
		Id:   id,
		Name: r.GetName(),
	}, nil
}

// FromProtogenRoleWithPermissions converts a protogen.RoleWithPermissions to a RoleWithPermissions
func FromProtogenRoleWithPermissions(r *protogen.RoleWithPermissions) (RoleWithPermissions, error) {
	if r == nil {
		return RoleWithPermissions{}, errors.New("role is nil")
	}

	role, err := FromProtogenRole(r.GetRole())
	if err != nil {
		return RoleWithPermissions{}, err
	}

	permissions, err := FromProtogenPermissions(r.GetPermissions())
	if err != nil {
		return RoleWithPermissions{
			Role: role,
		}, err
	}

	return RoleWithPermissions{
		Role:        role,
		Permissions: permissions,
	}, nil
}

// FromProtogenRolesWithPermissions converts a slice of protogen.RoleWithPermissions to a slice of RoleWithPermissions
func FromProtogenRolesWithPermissions(roles []*protogen.RoleWithPermissions) (RolesWithPermissions, error) {
	if roles == nil {
		return nil, errors.New("roles is nil")
	}

	result := make(RolesWithPermissions, len(roles))
	for i, role := range roles {
		r, err := FromProtogenRoleWithPermissions(role)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}

	return result, nil
}
