package password

import (
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"time"
)

// ResponseRequest represents the response for a password reset request
type ResponseRequest struct {
	Error     string    `json:"error,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	UserID    uuid.UUID `json:"user_id"`
}

type InputRequest struct {
	Email string `json:"email"`
}

type Request struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func InitRequest(userID uuid.UUID) (Request, time.Duration) {
	duration := viper.GetDuration("OTP_DURATION")
	if duration == 0 {
		duration = 15 * time.Minute
	}
	return Request{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     utils.RandDigitString(viper.GetInt("OTP_LENGTH")),
		ExpiresAt: time.Now().Add(duration),
		CreatedAt: time.Now(),
	}, duration
}
