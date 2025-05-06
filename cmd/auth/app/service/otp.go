package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Zapharaos/fihub-backend/cmd/auth/app/otp"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// GenerateOTP generates a one-time password (OTP) for the user
func (s *AuthService) GenerateOTP(ctx context.Context, req *authpb.GenerateOTPRequest) (*authpb.GenerateOTPResponse, error) {

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

	// TODO : move handlers middleware rate limiter to here? attempts count?

	// Check for existing OTP with userID and purpose
	userID := response.GetUser().GetId()
	otpKey := fmt.Sprintf("otp:%s:%s", userID, req.GetPurpose())
	exists, err := database.DB().Redis().Client.Exists(ctx, otpKey).Result()
	if err != nil {
		zap.L().Error("failed to check existing OTP", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to check existing OTP")
	}

	if exists == 1 {
		// If OTP exists, get its expiration time
		ttl, err := database.DB().Redis().Client.TTL(ctx, otpKey).Result()
		if err != nil {
			zap.L().Error("failed to get OTP expiration", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to get OTP expiration")
		}

		return &authpb.GenerateOTPResponse{
			ExpiresAt: timestamppb.New(time.Now().Add(ttl)),
		}, nil
	}

	// TODO : handle different duration and length depending on purpose?

	// Generate a new OTP
	timeLimit := viper.GetDuration("OTP_DURATION")
	if timeLimit == 0 {
		timeLimit = 15 * time.Minute
	}
	otpValue := utils.RandDigitString(viper.GetInt("OTP_LENGTH"))
	hashed := sha256.Sum256([]byte(otpValue))

	// Store OTP in Redis with email and purpose as key
	err = database.DB().Redis().Client.Set(ctx, otpKey, hashed, timeLimit).Err()
	if err != nil {
		zap.L().Error("failed to store OTP", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to store OTP")
	}

	// Prepare otp email content
	userLanguage := language.MustParse(req.GetLanguage())
	subject, plainTextContent, htmlContent, err := otp.BuildOtpEmailContents(userLanguage, otpValue, timeLimit)
	if err != nil {
		// Delete the request since the email could not be sent
		// TODO : delete OTP from Redis

		zap.L().Error("failed to build OTP email content", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to build OTP email content")
	}

	// Send email
	err = email.S().Send(req.GetEmail(), subject, plainTextContent, htmlContent)
	if err != nil {
		// Delete the request since the email could not be sent
		// TODO : delete OTP from Redis

		zap.L().Error("Failed to send OTP email", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to send OTP email")
	}

	return &authpb.GenerateOTPResponse{
		UserId:    userID,
		ExpiresAt: timestamppb.New(time.Now().Add(timeLimit)),
	}, nil
}

// ValidateOTP validates the one-time password (OTP) for the user
func (s *AuthService) ValidateOTP(ctx context.Context, req *authpb.ValidateOTPRequest) (*authpb.ValidateOTPResponse, error) {
	// Retrieve otp
	otpKey := fmt.Sprintf("otp:%s:%s", req.GetUserId(), req.GetPurpose())
	storedHashOtp, err := database.DB().Redis().Client.Get(ctx, otpKey).Result()
	if err != nil {
		// Differentiate between not found and other errors
		if errors.Is(err, redis.Nil) {
			return nil, status.Error(codes.NotFound, "OTP not found")
		}
		zap.L().Error("failed to get OTP", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get OTP")
	}

	// Compare hashes
	inputHashOtp := sha256.Sum256([]byte(req.GetOtp()))
	if storedHashOtp != hex.EncodeToString(inputHashOtp[:]) {
		return nil, status.Error(codes.InvalidArgument, "invalid OTP")
	}

	requestID := uuid.New().String()
	requestKey := fmt.Sprintf("otp_req:%s:%s", req.GetUserId(), req.GetPurpose())

	pipe := database.DB().Redis().Client.TxPipeline()
	pipe.Del(ctx, otpKey)
	// TODO : handle different duration depending on purpose?
	pipe.SetEx(ctx, requestKey, requestID, 15*time.Minute)
	if _, err = pipe.Exec(ctx); err != nil {
		zap.L().Error("failed to store request ID", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to store request ID")
	}

	return &authpb.ValidateOTPResponse{
		RequestId: requestID,
	}, nil
}
