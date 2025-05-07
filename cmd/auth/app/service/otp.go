package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/auth/app/otp"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (s *AuthService) handlePasswordOTP(ctx context.Context, req *authpb.GenerateOTPRequest) (*authpb.GenerateOTPResponse, error) {
	// Verify if the user exists
	response, err := s.userClient.GetByEmail(ctx, &userpb.GetByEmailRequest{
		Email: req.GetEmail(),
	})
	if err != nil {
		zap.L().Error("failed to check user existence", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to check user existence")
	}
	if response.GetUser() == nil || response.GetUser().GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user not found")
	}

	// Check for existing OTP with userID and purpose
	userID := response.GetUser().GetId()
	otpKey := otp.BuildOtpKey(userID, req.GetPurpose())
	ttl, err := otp.GetTtlForRedisKey(ctx, otpKey)
	if err != nil {
		zap.L().Error("failed to get OTP expiration", zap.Error(err))
		return nil, err
	}
	if ttl > 0 {
		return &authpb.GenerateOTPResponse{
			ExpiresAt: timestamppb.New(time.Now().Add(ttl)),
		}, nil
	}

	// Prepare OTP data
	otpTimeLimit := otp.GetTimeLimit()
	otpValue, otpHash := otp.GenerateOTPValueAndHash()

	// Store OTP in Redis with email and purpose as key
	err = database.DB().Redis().Client.Set(ctx, otpKey, otpHash, otpTimeLimit).Err()
	if err != nil {
		zap.L().Error("failed to store OTP", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to store OTP")
	}

	// Prepare otp email content
	userLanguage := language.MustParse(req.GetLanguage())
	subject, plainTextContent, htmlContent, err := otp.BuildOtpEmailContents(userLanguage, otpValue, otpTimeLimit)
	if err != nil {
		// Delete the request since the email could not be sent
		otp.CleanupRedisKey(ctx, otpKey)

		zap.L().Error("failed to build OTP email content", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to build OTP email content")
	}

	// Send email
	err = email.S().Send(req.GetEmail(), subject, plainTextContent, htmlContent)
	if err != nil {
		// Delete the request since the email could not be sent
		otp.CleanupRedisKey(ctx, otpKey)

		zap.L().Error("Failed to send OTP email", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to send OTP email")
	}

	return &authpb.GenerateOTPResponse{
		UserId:    userID,
		ExpiresAt: timestamppb.New(time.Now().Add(otpTimeLimit)),
	}, nil
}

// GenerateOTP generates a one-time password (OTP) for the user
func (s *AuthService) GenerateOTP(ctx context.Context, req *authpb.GenerateOTPRequest) (*authpb.GenerateOTPResponse, error) {
	// TODO : move handlers middleware rate limiter to here? attempts count?

	switch req.Purpose {
	case authpb.OtpPurpose_PASSWORD_CHANGE, authpb.OtpPurpose_PASSWORD_RESET:
		return s.handlePasswordOTP(ctx, req)
	case authpb.OtpPurpose_EMAIL_VERIFICATION:
		return nil, status.Error(codes.Unimplemented, "email verification not yet implemented")
		//return s.handleResetOTP(ctx, req)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid OTP purpose")
	}
}

// ValidateOTP validates the one-time password (OTP) for the user
func (s *AuthService) ValidateOTP(ctx context.Context, req *authpb.ValidateOTPRequest) (*authpb.ValidateOTPResponse, error) {
	// Retrieve otp
	otpKey := otp.BuildOtpKey(req.GetUserId(), req.GetPurpose())
	storedHashOtp, err := otp.GetRedisKey(ctx, otpKey)
	if err != nil {
		return nil, err
	}

	// Compare hashes
	if otp.CompareInputWithHash(req.GetOtp(), storedHashOtp) {
		return nil, status.Error(codes.InvalidArgument, "invalid OTP")
	}

	// Prepare next step data
	requestID := uuid.New().String()
	requestTimeLimit := 15 * time.Minute // TODO : handle different duration depending on purpose?
	requestKey := otp.BuildOtpRequestKey(req.GetUserId(), req.GetPurpose())

	// Prepare pipeline to store request ID and delete OTP
	pipe := database.DB().Redis().Client.TxPipeline()
	pipe.Del(ctx, otpKey)
	pipe.SetEx(ctx, requestKey, requestID, requestTimeLimit)

	if req.GetPurpose() == authpb.OtpPurpose_EMAIL_VERIFICATION {
		// TODO : handle user account activation
	}

	// Execute pipeline
	if _, err = pipe.Exec(ctx); err != nil {
		zap.L().Error("failed to store request ID", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to store request ID")
	}

	return &authpb.ValidateOTPResponse{
		RequestId: requestID,
	}, nil
}
