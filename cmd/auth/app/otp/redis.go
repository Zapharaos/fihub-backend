package otp

import (
	"context"
	"errors"
	"fmt"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func GetRedisKey(ctx context.Context, key string) (string, error) {
	value, err := database.DB().Redis().Client.Get(ctx, key).Result()
	if err != nil {
		// Differentiate between not found and other errors
		if errors.Is(err, redis.Nil) {
			return "", status.Error(codes.NotFound, "OTP not found")
		}
		return "", status.Error(codes.Internal, "failed to get OTP")
	}

	// If the key exists, return its value
	return value, nil
}

func GetTtlForRedisKey(ctx context.Context, key string) (time.Duration, error) {
	// Retrieve the key from Redis
	exists, err := database.DB().Redis().Client.Exists(ctx, key).Result()
	if err != nil {
		return -1, status.Error(codes.Internal, "failed to check key existence")
	}

	// Check if the key exists
	if exists == 1 {
		// If key exists, get its expiration time
		ttl, err := database.DB().Redis().Client.TTL(ctx, key).Result()
		if err != nil {
			return -1, status.Error(codes.Internal, "failed to get expiration")
		}

		return ttl, nil
	}

	return -1, nil
}

func CleanupRedisKey(ctx context.Context, key string) {
	if err := database.DB().Redis().Client.Del(ctx, key).Err(); err != nil {
		zap.L().Warn("failed to cleanup after failure", zap.Error(err))
	}
}

func BuildOtpKey(userID string, purpose authpb.OtpPurpose) string {
	return fmt.Sprintf("otp:%s:%s", userID, purpose)
}

func BuildOtpRequestKey(userID string, purpose authpb.OtpPurpose) string {
	return fmt.Sprintf("otp-request:%s:%s", userID, purpose)
}
