package service

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

// TestCreateUser tests the CreateUser service
func TestCreateUser(t *testing.T) {
	service := &Service{}
	validRequest := &userpb.CreateUserRequest{
		Email:        "email@example.com",
		Password:     "password",
		Confirmation: "password",
		Checkbox:     true,
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *userpb.CreateUserRequest
		expected        *userpb.CreateUserResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(ur)
			},
			request: nil,
			expected: &userpb.CreateUserResponse{
				User: nil,
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails at bad user input",
			mockSetup: func(ctrl *gomock.Controller) {
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(ur)
			},
			request: &userpb.CreateUserRequest{
				Email:        "emailbadexample.com",
				Password:     "",
				Confirmation: "password",
				Checkbox:     false,
			},
			expected: &userpb.CreateUserResponse{
				User: nil,
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "Fail to check existence",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Return(false, errors.New("error"))
				u.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(u)
			},
			request: validRequest,
			expected: &userpb.CreateUserResponse{
				User: nil,
			},
			expectedErrCode: codes.Internal,
		},
		{
			name: "User already exists",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Return(true, nil)
				u.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(u)
			},
			request: validRequest,
			expected: &userpb.CreateUserResponse{
				User: nil,
			},
			expectedErrCode: codes.AlreadyExists,
		},
		{
			name: "Fail at create",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Return(false, nil)
				u.EXPECT().Create(gomock.Any()).Return(uuid.Nil, errors.New("error"))
				u.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(u)
			},
			request: validRequest,
			expected: &userpb.CreateUserResponse{
				User: nil,
			},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Fails to retrieve the user",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Return(false, nil)
				u.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				u.EXPECT().Get(gomock.Any()).Return(models.User{}, false, errors.New("error"))
				repositories.ReplaceGlobals(u)
			},
			request: validRequest,
			expected: &userpb.CreateUserResponse{
				User: nil,
			},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Could not find the user",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Return(false, nil)
				u.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				u.EXPECT().Get(gomock.Any()).Return(models.User{}, false, nil)
				repositories.ReplaceGlobals(u)
			},
			request: validRequest,
			expected: &userpb.CreateUserResponse{
				User: nil,
			},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Return(false, nil)
				u.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				u.EXPECT().Get(gomock.Any()).Return(models.User{}, true, nil)
				repositories.ReplaceGlobals(u)
			},
			request: validRequest,
			expected: &userpb.CreateUserResponse{
				User: &userpb.User{},
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
			response, err := service.CreateUser(context.Background(), tt.request)

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

// TestGetUser tests the GetUser service
func TestGetUser(t *testing.T) {
	service := &Service{}
	validRequest := &userpb.GetUserRequest{
		Id: uuid.New().String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *userpb.GetUserRequest
		expected        *userpb.GetUserResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
			},
			request:         nil,
			expected:        &userpb.GetUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
			},
			request: &userpb.GetUserRequest{
				Id: "bad-uuid",
			},
			expected:        &userpb.GetUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(ur)
			},
			request:         validRequest,
			expected:        &userpb.GetUserResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "Fails to retrieve the user",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(models.User{}, false, errors.New("error"))
				repositories.ReplaceGlobals(u)
			},
			request:         validRequest,
			expected:        &userpb.GetUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Could not find the user",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(models.User{}, false, nil)
				repositories.ReplaceGlobals(u)
			},
			request:         validRequest,
			expected:        &userpb.GetUserResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "Succeeds",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(models.User{}, true, nil)
				repositories.ReplaceGlobals(u)
			},
			request: validRequest,
			expected: &userpb.GetUserResponse{
				User: &userpb.User{},
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
			response, err := service.GetUser(context.Background(), tt.request)

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

// TestUpdateUser tests the UpdateUser service
func TestUpdateUser(t *testing.T) {
	service := &Service{}
	validRequest := &userpb.UpdateUserRequest{
		Id:    uuid.New().String(),
		Email: "email@example.com",
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *userpb.UpdateUserRequest
		expected        *userpb.UpdateUserResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
			},
			request:         nil,
			expected:        &userpb.UpdateUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
			},
			request: &userpb.UpdateUserRequest{
				Id: "bad-uuid",
			},
			expected:        &userpb.UpdateUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(ur)
			},
			request:         validRequest,
			expected:        &userpb.UpdateUserResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "user input is not valid",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(ur)
			},
			request: &userpb.UpdateUserRequest{
				Id:    uuid.New().String(),
				Email: "bad-email",
			},
			expected:        &userpb.UpdateUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to update user",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().Update(gomock.Any()).Return(errors.New("error"))
				repositories.ReplaceGlobals(ur)
			},
			request:         validRequest,
			expected:        &userpb.UpdateUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Fails to retrieve the user",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Update(gomock.Any()).Return(nil)
				u.EXPECT().Get(gomock.Any()).Return(models.User{}, false, errors.New("error"))
				repositories.ReplaceGlobals(u)
			},
			request:         validRequest,
			expected:        &userpb.UpdateUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Could not find the user",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Update(gomock.Any()).Return(nil)
				u.EXPECT().Get(gomock.Any()).Return(models.User{}, false, nil)
				repositories.ReplaceGlobals(u)
			},
			request:         validRequest,
			expected:        &userpb.UpdateUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Succeeds",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Update(gomock.Any()).Return(nil)
				u.EXPECT().Get(gomock.Any()).Return(models.User{}, true, nil)
				repositories.ReplaceGlobals(u)
			},
			request: validRequest,
			expected: &userpb.UpdateUserResponse{
				User: &userpb.User{},
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
			response, err := service.UpdateUser(context.Background(), tt.request)

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

// TestUpdateUserPassword tests the UpdateUserPassword service
func TestUpdateUserPassword(t *testing.T) {
	service := &Service{}
	validRequest := &userpb.UpdateUserPasswordRequest{
		Id:           uuid.New().String(),
		Password:     "password",
		Confirmation: "password",
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *userpb.UpdateUserPasswordRequest
		expected        *userpb.UpdateUserPasswordResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
			},
			request:         nil,
			expected:        &userpb.UpdateUserPasswordResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
			},
			request: &userpb.UpdateUserPasswordRequest{
				Id: "bad-uuid",
			},
			expected:        &userpb.UpdateUserPasswordResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().UpdateWithPassword(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(ur)
			},
			request:         validRequest,
			expected:        &userpb.UpdateUserPasswordResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "invalid password input",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().UpdateWithPassword(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(ur)
			},
			request: &userpb.UpdateUserPasswordRequest{
				Id:           uuid.New().String(),
				Password:     "password",
				Confirmation: "bad-confirmation",
			},
			expected:        &userpb.UpdateUserPasswordResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to update user password",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().UpdateWithPassword(gomock.Any()).Return(errors.New("bad-error"))
				repositories.ReplaceGlobals(ur)
			},
			request:         validRequest,
			expected:        &userpb.UpdateUserPasswordResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeds",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().UpdateWithPassword(gomock.Any()).Return(nil)
				repositories.ReplaceGlobals(ur)
			},
			request: validRequest,
			expected: &userpb.UpdateUserPasswordResponse{
				Success: true,
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
			response, err := service.UpdateUserPassword(context.Background(), tt.request)

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

// TestDeleteUser tests the DeleteUser service
func TestDeleteUser(t *testing.T) {
	service := &Service{}
	validRequest := &userpb.DeleteUserRequest{
		Id: uuid.New().String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *userpb.DeleteUserRequest
		expected        *userpb.DeleteUserResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
			},
			request:         nil,
			expected:        &userpb.DeleteUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
			},
			request: &userpb.DeleteUserRequest{
				Id: "bad-uuid",
			},
			expected:        &userpb.DeleteUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(ur)
			},
			request:         validRequest,
			expected:        &userpb.DeleteUserResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails to delete user",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().Delete(gomock.Any()).Return(errors.New("some error"))
				repositories.ReplaceGlobals(ur)
			},
			request:         validRequest,
			expected:        &userpb.DeleteUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeds",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().Delete(gomock.Any()).Return(nil)
				repositories.ReplaceGlobals(ur)
			},
			request:         validRequest,
			expected:        &userpb.DeleteUserResponse{},
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
			response, err := service.DeleteUser(context.Background(), tt.request)

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

// TestListUsers tests the ListUsers service
func TestListUsers(t *testing.T) {
	service := &Service{}
	validRequest := &userpb.ListUsersRequest{}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		permissionValue bool
		request         *userpb.ListUsersRequest
		expected        *userpb.ListUsersResponse
		expectedErrCode codes.Code
	}{
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().List().Times(0)
				repositories.ReplaceGlobals(ur)
			},
			permissionValue: false,
			request:         validRequest,
			expected:        &userpb.ListUsersResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails to delete user",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().List().Return(nil, errors.New("some error"))
				repositories.ReplaceGlobals(ur)
			},
			request:         validRequest,
			expected:        &userpb.ListUsersResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeds",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the user repository
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().List().Return(models.Users{}, nil)
				repositories.ReplaceGlobals(ur)
			},
			request: validRequest,
			expected: &userpb.ListUsersResponse{
				Users: []*userpb.User{},
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
			response, err := service.ListUsers(context.Background(), tt.request)

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

// TestAuthenticateUser tests the AuthenticateUser service
func TestAuthenticateUser(t *testing.T) {
	service := &Service{}
	validRequest := &userpb.AuthenticateUserRequest{
		Email:    "email",
		Password: "password",
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *userpb.AuthenticateUserRequest
		expected        *userpb.AuthenticateUserResponse
		expectedErrCode codes.Code
	}{
		{
			name: "Fails to authenticate",
			mockSetup: func(ctrl *gomock.Controller) {
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(models.User{}, false, errors.New("error"))
				repositories.ReplaceGlobals(ur)
			},
			request: validRequest,
			expected: &userpb.AuthenticateUserResponse{
				User: nil,
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				ur := mocks.NewUserRepository(ctrl)
				ur.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(models.User{
					ID:        uuid.New(),
					Email:     "email",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, true, nil)
				repositories.ReplaceGlobals(ur)
			},
			request: validRequest,
			expected: &userpb.AuthenticateUserResponse{
				User: &userpb.User{},
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
			response, err := service.AuthenticateUser(context.Background(), tt.request)

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
