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

func TestService_AddUsersToRole(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.AddUsersToRoleRequest{
		RoleId:  uuid.New().String(),
		UserIds: []string{uuid.New().String()},
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.AddUsersToRoleRequest
		expected        *securitypb.AddUsersToRoleResponse
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
				rr.EXPECT().AddToUsers(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.AddUsersToRoleRequest{},
			expected:        &securitypb.AddUsersToRoleResponse{},
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
				rr.EXPECT().AddToUsers(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.AddUsersToRoleRequest{
				RoleId: "bad-uuid",
			},
			expected:        &securitypb.AddUsersToRoleResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse string to UUIDs from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().AddToUsers(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.AddUsersToRoleRequest{
				RoleId:  uuid.New().String(),
				UserIds: []string{"bad-uuid"},
			},
			expected:        &securitypb.AddUsersToRoleResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to add users",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().AddToUsers(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
				rr.EXPECT().ListUsersByRoleId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.AddUsersToRoleResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to list users",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().AddToUsers(gomock.Any(), gomock.Any()).Return(nil)
				rr.EXPECT().ListUsersByRoleId(gomock.Any()).Return(nil, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.AddUsersToRoleResponse{},
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
				rr.EXPECT().AddToUsers(gomock.Any(), gomock.Any()).Return(nil)
				rr.EXPECT().ListUsersByRoleId(gomock.Any()).Return([]string{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.AddUsersToRoleResponse{UserIds: []string{}},
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
			response, err := service.AddUsersToRole(context.Background(), tt.request)

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

func TestService_RemoveUsersFromRole(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.RemoveUsersFromRoleRequest{
		RoleId:  uuid.New().String(),
		UserIds: []string{uuid.New().String()},
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.RemoveUsersFromRoleRequest
		expected        *securitypb.RemoveUsersFromRoleResponse
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
				rr.EXPECT().RemoveFromUsers(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.RemoveUsersFromRoleRequest{},
			expected:        &securitypb.RemoveUsersFromRoleResponse{},
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
				rr.EXPECT().RemoveFromUsers(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.RemoveUsersFromRoleRequest{
				RoleId: "bad-uuid",
			},
			expected:        &securitypb.RemoveUsersFromRoleResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse string to UUIDs from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().RemoveFromUsers(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.RemoveUsersFromRoleRequest{
				RoleId:  uuid.New().String(),
				UserIds: []string{"bad-uuid"},
			},
			expected:        &securitypb.RemoveUsersFromRoleResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to add users",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().RemoveFromUsers(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
				rr.EXPECT().ListUsersByRoleId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.RemoveUsersFromRoleResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to list users",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().RemoveFromUsers(gomock.Any(), gomock.Any()).Return(nil)
				rr.EXPECT().ListUsersByRoleId(gomock.Any()).Return(nil, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.RemoveUsersFromRoleResponse{},
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
				rr.EXPECT().RemoveFromUsers(gomock.Any(), gomock.Any()).Return(nil)
				rr.EXPECT().ListUsersByRoleId(gomock.Any()).Return([]string{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.RemoveUsersFromRoleResponse{UserIds: []string{}},
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
			response, err := service.RemoveUsersFromRole(context.Background(), tt.request)

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

func TestService_ListUsersForRole(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.ListUsersForRoleRequest{
		RoleId: uuid.New().String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.ListUsersForRoleRequest
		expected        *securitypb.ListUsersForRoleResponse
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
				rr.EXPECT().ListUsersByRoleId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.ListUsersForRoleRequest{},
			expected:        &securitypb.ListUsersForRoleResponse{},
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
				rr.EXPECT().ListUsersByRoleId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.ListUsersForRoleRequest{
				RoleId: "bad-uuid",
			},
			expected:        &securitypb.ListUsersForRoleResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to list users",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().ListUsersByRoleId(gomock.Any()).Return(nil, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.ListUsersForRoleResponse{},
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
				rr.EXPECT().ListUsersByRoleId(gomock.Any()).Return([]string{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.ListUsersForRoleResponse{UserIds: []string{}},
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
			response, err := service.ListUsersForRole(context.Background(), tt.request)

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

func TestService_SetRolesForUser(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.SetRolesForUserRequest{
		UserId:  uuid.New().String(),
		RoleIds: []string{uuid.New().String()},
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.SetRolesForUserRequest
		expected        *securitypb.SetRolesForUserResponse
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
				rr.EXPECT().SetForUser(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.SetRolesForUserRequest{},
			expected:        &securitypb.SetRolesForUserResponse{},
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
				rr.EXPECT().SetForUser(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.SetRolesForUserRequest{
				UserId: "bad-uuid",
			},
			expected:        &securitypb.SetRolesForUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse string to UUIDs from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().SetForUser(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.SetRolesForUserRequest{
				UserId:  uuid.New().String(),
				RoleIds: []string{"bad-uuid"},
			},
			expected:        &securitypb.SetRolesForUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to set user roles",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().SetForUser(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
				rr.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.SetRolesForUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to list user roles with permissions",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().SetForUser(gomock.Any(), gomock.Any()).Return(nil)
				rr.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Return(nil, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.SetRolesForUserResponse{},
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
				rr.EXPECT().SetForUser(gomock.Any(), gomock.Any()).Return(nil)
				rr.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Return(models.RolesWithPermissions{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: validRequest,
			expected: &securitypb.SetRolesForUserResponse{
				Roles: []*securitypb.RoleWithPermissions{},
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
			response, err := service.SetRolesForUser(context.Background(), tt.request)

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

func TestService_ListRolesForUser(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.ListRolesForUserRequest{
		UserId: uuid.New().String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.ListRolesForUserRequest
		expected        *securitypb.ListRolesForUserResponse
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
				rr.EXPECT().ListByUserId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.ListRolesForUserRequest{},
			expected:        &securitypb.ListRolesForUserResponse{},
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
				rr.EXPECT().ListByUserId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.ListRolesForUserRequest{
				UserId: "bad-uuid",
			},
			expected:        &securitypb.ListRolesForUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to list by user ID",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().ListByUserId(gomock.Any()).Return(nil, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.ListRolesForUserResponse{},
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
				rr.EXPECT().ListByUserId(gomock.Any()).Return(models.Roles{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: validRequest,
			expected: &securitypb.ListRolesForUserResponse{
				Roles: []*securitypb.Role{},
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
			response, err := service.ListRolesForUser(context.Background(), tt.request)

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

func TestService_ListRolesWithPermissionsForUser(t *testing.T) {
	service := &Service{}
	validRequest := &securitypb.ListRolesWithPermissionsForUserRequest{
		UserId: uuid.New().String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.ListRolesWithPermissionsForUserRequest
		expected        *securitypb.ListRolesWithPermissionsForUserResponse
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
				rr.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.ListRolesWithPermissionsForUserRequest{},
			expected:        &securitypb.ListRolesWithPermissionsForUserResponse{},
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
				rr.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.ListRolesWithPermissionsForUserRequest{
				UserId: "bad-uuid",
			},
			expected:        &securitypb.ListRolesWithPermissionsForUserResponse{},
			expectedErrCode: codes.InvalidArgument,
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
				rr.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Return(nil, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         validRequest,
			expected:        &securitypb.ListRolesWithPermissionsForUserResponse{},
			expectedErrCode: codes.Internal,
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
				rr.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Return(models.RolesWithPermissions{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: validRequest,
			expected: &securitypb.ListRolesWithPermissionsForUserResponse{
				Roles: []*securitypb.RoleWithPermissions{},
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
			response, err := service.ListRolesWithPermissionsForUser(context.Background(), tt.request)

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

func TestService_ListUsersFull(t *testing.T) {
	service := &Service{}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *securitypb.ListUsersFullRequest
		expected        *securitypb.ListUsersFullResponse
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
				rr.EXPECT().ListUsers().Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.ListUsersFullRequest{},
			expected:        &securitypb.ListUsersFullResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails to list users",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().ListUsers().Return(nil, errors.New("some error"))
				rr.EXPECT().ListByUserId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.ListUsersFullRequest{},
			expected:        &securitypb.ListUsersFullResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to list user details",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the role repository
				rr := mocks.NewSecurityRoleRepository(ctrl)
				rr.EXPECT().ListUsers().Return([]string{uuid.New().String()}, nil)
				rr.EXPECT().ListByUserId(gomock.Any()).Return(models.Roles{}, errors.New("some error"))
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request:         &securitypb.ListUsersFullRequest{},
			expected:        &securitypb.ListUsersFullResponse{},
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
				rr.EXPECT().ListUsers().Return([]string{uuid.New().String()}, nil)
				rr.EXPECT().ListByUserId(gomock.Any()).Return(models.Roles{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(rr, nil))
			},
			request: &securitypb.ListUsersFullRequest{},
			expected: &securitypb.ListUsersFullResponse{
				Users: []*securitypb.UserWithRoles{},
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
			response, err := service.ListUsersFull(context.Background(), tt.request)

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
