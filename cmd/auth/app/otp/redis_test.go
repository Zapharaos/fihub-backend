package otp

import (
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"testing"
)

func TestBuildOtpKey(t *testing.T) {
	identifier := "user123"
	purpose := authpb.OtpPurpose_UNSPECIFIED
	expected := "otp:user123:UNSPECIFIED"

	result := BuildOtpKey(identifier, purpose)
	if result != expected {
		t.Errorf("BuildOtpKey() = %v, want %v", result, expected)
	}
}

func TestBuildOtpRequestKey(t *testing.T) {
	identifier := "user456"
	purpose := authpb.OtpPurpose_UNSPECIFIED
	expected := "otp-request:user456:UNSPECIFIED"

	result := BuildOtpRequestKey(identifier, purpose)
	if result != expected {
		t.Errorf("BuildOtpRequestKey() = %v, want %v", result, expected)
	}
}
