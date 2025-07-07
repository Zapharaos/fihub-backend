package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

// TestGenerateToken tests the AuthService.GenerateToken service
func TestFindUsersIdentifiers_PasswordReset(t *testing.T) {
	validEmail := "email@test.com"
	validIdentifier := "test-identifier"

	// Define tests
	tests := []struct {
		name             string
		serviceSetup     func(ctrl *gomock.Controller) *AuthService
		request          *authpb.GenerateOTPRequest
		expectEmail      string
		expectIdentifier string
		expectedErrCode  codes.Code
	}{
		{
			name: "invalid purpose argument",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				return NewAuthService(userClient)
			},
			request: &authpb.GenerateOTPRequest{
				Purpose: authpb.OtpPurpose_UNSPECIFIED, // Invalid purpose
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "PasswordReset - error while searching for user by email",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Internal, "internal error"))
				return NewAuthService(userClient)
			},
			request: &authpb.GenerateOTPRequest{
				Purpose:    authpb.OtpPurpose_PASSWORD_RESET,
				Identifier: validEmail,
			},
			expectedErrCode: codes.Internal,
		},
		{
			name: "PasswordReset - fails to find user by email",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(&userpb.GetByEmailResponse{}, nil)
				return NewAuthService(userClient)
			},
			request: &authpb.GenerateOTPRequest{
				Purpose:    authpb.OtpPurpose_PASSWORD_RESET,
				Identifier: validEmail,
			},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "PasswordReset - successfully finds user by email",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(&userpb.GetByEmailResponse{
					User: &userpb.User{
						Id:    validIdentifier,
						Email: validEmail,
					},
				}, nil)
				return NewAuthService(userClient)
			},
			request: &authpb.GenerateOTPRequest{
				Purpose:    authpb.OtpPurpose_PASSWORD_RESET,
				Identifier: validEmail,
			},
			expectEmail:      validEmail,
			expectIdentifier: validIdentifier,
			expectedErrCode:  codes.OK,
		},
		{
			name: "PasswordChange - error while searching for user",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Internal, "internal error"))
				return NewAuthService(userClient)
			},
			request: &authpb.GenerateOTPRequest{
				Purpose:    authpb.OtpPurpose_PASSWORD_CHANGE,
				Identifier: validIdentifier,
			},
			expectedErrCode: codes.Internal,
		},
		{
			name: "PasswordChange - fails to find user",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&userpb.GetUserResponse{}, nil)
				return NewAuthService(userClient)
			},
			request: &authpb.GenerateOTPRequest{
				Purpose:    authpb.OtpPurpose_PASSWORD_CHANGE,
				Identifier: validIdentifier,
			},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "PasswordChange - successfully finds user",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&userpb.GetUserResponse{
					User: &userpb.User{
						Id:    validIdentifier,
						Email: validEmail,
					},
				}, nil)
				return NewAuthService(userClient)
			},
			request: &authpb.GenerateOTPRequest{
				Purpose:    authpb.OtpPurpose_PASSWORD_CHANGE,
				Identifier: validIdentifier,
			},
			expectEmail:      validEmail,
			expectIdentifier: validIdentifier,
			expectedErrCode:  codes.OK,
		},
		{
			name: "UserSignup - silenced error while searching for user by email",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.NotFound, "not found error"))
				return NewAuthService(userClient)
			},
			request: &authpb.GenerateOTPRequest{
				Purpose:    authpb.OtpPurpose_USER_SIGNUP,
				Identifier: validEmail,
			},
			expectEmail:      validEmail,
			expectIdentifier: validEmail,
			expectedErrCode:  codes.OK,
		},
		{
			name: "UserSignup - error while searching for user by email",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Internal, "internal error"))
				return NewAuthService(userClient)
			},
			request: &authpb.GenerateOTPRequest{
				Purpose:    authpb.OtpPurpose_USER_SIGNUP,
				Identifier: validEmail,
			},
			expectedErrCode: codes.Internal,
		},
		{
			name: "UserSignup - user already exists",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(&userpb.GetByEmailResponse{
					User: &userpb.User{
						Id: validIdentifier,
					},
				}, nil)
				return NewAuthService(userClient)
			},
			request: &authpb.GenerateOTPRequest{
				Purpose:    authpb.OtpPurpose_USER_SIGNUP,
				Identifier: validEmail,
			},
			expectedErrCode: codes.AlreadyExists,
		},
		{
			name: "UserSignup - successfully verifies user does not exist",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(&userpb.GetByEmailResponse{}, nil)
				return NewAuthService(userClient)
			},
			request: &authpb.GenerateOTPRequest{
				Purpose:    authpb.OtpPurpose_USER_SIGNUP,
				Identifier: validEmail,
			},
			expectEmail:      validEmail,
			expectIdentifier: validEmail,
			expectedErrCode:  codes.OK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Call service
			service := tt.serviceSetup(ctrl)
			responseEmail, responseIdentifier, err := service.findUserIdentifiers(context.Background(), tt.request)

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
				assert.Equal(t, tt.expectEmail, responseEmail)
				assert.Equal(t, tt.expectIdentifier, responseIdentifier)
			} else {
				assert.Empty(t, responseEmail)
				assert.Empty(t, responseIdentifier)
			}
		})
	}
}
