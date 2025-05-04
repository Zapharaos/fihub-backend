package service

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestService_CreatePermission(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.CreatePermissionRequest{
		Value:       "value",
		Scope:       "admin",
		Description: "description",
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.CreatePermissionRequest
		expected        *securitypb.CreatePermissionResponse
		expectedErrCode codes.Code
	}{
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         &securitypb.CreatePermissionRequest{},
			expected:        &securitypb.CreatePermissionResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "invalid permission input",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request: &securitypb.CreatePermissionRequest{
				Value: "",
			},
			expected:        &securitypb.CreatePermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to create",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Create(gomock.Any()).Return(uuid.Nil, errors.New("error"))
				pr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         validRequest,
			expected:        &securitypb.CreatePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to retrieve permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				pr.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         validRequest,
			expected:        &securitypb.CreatePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				pr.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         validRequest,
			expected:        &securitypb.CreatePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "success",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				pr.EXPECT().Get(gomock.Any()).Return(models.Permission{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request: validRequest,
			expected: &securitypb.CreatePermissionResponse{
				Permission: &securitypb.Permission{},
			},
			expectedErrCode: codes.OK,
		},
	}

	// Run tests
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

func TestService_GetPermission(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.GetPermissionRequest{
		Id: uuid.New().String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.GetPermissionRequest
		expected        *securitypb.GetPermissionResponse
		expectedErrCode codes.Code
	}{
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         &securitypb.GetPermissionRequest{},
			expected:        &securitypb.GetPermissionResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request: &securitypb.GetPermissionRequest{
				Id: "bad-uuid",
			},
			expected:        &securitypb.GetPermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to retrieve permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         validRequest,
			expected:        &securitypb.GetPermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         validRequest,
			expected:        &securitypb.GetPermissionResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "success",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Get(gomock.Any()).Return(models.Permission{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request: validRequest,
			expected: &securitypb.GetPermissionResponse{
				Permission: &securitypb.Permission{},
			},
			expectedErrCode: codes.OK,
		},
	}

	// Run tests
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

func TestService_UpdatePermission(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.UpdatePermissionRequest{
		Id:          uuid.New().String(),
		Value:       "value",
		Scope:       "admin",
		Description: "description",
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.UpdatePermissionRequest
		expected        *securitypb.UpdatePermissionResponse
		expectedErrCode codes.Code
	}{
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         &securitypb.UpdatePermissionRequest{},
			expected:        &securitypb.UpdatePermissionResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request: &securitypb.UpdatePermissionRequest{
				Id: "bad-uuid",
			},
			expected:        &securitypb.UpdatePermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "invalid permission input",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request: &securitypb.UpdatePermissionRequest{
				Id:    uuid.New().String(),
				Value: "",
			},
			expected:        &securitypb.UpdatePermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to update",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Update(gomock.Any()).Return(errors.New("error"))
				pr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         validRequest,
			expected:        &securitypb.UpdatePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to retrieve permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Update(gomock.Any()).Return(nil)
				pr.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         validRequest,
			expected:        &securitypb.UpdatePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Update(gomock.Any()).Return(nil)
				pr.EXPECT().Get(gomock.Any()).Return(models.Permission{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         validRequest,
			expected:        &securitypb.UpdatePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "success",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Update(gomock.Any()).Return(nil)
				pr.EXPECT().Get(gomock.Any()).Return(models.Permission{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request: validRequest,
			expected: &securitypb.UpdatePermissionResponse{
				Permission: &securitypb.Permission{},
			},
			expectedErrCode: codes.OK,
		},
	}

	// Run tests
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

func TestService_DeletePermission(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.DeletePermissionRequest{
		Id: uuid.New().String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.DeletePermissionRequest
		expected        *securitypb.DeletePermissionResponse
		expectedErrCode codes.Code
	}{
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         &securitypb.DeletePermissionRequest{},
			expected:        &securitypb.DeletePermissionResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request: &securitypb.DeletePermissionRequest{
				Id: "bad-uuid",
			},
			expected:        &securitypb.DeletePermissionResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to delete permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         validRequest,
			expected:        &securitypb.DeletePermissionResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "success",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().Delete(gomock.Any()).Return(nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         validRequest,
			expected:        &securitypb.DeletePermissionResponse{},
			expectedErrCode: codes.OK,
		},
	}

	// Run tests
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

func TestService_ListPermissions(t *testing.T) {
	service := &Service{}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.ListPermissionsRequest
		expected        *securitypb.ListPermissionsResponse
		expectedErrCode codes.Code
	}{
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().List().Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         &securitypb.ListPermissionsRequest{},
			expected:        &securitypb.ListPermissionsResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails to list permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().List().Return(models.Permissions{}, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request:         &securitypb.ListPermissionsRequest{},
			expected:        &securitypb.ListPermissionsResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "success",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				pr := mocks.NewSecurityPermissionRepository(ctrl)
				pr.EXPECT().List().Return(models.Permissions{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, pr))
			},
			request: &securitypb.ListPermissionsRequest{},
			expected: &securitypb.ListPermissionsResponse{
				Permissions: []*securitypb.Permission{},
			},
			expectedErrCode: codes.OK,
		},
	}

	// Run tests
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
