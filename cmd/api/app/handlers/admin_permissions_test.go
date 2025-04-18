package handlers_test

import (
	"bytes"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/users/permissions"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreatePermission tests the CreatePermission function
func TestCreatePermission(t *testing.T) {
	// Declare the data
	permission := []byte(`{"value":"value","scope":"all","description":"description"}`)

	// Define the test cases
	tests := []struct {
		name           string
		body           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Without permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not because of mock
		},
		{
			name: "With decode error",
			body: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "With invalid permission data",
			body: []byte(`{"value":"","scope":"all","description":"description"}`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "With create error",
			body: permission,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Return(uuid.Nil, errors.New("error"))
				p.EXPECT().Get(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "With retrieve error",
			body: permission,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, errors.New("error"))
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "With retrieve not found",
			body: permission,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, nil)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "With success",
			body: permission,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, true, nil)
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
			r := httptest.NewRequest("POST", apiBasePath+"/permissions", bytes.NewBuffer(tt.body))

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
				p.EXPECT().Get(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not because of mock
		},
		{
			name: "With get error",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, errors.New("error"))
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "With retrieve not found",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, nil)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "With success",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, true, nil)
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
			r := httptest.NewRequest("GET", apiBasePath+"/permissions/{id}", nil)

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

// TestGetPermissions tests the GetPermissions function
func TestGetPermissions(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Without permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAll().Times(0)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not because of mock
		},
		{
			name: "With get all error",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAll().Return(nil, errors.New("error"))
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "With success",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAll().Return([]models.Permission{}, nil)
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
			r := httptest.NewRequest("GET", apiBasePath+"/permissions", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call the function
			handlers.GetPermissions(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestUpdatePermission tests the UpdatePermission function
func TestUpdatePermission(t *testing.T) {
	// Declare the data
	permission := []byte(`{"value":"value","scope":"all","description":"description"}`)

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
				p.EXPECT().Update(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not because of mock
		},
		{
			name:       "Wrong input",
			permission: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid input",
			permission: []byte(`{"value":"","scope":"all","description":"description"}`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "With update error",
			permission: permission,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Return(errors.New("error"))
				p.EXPECT().Get(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:       "With retrieve error",
			permission: permission,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Return(nil)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, errors.New("error"))
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:       "With retrieve not found",
			permission: permission,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Return(nil)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, nil)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:       "With retrieve not found",
			permission: permission,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Return(nil)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, true, nil)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", apiBasePath+"/permissions", bytes.NewBuffer(tt.permission))

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
				p.EXPECT().Delete(gomock.Any()).Times(0)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not because of mock
		},
		{
			name: "With delete error",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "With success",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true),
				)
				handlers.ReplaceGlobals(m)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().Delete(gomock.Any()).Return(nil)
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
			r := httptest.NewRequest("DELETE", apiBasePath+"/permissions", nil)

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
