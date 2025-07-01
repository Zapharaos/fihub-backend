package models

import (
	"github.com/google/uuid"
	"time"
)

// RequestUserOtp represents the request for a user otp
type RequestUserOtp struct {
	Email string `json:"email"`
}

// ResponseRequestUserOtp represents the response for a user otp request
type ResponseRequestUserOtp struct {
	Error      string    `json:"error,omitempty"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
	Identifier string    `json:"identifier"`
}

// ValidateUserOtp represents the request for validating a user otp
type ValidateUserOtp struct {
	UserID uuid.UUID `json:"user_id"`
	Otp    string    `json:"otp"`
}

// ResponseValidateUserOtp represents the response for a user otp validation
type ResponseValidateUserOtp struct {
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	RequestID string    `json:"request_id"`
}

type UserInputResetPassword struct {
	UserID       uuid.UUID `json:"user_id"`
	OtpRequestID uuid.UUID `json:"otp_request_id"`
	Password     string    `json:"password"`
	Confirmation string    `json:"confirmation"`
}

type UserInputChangePassword struct {
	OtpRequestID uuid.UUID `json:"otp_request_id"`
	Password     string    `json:"password"`
	Confirmation string    `json:"confirmation"`
}
