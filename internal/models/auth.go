package models

import (
	"github.com/google/uuid"
	"time"
)

// RequestUserOtp represents the request for a user otp
type RequestUserOtp struct {
	Email string `json:"email"`
}

// ResponseUserOtp represents the response for a user otp request
type ResponseUserOtp struct {
	Error     string    `json:"error,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	UserID    uuid.UUID `json:"user_id"`
}

// ValidateUserOtp represents the request for validating a user otp
type ValidateUserOtp struct {
	UserID uuid.UUID `json:"user_id"`
	Otp    string    `json:"otp"`
}
