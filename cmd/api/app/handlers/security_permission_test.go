package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreatePermission tests the CreatePermission function
func TestCreatePermission(t *testing.T) {
	// Declare the data
	validPermission := models.Permission{
		Value:       "value",
		Scope:       models.AdminScope,
		Description: "description",
	}
	validPermissionBody, _ := json.Marshal(validPermission)
	validResponse := &securitypb.CreatePermissionResponse{
		Permission: &securitypb.Permission{
			Id:          uuid.New().String(),
			Value:       validPermission.Value,
			Scope:       validPermission.Scope,
			Description: validPermission.Description,
		},
	}

	// Define the test cases
	tests := []struct {
		name           string
		body           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "With decode error",
			body: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().CreatePermission(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to create permission",
			body: validPermissionBody,
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().CreatePermission(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "With success",
			body: validPermissionBody,
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().CreatePermission(gomock.Any(), gomock.Any()).Return(validResponse, nil)
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
			r := httptest.NewRequest("POST", apiBasePath+"/security/permission", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call the function
			handlers.CreatePermission(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetPermission tests the GetPermission function
func TestGetPermission(t *testing.T) {
	validPermission := models.Permission{
		Id:          uuid.New(),
		Value:       "value",
		Scope:       models.AdminScope,
		Description: "description",
	}
	validResponse := &securitypb.GetPermissionResponse{
		Permission: &securitypb.Permission{
			Id:          validPermission.Id.String(),
			Value:       validPermission.Value,
			Scope:       validPermission.Scope,
			Description: validPermission.Description,
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
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not because of mock
		},
		{
			name: "With get error",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().GetPermission(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "With success",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().GetPermission(gomock.Any(), gomock.Any()).Return(validResponse, nil)
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
			r := httptest.NewRequest("GET", apiBasePath+"/security/permission/"+uuid.New().String(), nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call the function
			handlers.GetPermission(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestUpdatePermission tests the UpdatePermission function
func TestUpdatePermission(t *testing.T) {
	// Declare the data
	validPermission := models.Permission{
		Value:       "value",
		Scope:       models.AdminScope,
		Description: "description",
	}
	validPermissionBody, _ := json.Marshal(validPermission)
	validResponse := &securitypb.UpdatePermissionResponse{
		Permission: &securitypb.Permission{
			Id:          uuid.New().String(),
			Value:       validPermission.Value,
			Scope:       validPermission.Scope,
			Description: validPermission.Description,
		},
	}

	// Define the test cases
	tests := []struct {
		name           string
		permission     []byte
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
				sc.EXPECT().UpdatePermission(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not because of mock
		},
		{
			name:       "With decode error",
			permission: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().CreatePermission(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "fails to update permission",
			permission: validPermissionBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().UpdatePermission(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:       "With retrieve not found",
			permission: validPermissionBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().UpdatePermission(gomock.Any(), gomock.Any()).Return(validResponse, nil)
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
			r := httptest.NewRequest("PUT", apiBasePath+"/security/permission/"+uuid.New().String(), bytes.NewBuffer(tt.permission))

			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.UpdatePermission(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestDeletePermission tests the DeletePermission function
func TestDeletePermission(t *testing.T) {
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
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not because of mock
		},
		{
			name: "With delete error",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().DeletePermission(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "With success",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().DeletePermission(gomock.Any(), gomock.Any()).Return(&securitypb.DeletePermissionResponse{
					Success: true,
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
			r := httptest.NewRequest("DELETE", apiBasePath+"/security/permission/"+uuid.New().String(), nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call the function
			handlers.DeletePermission(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestListPermissions tests the ListPermissions function
func TestListPermissions(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "With get all error",
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListPermissions(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithSecurityClient(sc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "With success",
			mockSetup: func(ctrl *gomock.Controller) {
				sc := mocks.NewMockSecurityServiceClient(ctrl)
				sc.EXPECT().ListPermissions(gomock.Any(), gomock.Any()).Return(&securitypb.ListPermissionsResponse{
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
			r := httptest.NewRequest("GET", apiBasePath+"/security/permission", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call the function
			handlers.ListPermissions(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
