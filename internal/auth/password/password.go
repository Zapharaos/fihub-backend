package password

import (
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"github.com/google/uuid"
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
	duration := env.GetDuration("OTP_DURATION", 15*time.Minute)
	return Request{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     utils.RandDigitString(env.GetInt("OTP_LENGTH", 6)),
		ExpiresAt: time.Now().Add(duration),
		CreatedAt: time.Now(),
	}, duration
}
