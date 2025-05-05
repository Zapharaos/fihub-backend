package models

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestPermission_IsValid tests the IsValid method of the Permission struct
func TestPermission_IsValid(t *testing.T) {
	// Define valid values
	validValue := "read"

	// Define test cases
	tests := []struct {
		name     string     // Test case name
		perm     Permission // Permission instance to test
		expected bool       // Expected result
		err      error      // Expected error
	}{
		{
			name:     "Valid permission",
			perm:     Permission{Value: validValue, Scope: AdminScope},
			expected: true,
			err:      nil,
		},
		{
			name:     "Missing value",
			perm:     Permission{Scope: AdminScope},
			expected: false,
			err:      ErrValueRequired,
		},
		{
			name:     "Missing scope",
			perm:     Permission{Value: validValue},
			expected: false,
			err:      ErrScopeRequired,
		},
		{
			name:     "Invalid scope",
			perm:     Permission{Value: validValue, Scope: "invalid"},
			expected: false,
			err:      ErrScopeInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := tt.perm.IsValid()
			assert.Equal(t, tt.expected, valid)
			assert.Equal(t, tt.err, err)
		})
	}
}

// TestPermission_HasScope tests the HasScope method of the Permission struct
func TestPermission_HasScope(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string     // Test case name
		perm     Permission // Permission instance to test
		scope    Scope      // Scope to check
		expected bool       // Expected result
	}{
		{
			name:     "Has admin scope",
			perm:     Permission{Scope: AdminScope},
			scope:    AdminScope,
			expected: true,
		},
		{
			name:     "Does not have admin scope",
			perm:     Permission{Scope: AllScope},
			scope:    AdminScope,
			expected: false,
		},
		{
			name:     "Has all scope",
			perm:     Permission{Scope: AllScope},
			scope:    AllScope,
			expected: true,
		},
		{
			name:     "Does not have all scope",
			perm:     Permission{Scope: AdminScope},
			scope:    AllScope,
			expected: false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.perm.HasScope(tt.scope)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPermission_GetScopes tests the GetScopes method of the Permission struct
func TestPermission_GetScopes(t *testing.T) {
	perm := Permission{}
	expected := validScopes
	result := perm.GetScopes()
	assert.Equal(t, expected, result)
}

// TestPermission_Match tests the Match method of the Permission struct
func TestPermission_Match(t *testing.T) {
	// Define test cases
	tests := []struct {
		name       string     // Test case name
		perm       Permission // Permission instance to test
		permission string     // Permission to match
		expected   bool       // Expected result
	}{
		{
			name:       "Exact match",
			perm:       Permission{Value: "admin.users.read"},
			permission: "admin.users.read",
			expected:   true,
		},
		{
			name:       "Wildcard match",
			perm:       Permission{Value: "admin.users.*"},
			permission: "admin.users.read",
			expected:   true,
		},
		{
			name:       "No match",
			perm:       Permission{Value: "admin.users.write"},
			permission: "admin.users.read",
			expected:   false,
		},
		{
			name:       "Global wildcard match",
			perm:       Permission{Value: "*"},
			permission: "admin.users.read",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.perm.Match(tt.permission)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPermissions_GetUUIDs tests the GetUUIDs method of the Permissions struct
func TestPermissions_GetUUIDs(t *testing.T) {
	// Define valid values
	validUUID1 := uuid.New()
	validUUID2 := uuid.New()

	// Define test cases
	tests := []struct {
		name     string      // Test case name
		perms    Permissions // Permissions instance to test
		expected []uuid.UUID // Expected result
	}{
		{
			name: "multiple permissions",
			perms: Permissions{
				{Id: validUUID1},
				{Id: validUUID2},
			},
			expected: []uuid.UUID{validUUID1, validUUID2},
		},
		{
			name:     "no permissions",
			perms:    Permissions{},
			expected: []uuid.UUID{},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuids := tt.perms.GetUUIDs()
			assert.Equal(t, tt.expected, uuids)
		})
	}
}

// TestPermissions_GetUUIDsAsStrings tests the GetUUIDsAsStrings method of the Permissions struct
func TestPermissions_GetUUIDsAsStrings(t *testing.T) {
	// Define valid values
	validUUID1 := uuid.New()
	validUUID2 := uuid.New()

	// Define test cases
	tests := []struct {
		name     string      // Test case name
		perms    Permissions // Permissions instance to test
		expected []string    // Expected result
	}{
		{
			name:     "nil permissions",
			perms:    nil,
			expected: []string{},
		},
		{
			name:     "no permissions",
			perms:    Permissions{},
			expected: []string{},
		},
		{
			name: "multiple permissions",
			perms: Permissions{
				{Id: validUUID1},
				{Id: validUUID2},
			},
			expected: []string{validUUID1.String(), validUUID2.String()},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuids := tt.perms.GetUUIDsAsStrings()
			assert.Equal(t, tt.expected, uuids)
		})
	}
}
