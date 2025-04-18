package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/users"
	"github.com/Zapharaos/fihub-backend/internal/users/permissions"
	"github.com/Zapharaos/fihub-backend/internal/users/roles"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreateRole tests the CreateRole handler
func TestCreateRole(t *testing.T) {
	// Declare the data
	pRoleCreate := "admin.roles.create"
	pRolePermUpdate := "admin.roles.permissions.update"
	permission := models.Permission{
		Id: uuid.Nil,
	}
	role := models.RoleWithPermissions{
		Role: models.Role{
			Id:   uuid.Nil,
			Name: "admin",
		},
		Permissions: models.Permissions{
			permission,
		},
	}
	roleBody, _ := json.Marshal(role)
	roleInvalid := models.RoleWithPermissions{
		Role: models.Role{
			Name: "",
		},
		Permissions: models.Permissions{},
	}
	roleInvalidBody, _ := json.Marshal(roleInvalid)
	invalidRolePerms := models.RoleWithPermissions{
		Role: models.Role{
			Id:   uuid.Nil,
			Name: "admin",
		},
		Permissions: make([]models.Permission, models.LimitMaxPermissions+1),
	}
	invalidRolePermsBody, _ := json.Marshal(invalidRolePerms)

	// Define the test cases
	tests := []struct {
		name           string
		role           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to check permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleCreate).Return(false)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to decode",
			role: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleCreate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid role input",
			role: roleInvalidBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleCreate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid role permissions",
			role: invalidRolePermsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleCreate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Return(true)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at create",
			role: roleBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleCreate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.Nil, errors.New("error"))
				r.EXPECT().GetWithPermissions(gomock.Any()).Times(0)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve the role",
			role: roleBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleCreate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.New(), nil)
				r.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, false, errors.New("error"))
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "could not find the role",
			role: roleBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleCreate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.New(), nil)
				r.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, false, nil)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			role: roleBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleCreate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.New(), nil)
				r.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, true, nil)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", apiBasePath+"/roles", bytes.NewBuffer(tt.role))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.CreateRole(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetRoles tests the GetRoles handler
func TestGetRole(t *testing.T) {
	// Declare the data

	// Define the test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Without UUID param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
				)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not because of mock
		},
		{
			name: "Without permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false),
				)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetAllWithPermissions().Times(0)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not because of mock
		},
		{
			name: "Failed to get roles",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, false, errors.New("error"))
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Could not find the role",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, false, nil)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, true, nil)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/roles", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetRole(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetRoles tests the GetRoles handler
func TestGetRoles(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to check permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetAllWithPermissions().Times(0)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to retrieve the role",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetAllWithPermissions().Return(models.RolesWithPermissions{}, errors.New("error"))
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetAllWithPermissions().Return([]models.RoleWithPermissions{}, nil)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/roles", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetRoles(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestUpdateRole tests the UpdateRole handler
func TestUpdateRole(t *testing.T) {
	// Declare the data
	pRoleUpdate := "admin.roles.update"
	pRolePermUpdate := "admin.roles.permissions.update"
	permission := models.Permission{
		Id: uuid.Nil,
	}
	role := models.RoleWithPermissions{
		Role: models.Role{
			Id:   uuid.Nil,
			Name: "admin",
		},
		Permissions: models.Permissions{
			permission,
		},
	}
	roleBody, _ := json.Marshal(role)
	roleInvalid := models.RoleWithPermissions{
		Role: models.Role{
			Name: "",
		},
		Permissions: models.Permissions{},
	}
	roleInvalidBody, _ := json.Marshal(roleInvalid)
	invalidRolePerms := models.RoleWithPermissions{
		Role: models.Role{
			Id:   uuid.Nil,
			Name: "admin",
		},
		Permissions: make([]models.Permission, models.LimitMaxPermissions+1),
	}
	invalidRolePermsBody, _ := json.Marshal(invalidRolePerms)

	// Define the test cases
	tests := []struct {
		name           string
		role           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to parse the UUID",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleUpdate).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to check permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleUpdate).Return(false)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to decode",
			role: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleUpdate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid role input",
			role: roleInvalidBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleUpdate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid role permissions",
			role: invalidRolePermsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleUpdate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Return(true)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at update",
			role: roleBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleUpdate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				r.EXPECT().GetWithPermissions(gomock.Any()).Times(0)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve the role",
			role: roleBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleUpdate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				r.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, false, errors.New("error"))
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "could not find the role",
			role: roleBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleUpdate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				r.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, false, nil)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			role: roleBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRoleUpdate).Return(true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), pRolePermUpdate).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				r.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, true, nil)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", apiBasePath+"/roles", bytes.NewBuffer(tt.role))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.UpdateRole(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestDeleteRole tests the DeleteRole handler
func TestDeleteRole(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Without UUID param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
				)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not because of mock
		},
		{
			name: "Without permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false),
				)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().Delete(gomock.Any()).Times(0)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not because of mock
		},
		{
			name: "Failed to delete role",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().Delete(gomock.Any()).Return(nil)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/roles", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.DeleteRole(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetRolePermissions tests the GetRolePermissions handler
func TestGetRolePermissions(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Without UUID param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
				)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not because of mock
		},
		{
			name: "Without permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not because of mock
		},
		{
			name: "Failed to retrieve role permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Return(models.Permissions{}, errors.New("error"))
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Return(models.Permissions{}, nil)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/roles", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetRolePermissions(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestSetRolePermissions tests the SetRolePermissions handler
func TestSetRolePermissions(t *testing.T) {
	// Declare the data
	validPerms := models.Permissions{
		models.Permission{Id: uuid.New()},
	}
	validPermsBody, _ := json.Marshal(validPerms)
	invalidPerms := make([]models.Permission, models.LimitMaxPermissions+1)
	invalidPermsBody, _ := json.Marshal(invalidPerms)

	// Define the test cases
	tests := []struct {
		name           string
		body           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to check permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().SetRolePermissions(gomock.Any(), gomock.Any()).Times(0)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to decode",
			body: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().SetRolePermissions(gomock.Any(), gomock.Any()).Times(0)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at bad role input",
			body: invalidPermsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().SetRolePermissions(gomock.Any(), gomock.Any()).Times(0)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at set",
			body: validPermsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().SetRolePermissions(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validPermsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().SetRolePermissions(gomock.Any(), gomock.Any()).Return(nil)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", apiBasePath+"/roles/{id}/permissions", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.SetRolePermissions(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetRoleUsers tests the GetRoleUsers handler
func TestGetRoleUsers(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Without UUID param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
				)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not because of mock
		},
		{
			name: "Without permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false),
				)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetUsersByRoleID(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not because of mock
		},
		{
			name: "Failed to retrieve role permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetUsersByRoleID(gomock.Any()).Return([]models.User{}, errors.New("error"))
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetUsersByRoleID(gomock.Any()).Return([]models.User{}, nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/roles", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetRoleUsers(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestPutUsersRole tests the PutUsersRole handler
func TestPutUsersRole(t *testing.T) {
	// Declare the data
	validUUIDs := []uuid.UUID{uuid.New()}
	validUUIDsBody, _ := json.Marshal(validUUIDs)

	// Define the test cases
	tests := []struct {
		name           string
		body           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to check permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().AddUsersRole(gomock.Any(), gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to decode",
			body: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().AddUsersRole(gomock.Any(), gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at uuids length = 0",
			body: []byte(`[]`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().AddUsersRole(gomock.Any(), gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at add users role",
			body: validUUIDsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().AddUsersRole(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validUUIDsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().AddUsersRole(gomock.Any(), gomock.Any()).Return(nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", apiBasePath+"/roles/{id}/users", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.PutUsersRole(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestDeleteUsersRole tests the DeleteUsersRole handler
func TestDeleteUsersRole(t *testing.T) {
	// Declare the data
	validUUIDs := []uuid.UUID{uuid.New()}
	validUUIDsBody, _ := json.Marshal(validUUIDs)

	// Define the test cases
	tests := []struct {
		name           string
		body           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to check permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().RemoveUsersRole(gomock.Any(), gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to decode",
			body: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().RemoveUsersRole(gomock.Any(), gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at uuids length = 0",
			body: []byte(`[]`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().RemoveUsersRole(gomock.Any(), gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at RemoveUsersRole",
			body: validUUIDsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().RemoveUsersRole(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validUUIDsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().RemoveUsersRole(gomock.Any(), gomock.Any()).Return(nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", apiBasePath+"/roles/{id}/users", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.DeleteUsersRole(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
