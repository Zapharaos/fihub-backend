package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
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
	validResponse := &securitypb.CreateRoleResponse{
		Role: &securitypb.RoleWithPermissions{
			Role: &securitypb.Role{
				Id: uuid.New().String(),
			},
		},
	}

	// Define the test cases
	tests := []struct {
		name           string
		role           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to decode",
			role: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at create",
			role: roleBody,
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "success",
			role: roleBody,
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
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
			r := httptest.NewRequest("POST", apiBasePath+"/security/role", bytes.NewBuffer(tt.role))

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
	validResponse := &securitypb.GetRoleResponse{
		Role: &securitypb.RoleWithPermissions{
			Role: &securitypb.Role{
				Id: uuid.New().String(),
			},
		},
	}

	// Define the test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Without UUID param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().GetRole(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not because of mock
		},
		{
			name: "Failed to get roles",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().GetRole(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Success",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().GetRole(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
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
			r := httptest.NewRequest("GET", apiBasePath+"/security/role/"+uuid.NewString(), nil)

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

// TestListRoles tests the ListRoles handler
func TestListRoles(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to list the roles",
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListRoles(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "success",
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListRoles(gomock.Any(), gomock.Any()).Return(&securitypb.ListRolesResponse{
					Roles: []*securitypb.RoleWithPermissions{},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
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
			r := httptest.NewRequest("GET", apiBasePath+"/security/role", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.ListRoles(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestUpdateRole tests the UpdateRole handler
func TestUpdateRole(t *testing.T) {
	// Declare the data
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
	validResponse := &securitypb.UpdateRoleResponse{
		Role: &securitypb.RoleWithPermissions{
			Role: &securitypb.Role{
				Id: uuid.New().String(),
			},
		},
	}

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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().UpdateRole(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to decode",
			role: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().UpdateRole(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at update",
			role: roleBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().UpdateRole(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			role: roleBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().UpdateRole(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().DeleteRole(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not because of mock
		},
		{
			name: "Failed to delete role",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().DeleteRole(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().DeleteRole(gomock.Any(), gomock.Any()).Return(&securitypb.DeleteRoleResponse{}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListRolePermissions(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not because of mock
		},
		{
			name: "Failed to retrieve role permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListRolePermissions(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListRolePermissions(gomock.Any(), gomock.Any()).Return(&securitypb.ListRolePermissionsResponse{
					Permissions: []*securitypb.Permission{},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().SetRolePermissions(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to decode",
			body: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().SetRolePermissions(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failed to delete role",
			body: validPermsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().SetRolePermissions(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			body: validPermsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().SetRolePermissions(gomock.Any(), gomock.Any()).Return(&securitypb.SetRolePermissionsResponse{
					Permissions: []*securitypb.Permission{},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", apiBasePath+"/security/role/"+uuid.New().String()+"/permission", bytes.NewBuffer(tt.body))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt.mockSetup(ctrl)

			handlers.SetRolePermissions(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestAddUsersToRole tests the AddUsersToRole handler
func TestAddUsersToRole(t *testing.T) {
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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().AddUsersToRole(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to decode",
			body: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().AddUsersToRole(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at add users role",
			body: validUUIDsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().AddUsersToRole(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validUUIDsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().AddUsersToRole(gomock.Any(), gomock.Any()).Return(&securitypb.AddUsersToRoleResponse{
					UserIds: []string{},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
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
			r := httptest.NewRequest("PUT", apiBasePath+"/security/role/"+uuid.New().String()+"/user", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.AddUsersToRole(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestRemoveUsersFromRole tests the RemoveUsersFromRole handler
func TestRemoveUsersFromRole(t *testing.T) {
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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().RemoveUsersFromRole(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to decode",
			body: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().RemoveUsersFromRole(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to remove users from role",
			body: validUUIDsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().RemoveUsersFromRole(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validUUIDsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().RemoveUsersFromRole(gomock.Any(), gomock.Any()).Return(&securitypb.RemoveUsersFromRoleResponse{
					UserIds: []string{},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
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
			r := httptest.NewRequest("DELETE", apiBasePath+"/security/role/"+uuid.New().String()+"/user", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.RemoveUsersFromRole(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestListUsersForRole tests the ListUsersForRole handler
func TestListUsersForRole(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Without UUID param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListUsersForRole(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not because of mock
		},
		{
			name: "Failed to retrieve role permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListUsersForRole(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListUsersForRole(gomock.Any(), gomock.Any()).Return(&securitypb.ListUsersForRoleResponse{
					UserIds: []string{},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
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
			r := httptest.NewRequest("GET", apiBasePath+"/security/role/"+uuid.New().String()+"/user", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.ListUsersForRole(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestSetRolesForUser tests the SetRolesForUser handler
func TestSetRolesForUser(t *testing.T) {
	// Prepare data
	validUser := []string{uuid.New().String(), uuid.New().String()}
	validUserBody, _ := json.Marshal(validUser)

	// Test cases
	tests := []struct {
		name           string
		user           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to parse param",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().SetRolesForUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "Fails to decode",
			user: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().SetRolesForUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail at set roles",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().SetRolesForUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().SetRolesForUser(gomock.Any(), gomock.Any()).Return(&securitypb.SetRolesForUserResponse{}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", apiBasePath+"/security/role/user/"+uuid.New().String(), bytes.NewBuffer(tt.user))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.SetRolesForUser(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestListRolesWithPermissionsForUser tests the ListRolesWithPermissionsForUser handler
func TestListRolesWithPermissionsForUser(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListRolesWithPermissionsForUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "Fail at get roles",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListRolesWithPermissionsForUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListRolesWithPermissionsForUser(gomock.Any(), gomock.Any()).Return(&securitypb.ListRolesWithPermissionsForUserResponse{
					Roles: []*securitypb.RoleWithPermissions{},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/security/role/user/"+uuid.New().String(), nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.ListRolesWithPermissionsForUser(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestListUsersWithRoles tests the ListUsersWithRoles handler
func TestListUsersWithRoles(t *testing.T) {
	userID := uuid.New().String()
	validSecurityUsers := []*securitypb.UserWithRoles{
		{
			UserId: userID,
			Roles: []*securitypb.Role{
				{
					Id: uuid.New().String(),
				},
			},
		},
	}
	validUserUsers := []*userpb.User{
		{
			Id: userID,
		},
	}

	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to list users from roles",
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListUsersFull(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Fails to list users from users",
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListUsersFull(gomock.Any(), gomock.Any()).Return(&securitypb.ListUsersFullResponse{
					Users: validSecurityUsers,
				}, nil)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListUsersFull(gomock.Any(), gomock.Any()).Return(&securitypb.ListUsersFullResponse{
					Users: validSecurityUsers,
				}, nil)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(&userpb.ListUsersResponse{
					Users: validUserUsers,
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/security/role/user", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.ListUsersWithRoles(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
