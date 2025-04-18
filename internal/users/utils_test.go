package users_test

import (
	"errors"
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/users"
	"github.com/Zapharaos/fihub-backend/internal/users/permissions"
	"github.com/Zapharaos/fihub-backend/internal/users/roles"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

// TestLoadUserRoles tests the LoadUserRoles function
func TestLoadUserRoles(t *testing.T) {
	// Define test data
	role := models.Role{Id: uuid.New(), Name: "admin"}
	perm := models.Permission{Id: uuid.New(), Value: "read", Scope: "scope"}
	perms := []models.Permission{perm, perm, perm}

	// Define test cases
	tests := []struct {
		name      string                      // Test case name
		mockSetup func(*gomock.Controller)    // Mock setup function
		expected  models.RolesWithPermissions // Expected result
		error     error                       // Expected error
	}{
		{
			name: "can't retrieve roles",
			mockSetup: func(ctrl *gomock.Controller) {
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return(nil, errors.New("error"))
				roles.ReplaceGlobals(r)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			error: fmt.Errorf("error"),
		},
		{
			name: "user has no roles",
			mockSetup: func(ctrl *gomock.Controller) {
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return([]models.Role{}, nil)
				roles.ReplaceGlobals(r)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expected: make(models.RolesWithPermissions, 0),
			error:    nil,
		},
		{
			name: "user has one role but can't retrieve permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return([]models.Role{role}, nil)
				roles.ReplaceGlobals(r)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Return(nil, errors.New("error"))
				permissions.ReplaceGlobals(p)
			},
			expected: models.RolesWithPermissions(nil),
			error:    fmt.Errorf("error"),
		},
		{
			name: "user has one role without permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return([]models.Role{role}, nil)
				roles.ReplaceGlobals(r)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Return([]models.Permission{}, nil)
				permissions.ReplaceGlobals(p)
			},
			expected: models.RolesWithPermissions{
				models.RoleWithPermissions{Role: role, Permissions: make([]models.Permission, 0)},
			},
			error: nil,
		},
		{
			name: "user has multiple roles with multiple permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return([]models.Role{role, role, role}, nil)
				roles.ReplaceGlobals(r)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Return(perms, nil).Times(3)
				permissions.ReplaceGlobals(p)
			},
			expected: models.RolesWithPermissions{
				models.RoleWithPermissions{Role: role, Permissions: perms},
				models.RoleWithPermissions{Role: role, Permissions: perms},
				models.RoleWithPermissions{Role: role, Permissions: perms},
			},
			error: nil,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tt.mockSetup(ctrl)

			// Call the function
			result, err := users.LoadUserRoles(uuid.Nil)

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
	user := models.User{ID: uuid.New()}
	role := models.Role{Id: uuid.New(), Name: "admin"}
	perm := models.Permission{Id: uuid.New(), Value: "read", Scope: "scope"}
	perms := []models.Permission{perm, perm, perm}
	roleWP := models.RoleWithPermissions{Role: role, Permissions: perms}

	// Define test cases
	tests := []struct {
		name      string                   // Test case name
		mockSetup func(*gomock.Controller) // Mock setup function
		expected  *models.UserWithRoles    // Expected result
		found     bool                     // Expected found
		error     error                    // Expected error
	}{
		{
			name: "can't retrieve user",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(models.User{}, false, errors.New("error"))
				users.ReplaceGlobals(u)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Times(0)
				roles.ReplaceGlobals(r)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expected: nil,
			found:    false,
			error:    fmt.Errorf("error"),
		},
		{
			name: "user not found",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(models.User{}, false, nil)
				users.ReplaceGlobals(u)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Times(0)
				roles.ReplaceGlobals(r)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expected: nil,
			found:    false,
			error:    users.ErrorUserNotFound,
		},
		{
			name: "can't retrieve roles",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(user, true, nil)
				users.ReplaceGlobals(u)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return([]models.Role{}, errors.New("error"))
				roles.ReplaceGlobals(r)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expected: (*models.UserWithRoles)(nil),
			found:    false,
			error:    fmt.Errorf("error"),
		},
		{
			name: "user has multiple roles",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(user, true, nil)
				users.ReplaceGlobals(u)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return([]models.Role{role, role, role}, nil)
				roles.ReplaceGlobals(r)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Return(perms, nil).Times(3)
				permissions.ReplaceGlobals(p)
			},
			expected: &models.UserWithRoles{
				User: user, Roles: models.RolesWithPermissions{
					roleWP, roleWP, roleWP,
				},
			},
			found: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tt.mockSetup(ctrl)

			// Call the function
			result, found, err := users.LoadFullUser(uuid.Nil)

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
