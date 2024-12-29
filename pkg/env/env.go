package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

// Load attempts to load environment variables from a file named ".env" in the current working directory
func Load() error {
	return godotenv.Load()
}

// GetString retrieves an environment variable by key and returns its value as a string.
func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

// GetInt retrieves an environment variable by key and returns its value as an integer.
func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("Error converting environment variable", key, "to int:", err)
		return fallback
	}

	return valAsInt
}

// GetBool retrieves an environment variable by key and returns its value as a boolean.
func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		// fmt.Println("Error converting environment variable", key, "to bool:", err)
		return fallback
	}

	return boolVal
}

// GetDuration retrieves an environment variable by key and returns its value as a time.Duration.
func GetDuration(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	durationVal, err := time.ParseDuration(val)
	if err != nil {
		fmt.Println("Error converting environment variable", key, "to duration:", err)
		return fallback
	}

	return durationVal
}
