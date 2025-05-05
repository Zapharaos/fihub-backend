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
	Id   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}

// Roles represents a list of Role
type Roles []Role

// RoleWithPermissions represents a Role with its permissions.Permissions
type RoleWithPermissions struct {
	Role
	Permissions Permissions `json:"permissions"`
}

// RolesWithPermissions represents a list of RoleWithPermissions
type RolesWithPermissions []RoleWithPermissions

// RolePermissionsInput represents a list of Permission UUIDs at input
type RolePermissionsInput []uuid.UUID

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

// IsValid checks if a RolePermissionsInput is valid
func (rpi RolePermissionsInput) IsValid() (bool, error) {
	if len(rpi) > LimitMaxRolePermissions {
		return false, ErrLimitExceeded
	}
	return true, nil
}

// GetUUIDsAsStrings converts a slice of RolePermissionsInput to a slice of string uuids
func (rpi RolePermissionsInput) GetUUIDsAsStrings() []string {
	if rpi == nil {
		return nil
	}

	permissions := make([]string, len(rpi))
	for i, perm := range rpi {
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
