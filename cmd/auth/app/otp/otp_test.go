package otp

import (
	"github.com/spf13/viper"
	"testing"
	"time"
)

func TestGetOtpTimeLimit(t *testing.T) {
	t.Run("default value", func(t *testing.T) {
		viper.Set("OTP_DURATION", 0)
		expected := 15 * time.Minute
		got := GetOtpTimeLimit()
		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})

	t.Run("config value", func(t *testing.T) {
		viper.Set("OTP_DURATION", 5*time.Minute)
		expected := 5 * time.Minute
		got := GetOtpTimeLimit()
		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
}

func TestGetFinalRequestTimeLimit(t *testing.T) {
	t.Run("default value", func(t *testing.T) {
		viper.Set("OTP_FINAL_REQUEST_DURATION", 0)
		expected := 15 * time.Minute
		got := GetFinalRequestTimeLimit()
		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})

	t.Run("config value", func(t *testing.T) {
		viper.Set("OTP_FINAL_REQUEST_DURATION", 10*time.Minute)
		expected := 10 * time.Minute
		got := GetFinalRequestTimeLimit()
		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
}
func TestHashAndCompareInputWithHash(t *testing.T) {
	input := "123456"
	hashed := hash(input)
	hashedStr := string(hashed)

	if !compareInputWithHash(input, hashedStr) {
		t.Errorf("compareInputWithHash should return true for correct input")
	}

	if compareInputWithHash("654321", hashedStr) {
		t.Errorf("compareInputWithHash should return false for incorrect input")
	}
}

func TestGenerate(t *testing.T) {
	viper.Set("OTP_LENGTH", 6)
	otp, hashed := Generate()

	if len(otp) != 6 {
		t.Errorf("expected OTP length 6, got %d", len(otp))
	}
	if !compareInputWithHash(otp, string(hashed)) {
		t.Errorf("generated OTP should match its hash")
	}
}
