package service

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

// TestCreatePermission tests the CreatePermission function
func TestCreatePermission(t *testing.T) {
	// Prepare data
	service := &Service{}
	validRequest := &protogen.CreatePermissionRequest{
		Value:       "test_permission",
		Scope:       models.AdminScope,
		Description: "test_description",
	}

	// Define the test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.CreatePermissionRequest
		expected        *protogen.CreatePermissionResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         nil,
			expected:        &protogen.CreatePermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails at bad permission input",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request: &protogen.CreatePermissionRequest{
				Value: "",
				Scope: "",
			},
			expected:        &protogen.CreatePermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "With create error",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Return(uuid.Nil, errors.New("error"))
				p.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.CreatePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "With retrieve error",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.CreatePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "With retrieve not found",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.CreatePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "With success",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.CreatePermissionResponse{},
			expectedErrCode: codes.OK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.CreatePermission(context.Background(), tt.request)

			// Handle errors
			if err != nil && tt.expectedErrCode == codes.OK {
				assert.Fail(t, "unexpected error", err)
			} else if err != nil {
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, tt.expectedErrCode, s.Code())
				} else {
					assert.Fail(t, "failed to get status from error")
				}
			}

			// Handle response
			if tt.expectedErrCode == codes.OK {
				assert.NotNil(t, response)
			} else {
				assert.Equal(t, tt.expected, response)
			}
		})
	}
}

// TestGetPermission tests the GetPermission function
func TestGetPermission(t *testing.T) {
	// Prepare data
	service := &Service{}
	validRequest := &protogen.GetPermissionRequest{
		Id: uuid.New().String(),
	}

	// Define the test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.GetPermissionRequest
		expected        *protogen.GetPermissionResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         nil,
			expected:        &protogen.GetPermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request: &protogen.GetPermissionRequest{
				Id: "bad-uuid",
			},
			expected:        &protogen.GetPermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "With get error",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.GetPermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "With retrieve not found",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.GetPermissionResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "With success",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.GetPermissionResponse{},
			expectedErrCode: codes.OK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.GetPermission(context.Background(), tt.request)

			// Handle errors
			if err != nil && tt.expectedErrCode == codes.OK {
				assert.Fail(t, "unexpected error", err)
			} else if err != nil {
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, tt.expectedErrCode, s.Code())
				} else {
					assert.Fail(t, "failed to get status from error")
				}
			}

			// Handle response
			if tt.expectedErrCode == codes.OK {
				assert.NotNil(t, response)
			} else {
				assert.Equal(t, tt.expected, response)
			}
		})
	}
}

// TestListPermissions tests the ListPermissions function
func TestListPermissions(t *testing.T) {
	// Prepare data
	service := &Service{}

	// Define the test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.ListPermissionsRequest
		expected        *protogen.ListPermissionsResponse
		expectedErrCode codes.Code
	}{
		{
			name: "With get all error",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().GetAll().Return(nil, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         nil,
			expected:        &protogen.ListPermissionsResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "With success",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().GetAll().Return(models.Permissions{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         nil,
			expected:        &protogen.ListPermissionsResponse{},
			expectedErrCode: codes.OK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.ListPermissions(context.Background(), tt.request)

			// Handle errors
			if err != nil && tt.expectedErrCode == codes.OK {
				assert.Fail(t, "unexpected error", err)
			} else if err != nil {
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, tt.expectedErrCode, s.Code())
				} else {
					assert.Fail(t, "failed to get status from error")
				}
			}

			// Handle response
			if tt.expectedErrCode == codes.OK {
				assert.NotNil(t, response)
			} else {
				assert.Equal(t, tt.expected, response)
			}
		})
	}
}

// TestUpdatePermission tests the UpdatePermission function
func TestUpdatePermission(t *testing.T) {
	// Prepare data
	service := &Service{}
	id := uuid.New()
	validRequest := &protogen.UpdatePermissionRequest{
		Id:          id.String(),
		Value:       "test_permission",
		Scope:       models.AdminScope,
		Description: "test_description",
	}

	// Define the test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.UpdatePermissionRequest
		expected        *protogen.UpdatePermissionResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         nil,
			expected:        &protogen.UpdatePermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request: &protogen.UpdatePermissionRequest{
				Id: "bad-uuid",
			},
			expected:        &protogen.UpdatePermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails at bad permission input",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         &protogen.UpdatePermissionRequest{},
			expected:        &protogen.UpdatePermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "With update error",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Return(errors.New("error"))
				p.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.UpdatePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "With retrieve error",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Return(nil)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.UpdatePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "With retrieve not found",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Return(nil)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.UpdatePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Succeed",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Update(gomock.Any()).Return(nil)
				p.EXPECT().Get(gomock.Any()).Return(models.Permission{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.UpdatePermissionResponse{},
			expectedErrCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.UpdatePermission(context.Background(), tt.request)

			// Handle errors
			if err != nil && tt.expectedErrCode == codes.OK {
				assert.Fail(t, "unexpected error", err)
			} else if err != nil {
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, tt.expectedErrCode, s.Code())
				} else {
					assert.Fail(t, "failed to get status from error")
				}
			}

			// Handle response
			if tt.expectedErrCode == codes.OK {
				assert.NotNil(t, response)
			} else {
				assert.Equal(t, tt.expected, response)
			}
		})
	}
}

// TestDeletePermission tests the DeletePermission function
func TestDeletePermission(t *testing.T) {
	// Prepare data
	service := &Service{}
	validRequest := &protogen.DeletePermissionRequest{
		Id: uuid.New().String(),
	}

	// Define the test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.DeletePermissionRequest
		expected        *protogen.DeletePermissionResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         nil,
			expected:        &protogen.DeletePermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request: &protogen.DeletePermissionRequest{
				Id: "bad-uuid",
			},
			expected:        &protogen.DeletePermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "With delete error",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.DeletePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "With success",
			mockSetup: func(ctrl *gomock.Controller) {
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().Delete(gomock.Any()).Return(nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, p))
			},
			request:         validRequest,
			expected:        &protogen.DeletePermissionResponse{},
			expectedErrCode: codes.OK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.DeletePermission(context.Background(), tt.request)

			// Handle errors
			if err != nil && tt.expectedErrCode == codes.OK {
				assert.Fail(t, "unexpected error", err)
			} else if err != nil {
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, tt.expectedErrCode, s.Code())
				} else {
					assert.Fail(t, "failed to get status from error")
				}
			}

			// Handle response
			if tt.expectedErrCode == codes.OK {
				assert.NotNil(t, response)
			} else {
				assert.Equal(t, tt.expected, response)
			}
		})
	}
}
