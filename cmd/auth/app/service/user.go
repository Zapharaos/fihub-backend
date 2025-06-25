package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/auth/app/otp"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ResetForgottenPassword resets the forgotten password for a user
func (s *AuthService) ResetForgottenPassword(ctx context.Context, req *authpb.ResetForgottenPasswordRequest) (*authpb.ResetForgottenPasswordResponse, error) {
	purpose := authpb.OtpPurpose_PASSWORD_RESET

	// Validate the request
	err := otp.IsFinalRequestValid(ctx, purpose, req.GetUserId(), req.GetRequestId())
	if err != nil {
		return nil, err
	}

	// Setup metadata for gRPC clients as context
	// We can trust the userID here because it has been validated in the OTP process
	md := metadata.Pairs("x-user-id", req.GetUserId())
	userClientCtx := metadata.NewOutgoingContext(ctx, md)

	// Update the user password
	_, err = s.userClient.UpdateUserPassword(userClientCtx, &userpb.UpdateUserPasswordRequest{
		Id:           req.GetUserId(),
		Password:     req.GetPassword(),
		Confirmation: req.GetConfirmation(),
	})
	if err != nil {
		zap.L().Error("failed to update user password", zap.Error(err))
		return nil, err
	}

	// Delete the key from Redis
	otp.CleanupRedisKey(ctx, otp.BuildOtpRequestKey(req.GetUserId(), purpose))

	return &authpb.ResetForgottenPasswordResponse{
		Success: true,
	}, nil
}

// UpdatePassword updates the current user password
func (s *AuthService) UpdatePassword(ctx context.Context, req *authpb.UpdatePasswordRequest) (*authpb.UpdatePasswordResponse, error) {
	purpose := authpb.OtpPurpose_PASSWORD_CHANGE

	// Retrieve the userID from the context
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user ID not found in context")
	}

	// Validate the request
	err := otp.IsFinalRequestValid(ctx, authpb.OtpPurpose_PASSWORD_CHANGE, userID, req.GetRequestId())
	if err != nil {
		return nil, err
	}

	// Update the user password
	_, err = s.userClient.UpdateUserPassword(ctx, &userpb.UpdateUserPasswordRequest{
		Id:           userID,
		Password:     req.GetPassword(),
		Confirmation: req.GetConfirmation(),
	})
	if err != nil {
		zap.L().Error("failed to update user password", zap.Error(err))
		return nil, err
	}

	// Delete the key from Redis
	otp.CleanupRedisKey(ctx, otp.BuildOtpRequestKey(userID, purpose))

	return &authpb.UpdatePasswordResponse{
		Success: true,
	}, nil
}

// CreateUser creates a new user account
func (s *AuthService) CreateUser(ctx context.Context, req *authpb.CreateUserRequest) (*authpb.CreateUserResponse, error) {
	purpose := authpb.OtpPurpose_USER_SIGNUP

	// Validate the request
	err := otp.IsFinalRequestValid(ctx, purpose, req.GetEmail(), req.GetRequestId())
	if err != nil {
		return nil, err
	}

	// Update the user password
	_, err = s.userClient.CreateUser(ctx, &userpb.CreateUserRequest{
		Email:        req.GetEmail(),
		Password:     req.GetPassword(),
		Confirmation: req.GetConfirmation(),
		Checkbox:     req.GetCheckbox(),
	})
	if err != nil {
		zap.L().Error("failed to create user", zap.Error(err))
		return nil, err
	}

	// Delete the key from Redis
	otp.CleanupRedisKey(ctx, otp.BuildOtpRequestKey(req.GetEmail(), purpose))

	return &authpb.CreateUserResponse{
		Success: true,
	}, nil
}
