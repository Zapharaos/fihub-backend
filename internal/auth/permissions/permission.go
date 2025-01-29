package permissions

import (
	"fmt"
	"github.com/google/uuid"
)

type Permission struct {
	Id          uuid.UUID `json:"id"`
	Value       string    `json:"value"`
	Scope       string    `json:"scope"`
	Description string    `json:"description"`
}

type Scope = string

const (
	AdminScope Scope = "admin"
	AllScope   Scope = "all"
)

var validScopes = []Scope{AdminScope, AllScope}

// IsValid checks if a permission is valid and has no missing mandatory fields
func (p Permission) IsValid() (bool, error) {
	if p.Value == "" {
		return false, fmt.Errorf("missing value")
	}
	if p.Scope == "" {
		return false, fmt.Errorf("missing scope")
	}
	// check if the scope is valid
	if !CheckScope(p.Scope) {
		return false, fmt.Errorf("invalid scope")
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

// CheckScope checks if a scope is valid
func CheckScope(scope Scope) bool {
	for _, s := range validScopes {
		if s == scope {
			return true
		}
	}
	return false
}
