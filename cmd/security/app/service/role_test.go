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

func TestService_CreateRole(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.CreateRoleRequest{
		Name:        "name",
		Permissions: []string{uuid.New().String()},
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.CreateRoleRequest
		expected        *securitypb.CreateRoleResponse
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.CreateRoleRequest{},
			expected:        &securitypb.CreateRoleResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "invalid role input",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.CreateRoleRequest{
				Name: "",
			},
			expected:        &securitypb.CreateRoleResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to create",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil).Times(2)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.Nil, errors.New("some error"))
				rr.EXPECT().GetWithPermissions(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.CreateRoleResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to get role with permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil).Times(2)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.New(), nil)
				rr.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, false, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.CreateRoleResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find role with permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil).Times(2)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.New(), nil)
				rr.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.CreateRoleResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "success",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil).Times(2)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.New(), nil)
				rr.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: validRequest,
			expected: &securitypb.CreateRoleResponse{
				Role: &securitypb.RoleWithPermissions{},
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
			response, err := service.CreateRole(context.Background(), tt.request)

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

func TestService_GetRole(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.GetRoleRequest{
		Id: uuid.New().String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.GetRoleRequest
		expected        *securitypb.GetRoleResponse
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().GetWithPermissions(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.GetRoleRequest{},
			expected:        &securitypb.GetRoleResponse{},
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().GetWithPermissions(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.GetRoleRequest{
				Id: "bad-uuid",
			},
			expected:        &securitypb.GetRoleResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to get role with permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, false, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.GetRoleResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to get role with permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.GetRoleResponse{},
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.GetRoleResponse{},
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
			response, err := service.GetRole(context.Background(), tt.request)

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

func TestService_UpdateRole(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.UpdateRoleRequest{
		Id:          uuid.New().String(),
		Name:        "name",
		Permissions: []string{uuid.New().String()},
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.UpdateRoleRequest
		expected        *securitypb.UpdateRoleResponse
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.UpdateRoleRequest{},
			expected:        &securitypb.UpdateRoleResponse{},
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.UpdateRoleRequest{
				Id: "bad-uuid",
			},
			expected:        &securitypb.UpdateRoleResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "invalid role input",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.UpdateRoleRequest{
				Id:   uuid.New().String(),
				Name: "",
			},
			expected:        &securitypb.UpdateRoleResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to update",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil).Times(2)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
				rr.EXPECT().GetWithPermissions(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.UpdateRoleResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to get role with permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil).Times(2)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				rr.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, false, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.UpdateRoleResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find role with permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil).Times(2)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				rr.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.UpdateRoleResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "success",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil).Times(2)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				rr.EXPECT().GetWithPermissions(gomock.Any()).Return(models.RoleWithPermissions{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: validRequest,
			expected: &securitypb.UpdateRoleResponse{
				Role: &securitypb.RoleWithPermissions{},
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
			response, err := service.UpdateRole(context.Background(), tt.request)

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

func TestService_DeleteRole(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.DeleteRoleRequest{
		Id: uuid.New().String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.DeleteRoleRequest
		expected        *securitypb.DeleteRoleResponse
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.DeleteRoleRequest{},
			expected:        &securitypb.DeleteRoleResponse{},
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.DeleteRoleRequest{
				Id: "bad-uuid",
			},
			expected:        &securitypb.DeleteRoleResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to delete role",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Delete(gomock.Any()).Return(errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.DeleteRoleResponse{},
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().Delete(gomock.Any()).Return(nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.DeleteRoleResponse{},
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
			response, err := service.DeleteRole(context.Background(), tt.request)

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

func TestService_ListRoles(t *testing.T) {
	service := &Service{}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.ListRolesRequest
		expected        *securitypb.ListRolesResponse
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().ListWithPermissions().Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.ListRolesRequest{},
			expected:        &securitypb.ListRolesResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails to list roles",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().ListWithPermissions().Return(models.RolesWithPermissions{}, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.ListRolesRequest{},
			expected:        &securitypb.ListRolesResponse{},
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().ListWithPermissions().Return(models.RolesWithPermissions{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.ListRolesRequest{},
			expected: &securitypb.ListRolesResponse{
				Roles: []*securitypb.RoleWithPermissions{},
			},
			expectedErrCode: codes.Internal,
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
			response, err := service.ListRoles(context.Background(), tt.request)

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

func TestService_ListRolePermissions(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.ListRolePermissionsRequest{
		Id: uuid.New().String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.ListRolePermissionsRequest
		expected        *securitypb.ListRolePermissionsResponse
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().ListPermissionsByRoleId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.ListRolePermissionsRequest{},
			expected:        &securitypb.ListRolePermissionsResponse{},
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().ListPermissionsByRoleId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.ListRolePermissionsRequest{
				Id: "bad-uuid",
			},
			expected:        &securitypb.ListRolePermissionsResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to list permissions by role ID",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().ListPermissionsByRoleId(gomock.Any()).Return(models.Permissions{}, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.ListRolePermissionsResponse{},
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().ListPermissionsByRoleId(gomock.Any()).Return(models.Permissions{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: validRequest,
			expected: &securitypb.ListRolePermissionsResponse{
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
			response, err := service.ListRolePermissions(context.Background(), tt.request)

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

func TestService_SetRolePermissions(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.SetRolePermissionsRequest{
		Id:          uuid.New().String(),
		Permissions: []string{uuid.New().String()},
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.SetRolePermissionsRequest
		expected        *securitypb.SetRolePermissionsResponse
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().SetPermissionsByRoleId(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.SetRolePermissionsRequest{},
			expected:        &securitypb.SetRolePermissionsResponse{},
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().SetPermissionsByRoleId(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.SetRolePermissionsRequest{
				Id: "bad-uuid",
			},
			expected:        &securitypb.SetRolePermissionsResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse permissions ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().SetPermissionsByRoleId(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.SetRolePermissionsRequest{
				Id:          uuid.New().String(),
				Permissions: []string{"bad-uuid"},
			},
			expected:        &securitypb.SetRolePermissionsResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to set permissions by role ID",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().SetPermissionsByRoleId(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
				rr.EXPECT().ListPermissionsByRoleId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.SetRolePermissionsResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to list permissions by role ID",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().SetPermissionsByRoleId(gomock.Any(), gomock.Any()).Return(nil)
				rr.EXPECT().ListPermissionsByRoleId(gomock.Any()).Return(models.Permissions{}, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.SetRolePermissionsResponse{},
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
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().SetPermissionsByRoleId(gomock.Any(), gomock.Any()).Return(nil)
				rr.EXPECT().ListPermissionsByRoleId(gomock.Any()).Return(models.Permissions{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: validRequest,
			expected: &securitypb.SetRolePermissionsResponse{
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
			response, err := service.SetRolePermissions(context.Background(), tt.request)

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
