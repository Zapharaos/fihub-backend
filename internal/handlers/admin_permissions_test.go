package handlers_test

import (
	"bytes"
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/handlers"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
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
				p.EXPECT().Get(gomock.Any()).Return(permissions.Permission{}, false, errors.New("error"))
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
				p.EXPECT().Get(gomock.Any()).Return(permissions.Permission{}, false, nil)
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
				p.EXPECT().Get(gomock.Any()).Return(permissions.Permission{}, true, nil)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/permissions", bytes.NewBuffer(tt.body))

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
				p.EXPECT().Get(gomock.Any()).Return(permissions.Permission{}, false, errors.New("error"))
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
				p.EXPECT().Get(gomock.Any()).Return(permissions.Permission{}, false, nil)
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
				p.EXPECT().Get(gomock.Any()).Return(permissions.Permission{}, true, nil)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/permissions/{id}", nil)

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
				p.EXPECT().GetAll().Return([]permissions.Permission{}, nil)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/permissions", nil)

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
				p.EXPECT().Get(gomock.Any()).Return(permissions.Permission{}, false, errors.New("error"))
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
				p.EXPECT().Get(gomock.Any()).Return(permissions.Permission{}, false, nil)
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
				p.EXPECT().Get(gomock.Any()).Return(permissions.Permission{}, true, nil)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/permissions", bytes.NewBuffer(tt.permission))

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
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/permissions", nil)

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
