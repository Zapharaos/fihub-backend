package auth

import (
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/test/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestLoadUserRoles tests the LoadUserRoles function
func TestLoadUserRoles(t *testing.T) {
	// Define test data
	role := roles.Role{Id: uuid.New(), Name: "admin"}
	perm := permissions.Permission{Id: uuid.New(), Value: "read", Scope: "scope"}
	perms := []permissions.Permission{perm, perm, perm}

	// Define test cases
	tests := []struct {
		name        string                     // Test case name
		roles       mock.RolesRepository       // Roles repository mock
		permissions mock.PermissionsRepository // Permissions repository mock
		expected    roles.RolesWithPermissions // Expected result
		error       error                      // Expected error
	}{
		{
			name: "can't retrieve roles",
			roles: mock.RolesRepository{
				Error: fmt.Errorf("error"),
			},
			expected: roles.RolesWithPermissions(nil),
			error:    fmt.Errorf("error"),
		},
		{
			name: "user has no roles",
			roles: mock.RolesRepository{
				Roles: make([]roles.Role, 0),
			},
			expected: make(roles.RolesWithPermissions, 0),
			error:    nil,
		},
		{
			name: "user has one role but can't retrieve permissions",
			roles: mock.RolesRepository{
				Roles: []roles.Role{role},
			},
			permissions: mock.PermissionsRepository{
				Error: fmt.Errorf("error"),
			},
			expected: roles.RolesWithPermissions(nil),
			error:    fmt.Errorf("error"),
		},
		{
			name: "user has one role without permissions",
			roles: mock.RolesRepository{
				Roles: []roles.Role{role},
			},
			permissions: mock.PermissionsRepository{
				Perms: make([]permissions.Permission, 0),
			},
			expected: roles.RolesWithPermissions{
				roles.RoleWithPermissions{Role: role, Permissions: make([]permissions.Permission, 0)},
			},
			error: nil,
		},
		{
			name: "user has multiple roles with multiple permissions",
			roles: mock.RolesRepository{
				Roles: []roles.Role{role, role, role},
			},
			permissions: mock.PermissionsRepository{
				Perms: perms,
			},
			expected: roles.RolesWithPermissions{
				roles.RoleWithPermissions{Role: role, Permissions: perms},
				roles.RoleWithPermissions{Role: role, Permissions: perms},
				roles.RoleWithPermissions{Role: role, Permissions: perms},
			},
			error: nil,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the repositories
			roles.ReplaceGlobals(mock.NewRolesRepository(tt.roles))
			permissions.ReplaceGlobals(mock.NewPermissionsRepository(tt.permissions))

			// Call the function
			result, err := LoadUserRoles(uuid.Nil)

			// Assert results
			assert.Equal(t, tt.expected, result)
			if tt.error != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.error, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestLoadFullUser tests the LoadFullUser function
func TestLoadFullUser(t *testing.T) {
	// Define test data
	user := users.User{ID: uuid.New()}
	role := roles.Role{Id: uuid.New(), Name: "admin"}
	perm := permissions.Permission{Id: uuid.New(), Value: "read", Scope: "scope"}
	perms := []permissions.Permission{perm, perm, perm}
	roleWP := roles.RoleWithPermissions{Role: role, Permissions: perms}

	// Define test cases
	tests := []struct {
		name        string                     // Test case name
		users       mock.UsersRepository       // Users repository mock
		roles       mock.RolesRepository       // Roles repository mock
		permissions mock.PermissionsRepository // Permissions repository mock
		expected    *users.UserWithRoles       // Expected result
		found       bool                       // Expected found
		error       error                      // Expected error
	}{
		{
			name: "can't retrieve user",
			users: mock.UsersRepository{
				Err:   fmt.Errorf("error"),
				Found: false,
			},
			expected: nil,
			found:    false,
			error:    fmt.Errorf("error"),
		},
		{
			name: "user not found",
			users: mock.UsersRepository{
				Found: false,
			},
			expected: nil,
			found:    false,
			error:    ErrorUserNotFound,
		},
		{
			name: "can't retrieve roles",
			users: mock.UsersRepository{
				User:  user,
				Found: true,
			},
			roles: mock.RolesRepository{
				Error: fmt.Errorf("error"),
			},
			expected: (*users.UserWithRoles)(nil),
			found:    false,
			error:    fmt.Errorf("error"),
		},
		{
			name: "user has multiple roles",
			users: mock.UsersRepository{
				User:  user,
				Found: true,
			},
			roles: mock.RolesRepository{
				Roles: []roles.Role{role, role, role},
			},
			permissions: mock.PermissionsRepository{
				Perms: perms,
			},
			expected: &users.UserWithRoles{
				User: user, Roles: roles.RolesWithPermissions{
					roleWP, roleWP, roleWP,
				},
			},
			found: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the repositories
			users.ReplaceGlobals(mock.NewUsersRepository(tt.users))
			roles.ReplaceGlobals(mock.NewRolesRepository(tt.roles))
			permissions.ReplaceGlobals(mock.NewPermissionsRepository(tt.permissions))

			// Call the function
			result, found, err := LoadFullUser(uuid.Nil)

			// Assert results
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.found, found)
			if tt.error != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.error, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
