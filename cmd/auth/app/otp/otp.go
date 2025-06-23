package otp

import (
	"crypto/sha256"
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

func hash(value string) []byte {
	// Hash the OTP value using SHA-256
	hashed := sha256.Sum256([]byte(value))
	return hashed[:]
}

func Generate() (string, []byte) {
	// TODO : handle different length depending on purpose?
	otpValue := utils.RandDigitString(viper.GetInt("OTP_LENGTH"))
	hashed := hash(otpValue)
	return otpValue, hashed
}

func CompareInputWithHash(input string, hashed string) bool {
	// Hash the input
	hashedInput := hash(input)
	return string(hashedInput) == hashed
}
