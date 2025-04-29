package models

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
	"regexp"
	"strings"
)

var (
	ErrValueRequired = errors.New("value-required")
	ErrScopeRequired = errors.New("scope-required")
	ErrScopeInvalid  = errors.New("scope-invalid")
	ErrLimitExceeded = errors.New("permissions-limit-exceeded")
)

type Permission struct {
	Id          uuid.UUID `json:"id"`
	Value       string    `json:"value"`
	Scope       string    `json:"scope"`
	Description string    `json:"description"`
}

type Permissions []Permission

type Scope = string

const (
	AdminScope Scope = "admin"
	AllScope   Scope = "all"
)

var validScopes = []Scope{AdminScope, AllScope}

// IsValid checks if a permission is valid and has no missing mandatory fields
func (p Permission) IsValid() (bool, error) {
	if p.Value == "" {
		return false, ErrValueRequired
	}
	if p.Scope == "" {
		return false, ErrScopeRequired
	}
	// check if the scope is valid
	if !CheckScope(p.Scope) {
		return false, ErrScopeInvalid
	}
	return true, nil
}

// HasScope checks if a permission has a specific scope
func (p Permission) HasScope(scope Scope) bool {
	return p.Scope == scope
}

// GetScopes returns the list of valid scopes
func (p Permission) GetScopes() []Scope {
	return validScopes
}

// Match checks if the given permission matches the pattern,
// supporting wildcards (*) in the pattern using regular expressions.
func (p Permission) Match(permission string) bool {
	if p.Value == "*" {
		return true
	}

	// Escape special characters in the pattern
	pattern := regexp.QuoteMeta(p.Value)

	// Replace wildcard (*) with a regex-friendly wildcard (.*)
	pattern = strings.ReplaceAll(pattern, "\\*", ".*")

	// Compile the regular expression
	re := regexp.MustCompile("^" + pattern + "$")

	// Check if the permission matches the compiled regex
	return re.MatchString(permission)
}

// ToProtogenPermission converts a Permission to a protogen.Permission
func (p Permission) ToProtogenPermission() *protogen.Permission {
	return &protogen.Permission{
		Id:          p.Id.String(),
		Value:       p.Value,
		Scope:       p.Scope,
		Description: p.Description,
	}
}

// FromProtogenPermission converts a protogen.Permission to a Permission
func FromProtogenPermission(p *protogen.Permission) (Permission, error) {
	if p == nil {
		return Permission{}, errors.New("permission is nil")
	}

	id, err := uuid.Parse(p.GetId())
	if err != nil {
		return Permission{}, err
	}

	return Permission{
		Id:          id,
		Value:       p.GetValue(),
		Scope:       p.GetScope(),
		Description: p.GetDescription(),
	}, nil
}

// ToProtogenPermissions converts a slice of Permission to a slice of protogen.Permission
func (p Permissions) ToProtogenPermissions() []*protogen.Permission {
	if p == nil {
		return []*protogen.Permission{}
	}

	permissions := make([]*protogen.Permission, len(p))
	for i, perm := range p {
		permissions[i] = perm.ToProtogenPermission()
	}

	return permissions
}

// ToProtogenPermissionsUuidInput converts a slice of Permission to a slice of string
func (p Permissions) ToProtogenPermissionsUuidInput() []string {
	if p == nil {
		return []string{}
	}

	permissions := make([]string, len(p))
	for i, perm := range p {
		permissions[i] = perm.Id.String()
	}

	return permissions
}

// FromProtogenPermissions converts a slice of protogen.Permission to a slice of Permission
func FromProtogenPermissions(p []*protogen.Permission) Permissions {
	if p == nil {
		return Permissions{}
	}

	permissions := make(Permissions, len(p))
	for i, protoPerm := range p {
		perm, err := FromProtogenPermission(protoPerm)
		if err != nil {
			continue
		}
		permissions[i] = perm
	}

	return permissions
}

// CheckScope checks if a scope is valid
func CheckScope(scope Scope) bool {
	for _, s := range validScopes {
		if s == scope {
			return true
		}
	}
	return false
}

// GetUUIDs returns the list of UUIDs for the permissions
func (p Permissions) GetUUIDs() []uuid.UUID {
	uuids := make([]uuid.UUID, 0)
	for _, perm := range p {
		uuids = append(uuids, perm.Id)
	}
	return uuids
}
