package models

import (
	"errors"
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

// ToUUIDs converts a slice of RolePermissionsInput to a slice of string uuids
func (r RolePermissionsInput) ToUUIDs() []string {
	if r == nil {
		return nil
	}

	permissions := make([]string, len(r))
	for i, perm := range r {
		permissions[i] = perm.String()
	}

	return permissions
}

// RolePermissionsInputFromUUIDs converts a slice of string uuids to RolePermissionsInput
func RolePermissionsInputFromUUIDs(s []string) RolePermissionsInput {
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
