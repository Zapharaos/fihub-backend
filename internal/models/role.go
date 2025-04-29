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
	NameMinLength           = 3
	LimitMaxRolePermissions = 250
)

// Role represents a role in the system
type Role struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Roles []Role

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

// HasPermission returns true if the User has the given permission.
// Wildcards (*) in permissions are supported.
func (r RolesWithPermissions) HasPermission(permission string) bool {
	for _, role := range r {
		for _, p := range role.Permissions {
			if p.Match(permission) {
				return true
			}
		}
	}
	return false
}

type RolePermissionsInput []uuid.UUID

// IsValid checks if a RolePermissionsInput is valid
func (rp RolePermissionsInput) IsValid() (bool, error) {
	if len(rp) > LimitMaxRolePermissions {
		return false, ErrLimitExceeded
	}
	return true, nil
}

// ToProtogenRole converts a Role to a protogen.Role
func (r Role) ToProtogenRole() *protogen.Role {
	return &protogen.Role{
		Id:   r.Id.String(),
		Name: r.Name,
	}
}

// ToProtogenRoles converts a slice of Roles to a slice of protogen.Role
func (r Roles) ToProtogenRoles() []*protogen.Role {
	roles := make([]*protogen.Role, 0)
	for _, role := range r {
		roles = append(roles, role.ToProtogenRole())
	}

	return roles
}

// ToProtogenString converts a slice of RolePermissionsInput to a slice of strings
func (r RolePermissionsInput) ToProtogenString() []string {
	if r == nil {
		return nil
	}

	permissions := make([]string, len(r))
	for i, perm := range r {
		permissions[i] = perm.String()
	}

	return permissions
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
func FromProtogenRole(r *protogen.Role) Role {
	if r == nil {
		return Role{}
	}

	id, err := uuid.Parse(r.GetId())
	if err != nil {
		return Role{}
	}

	return Role{
		Id:   id,
		Name: r.GetName(),
	}
}

// FromProtogenRoles converts a slice of protogen.Role to Roles
func FromProtogenRoles(roles []*protogen.Role) Roles {
	if roles == nil {
		return Roles{}
	}

	result := make(Roles, len(roles))
	for i, role := range roles {
		r := FromProtogenRole(role)
		result[i] = r
	}

	return result
}

// FromProtogenRolePermissionsInput converts a slice of string to RolePermissionsInput
func FromProtogenRolePermissionsInput(s []string) RolePermissionsInput {
	if s == nil {
		return RolePermissionsInput{}
	}

	permissions := make(RolePermissionsInput, len(s))
	for i, perm := range s {
		id, err := uuid.Parse(perm)
		if err != nil {
			return RolePermissionsInput{}
		}
		permissions[i] = id
	}

	return permissions
}

// FromProtogenRoleWithPermissions converts a protogen.RoleWithPermissions to a RoleWithPermissions
func FromProtogenRoleWithPermissions(r *protogen.RoleWithPermissions) RoleWithPermissions {
	if r == nil {
		return RoleWithPermissions{}
	}

	return RoleWithPermissions{
		Role:        FromProtogenRole(r.GetRole()),
		Permissions: FromProtogenPermissions(r.GetPermissions()),
	}
}

// FromProtogenRolesWithPermissions converts a slice of protogen.RoleWithPermissions to a slice of RoleWithPermissions
func FromProtogenRolesWithPermissions(roles []*protogen.RoleWithPermissions) RolesWithPermissions {
	if roles == nil {
		return RolesWithPermissions{}
	}

	result := make(RolesWithPermissions, len(roles))
	for i, role := range roles {
		r := FromProtogenRoleWithPermissions(role)
		result[i] = r
	}

	return result
}
