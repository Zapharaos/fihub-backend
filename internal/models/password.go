package models

import (
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"time"
)

// PasswordResponseRequest represents the response for a password reset request
type PasswordResponseRequest struct {
	Error     string    `json:"error,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	UserID    uuid.UUID `json:"user_id"`
}

type PasswordInputRequest struct {
	Email string `json:"email"`
}

type PasswordRequest struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func InitPasswordRequest(userID uuid.UUID) (PasswordRequest, time.Duration) {
	duration := viper.GetDuration("OTP_DURATION")
	if duration == 0 {
		duration = 15 * time.Minute
	}
	return PasswordRequest{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     utils.RandDigitString(viper.GetInt("OTP_LENGTH")),
		ExpiresAt: time.Now().Add(duration),
		CreatedAt: time.Now(),
	}, duration
}
