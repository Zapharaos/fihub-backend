package otp

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/spf13/viper"
	"time"
)

func GetTimeLimit() time.Duration {
	// TODO : handle different duration depending on purpose?
	timeLimit := viper.GetDuration("OTP_DURATION")
	if timeLimit == 0 {
		timeLimit = 15 * time.Minute
	}
	return timeLimit
}

func GenerateOTPValueAndHash() (string, [32]byte) {
	// TODO : handle different length depending on purpose?
	otpValue := utils.RandDigitString(viper.GetInt("OTP_LENGTH"))
	hashed := sha256.Sum256([]byte(otpValue))
	return otpValue, hashed
}

func CompareInputWithHash(input string, hash string) bool {
	// Hash the input
	hashedInput := sha256.Sum256([]byte(input))
	// Compare the hashes
	return hash != hex.EncodeToString(hashedInput[:])
}
