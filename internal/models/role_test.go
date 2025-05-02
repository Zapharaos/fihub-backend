package models

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TestRole_IsValid tests the IsValid method of the Role struct
func TestRole_IsValid(t *testing.T) {
	// Define valid values
	validUUID := uuid.New()
	validName := strings.Repeat("a", NameMinLength)

	// Define test cases
	tests := []struct {
		name     string
		role     Role
		expected bool
		err      error
	}{
		{
			name:     "valid role",
			role:     Role{Id: validUUID, Name: validName},
			expected: true,
			err:      nil,
		},
		{
			name:     "empty name",
			role:     Role{Id: validUUID, Name: ""},
			expected: false,
			err:      ErrNameRequired,
		},
		{
			name:     "name too short",
			role:     Role{Id: validUUID, Name: strings.Repeat("a", NameMinLength-1)},
			expected: false,
			err:      ErrNameInvalid,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := tt.role.IsValid()
			assert.Equal(t, tt.expected, valid)
			assert.Equal(t, tt.err, err)
		})
	}
}

// TestRolesWithPermissions_GetUUIDs tests the GetUUIDs method of the RolesWithPermissions struct
func TestRolesWithPermissions_GetUUIDs(t *testing.T) {
	// Define valid values
	validUUID1 := uuid.New()
	validUUID2 := uuid.New()

	// Define test cases
	tests := []struct {
		name     string
		roles    RolesWithPermissions
		expected []uuid.UUID
	}{
		{
			name: "multiple roles",
			roles: RolesWithPermissions{
				{Role: Role{Id: validUUID1}},
				{Role: Role{Id: validUUID2}},
			},
			expected: []uuid.UUID{validUUID1, validUUID2},
		},
		{
			name:     "no roles",
			roles:    RolesWithPermissions{},
			expected: []uuid.UUID{},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuids := tt.roles.GetUUIDs()
			assert.Equal(t, tt.expected, uuids)
		})
	}
}

// TestRolesWithPermissions_HasPermission tests the HasPermission method
func TestRolesWithPermissions_HasPermission(t *testing.T) {

	// Define test cases
	tests := []struct {
		name       string               // Test case name
		roles      RolesWithPermissions // UserWithRoles instance to test
		permission string               // Permission to check
		expected   bool                 // Expected result
	}{
		{
			name:       "User has permission",
			roles:      createRolesWithPermission("admin.Users.read"),
			permission: "admin.Users.read",
			expected:   true,
		},
		{
			name:       "User does not have permission",
			roles:      createRolesWithPermission(""),
			permission: "admin.Users.read",
			expected:   false,
		},
		{
			name:       "User has wildcard permission",
			roles:      createRolesWithPermission("admin.Users.*"),
			permission: "admin.Users.read",
			expected:   true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.roles.HasPermission(tt.permission)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// createRolesWithPermission creates a RolesWithPermissions instance with a single permission
func createRolesWithPermission(permission string) RolesWithPermissions {
	return RolesWithPermissions{
		{
			Permissions: Permissions{
				{Value: permission},
			},
		},
	}
}

// TestRolePermissionsInput_IsValid tests the IsValid method of the RolePermissionsInput struct
func TestRolePermissionsInput_IsValid(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string               // Test case name
		perms    RolePermissionsInput // Permissions instance to test
		expected bool                 // Expected result
		err      error                // Expected error
	}{
		{
			name:     "Valid permissions",
			perms:    make(RolePermissionsInput, LimitMaxRolePermissions),
			expected: true,
			err:      nil,
		},
		{
			name:     "Empty permissions",
			perms:    make(RolePermissionsInput, 0),
			expected: true,
			err:      nil,
		},
		{
			name:     "Exceeds limit",
			perms:    make(RolePermissionsInput, LimitMaxRolePermissions+1),
			expected: false,
			err:      ErrLimitExceeded,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := tt.perms.IsValid()
			assert.Equal(t, tt.expected, valid)
			assert.Equal(t, tt.err, err)
		})
	}
}

// TestRolePermissionsInput_GetUUIDsAsStrings tests the GetUUIDsAsStrings method of the RolesWithPermissions struct
func TestRolePermissionsInput_GetUUIDsAsStrings(t *testing.T) {
	// Define valid values
	validUUID1 := uuid.New()
	validUUID2 := uuid.New()

	// Define test cases
	tests := []struct {
		name     string
		roles    RolePermissionsInput
		expected []string
	}{
		{
			name:     "nil roles",
			roles:    nil,
			expected: []string{},
		},
		{
			name:     "no roles",
			roles:    RolePermissionsInput{},
			expected: []string{},
		},
		{
			name:     "multiple roles",
			roles:    RolePermissionsInput{validUUID1, validUUID2},
			expected: []string{validUUID1.String(), validUUID2.String()},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuids := tt.roles.GetUUIDsAsStrings()
			assert.Equal(t, tt.expected, uuids)
		})
	}
}

// TestRolePermissionsInputFromUUIDs tests the RolePermissionsInputFromUUIDs function
func TestRolePermissionsInputFromUUIDs(t *testing.T) {
	// Define valid values
	validUUID1 := uuid.New()
	validUUID2 := uuid.New()

	// Define test cases
	tests := []struct {
		name     string
		roles    []string
		expected RolePermissionsInput
	}{
		{
			name:     "nil roles",
			roles:    nil,
			expected: RolePermissionsInput{},
		},
		{
			name:     "empty roles",
			roles:    []string{},
			expected: RolePermissionsInput{},
		},
		{
			name:     "bad role",
			roles:    []string{"invalid-uuid"},
			expected: RolePermissionsInput{},
		},
		{
			name:     "multiple roles",
			roles:    []string{validUUID1.String(), validUUID2.String()},
			expected: RolePermissionsInput{validUUID1, validUUID2},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuids := RolePermissionsInputFromUUIDs(tt.roles)
			assert.Equal(t, tt.expected, uuids)
		})
	}
}
