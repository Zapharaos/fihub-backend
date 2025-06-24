package otp

import (
	"context"
	"crypto/sha256"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func GetOtpTimeLimit() time.Duration {
	// TODO : handle different duration depending on purpose?
	timeLimit := viper.GetDuration("OTP_DURATION")
	if timeLimit == 0 {
		timeLimit = 15 * time.Minute
	}
	return timeLimit
}

func GetFinalRequestTimeLimit() time.Duration {
	// TODO : add as config variable?
	// TODO : handle different duration depending on purpose?
	return 15 * time.Minute
}

func hash(value string) []byte {
	// Hash the OTP value using SHA-256
	hashed := sha256.Sum256([]byte(value))
	return hashed[:]
}

func compareInputWithHash(input string, hashed string) bool {
	// Hash the input
	hashedInput := hash(input)
	return string(hashedInput) == hashed
}

func Generate() (string, []byte) {
	// TODO : handle different length depending on purpose?
	otpValue := utils.RandDigitString(viper.GetInt("OTP_LENGTH"))
	hashed := hash(otpValue)
	return otpValue, hashed
}

func IsOtpValid(ctx context.Context, purpose authpb.OtpPurpose, userID, inputOtp string) error {
	// Retrieve otp
	otpKey := BuildOtpKey(userID, purpose)
	storedHashOtp, err := GetRedisKey(ctx, otpKey)
	if err != nil {
		return err
	}

	// Compare hashes
	if !compareInputWithHash(inputOtp, storedHashOtp) {
		return status.Error(codes.InvalidArgument, "invalid OTP")
	}

	return nil
}

func IsFinalRequestValid(ctx context.Context, purpose authpb.OtpPurpose, userID, inputRequestID string) error {
	// Validate the request
	requestKey := BuildOtpRequestKey(userID, purpose)
	requestID, err := GetRedisKey(ctx, requestKey)
	if err != nil {
		zap.L().Error("failed to get OTP request ID", zap.Error(err))
		return err
	}
	if requestID != inputRequestID {
		zap.L().Error("invalid OTP request ID", zap.String("request_id", inputRequestID))
		return status.Error(codes.InvalidArgument, "invalid OTP request ID")
	}

	return nil
}
