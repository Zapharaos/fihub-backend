package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/auth/app/otp"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ResetForgottenPassword resets the forgotten password for a user
func (s *AuthService) ResetForgottenPassword(ctx context.Context, req *authpb.ResetForgottenPasswordRequest) (*authpb.ResetForgottenPasswordResponse, error) {
	// Validate the request
	requestKey := otp.BuildOtpRequestKey(req.GetUserId(), authpb.OtpPurpose_PASSWORD_RESET)
	requestID, err := otp.GetRedisKey(ctx, requestKey)
	if err != nil {
		zap.L().Error("failed to get OTP request ID", zap.Error(err))
		return nil, err
	}
	if requestID != req.GetRequestId() {
		zap.L().Error("invalid OTP request ID", zap.String("request_id", req.GetRequestId()))
		return nil, status.Error(codes.InvalidArgument, "invalid OTP request ID")
	}

	// Update the user password
	_, err = s.userClient.UpdateUserPassword(ctx, &userpb.UpdateUserPasswordRequest{
		Id:           req.GetUserId(),
		Password:     req.GetPassword(),
		Confirmation: req.GetConfirmation(),
	})
	if err != nil {
		zap.L().Error("failed to update user password", zap.Error(err))
		return nil, err
	}

	// Delete the key from Redis
	otp.CleanupRedisKey(ctx, requestKey)

	return &authpb.ResetForgottenPasswordResponse{
		Success: true,
	}, nil
}

// UpdatePassword updates the current user password
func (s *AuthService) UpdatePassword(ctx context.Context, req *authpb.UpdatePasswordRequest) (*authpb.UpdatePasswordResponse, error) {
	// Validate the request
	requestKey := otp.BuildOtpRequestKey(req.GetUserId(), authpb.OtpPurpose_PASSWORD_CHANGE)
	requestID, err := otp.GetRedisKey(ctx, requestKey)
	if err != nil {
		zap.L().Error("failed to get OTP request ID", zap.Error(err))
		return nil, err
	}
	if requestID != req.GetRequestId() {
		zap.L().Error("invalid OTP request ID", zap.String("request_id", req.GetRequestId()))
		return nil, status.Error(codes.InvalidArgument, "invalid OTP request ID")
	}

	// Update the user password
	_, err = s.userClient.UpdateUserPassword(ctx, &userpb.UpdateUserPasswordRequest{
		Id:           req.GetUserId(),
		Password:     req.GetPassword(),
		Confirmation: req.GetConfirmation(),
	})
	if err != nil {
		zap.L().Error("failed to update user password", zap.Error(err))
		return nil, err
	}

	// Delete the key from Redis
	otp.CleanupRedisKey(ctx, requestKey)

	return &authpb.UpdatePasswordResponse{
		Success: true,
	}, nil
}
