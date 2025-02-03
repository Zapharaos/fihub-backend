package env

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// TestGetString tests the GetString function
// It verifies that the function returns the correct value for an existing key
// and the default value for a non-existent key.
func TestGetString(t *testing.T) {
	t.Run("Non-existent key", func(t *testing.T) {
		fallback := "default"

		// Check if GetString returns the default value
		if val := GetString("NON_EXISTENT", fallback); val != fallback {
			t.Errorf("Expected '%s', got '%s'", fallback, val)
		}
	})

	t.Run("Unexpected type", func(t *testing.T) {
		fallback := "default"
		value := "123"

		// Set an environment variable to an unexpected type
		os.Setenv("TEST_STRING", value)
		defer os.Unsetenv("TEST_STRING")

		// Check if GetString returns the correct value
		if val := GetString("TEST_STRING", fallback); val != value {
			t.Errorf("Expected '%s', got '%s'", value, val)
		}
	})

	t.Run("Existing key", func(t *testing.T) {
		fallback := "default"
		value := "value"

		// Set an environment variable
		os.Setenv("TEST_STRING", value)
		defer os.Unsetenv("TEST_STRING")

		// Check if GetString returns the correct value
		if val := GetString("TEST_STRING", fallback); val != value {
			t.Errorf("Expected '%s', got '%s'", value, val)
		}
	})
}

// TestGetInt tests the GetInt function
// It verifies that the function returns the correct integer value for an existing key
// and the default integer value for a non-existent key.
func TestGetInt(t *testing.T) {
	t.Run("Non-existent key", func(t *testing.T) {
		fallback := 123

		// Check if GetInt returns the default integer value
		if val := GetInt("NON_EXISTENT", fallback); val != fallback {
			t.Errorf("Expected %d, got %d", fallback, val)
		}
	})

	t.Run("Unexpected type", func(t *testing.T) {
		fallback := 123
		value := "not_an_int"

		// Set an environment variable to an unexpected type
		os.Setenv("TEST_INT", value)
		defer os.Unsetenv("TEST_INT")

		// Check if GetInt returns the default value
		if val := GetInt("TEST_INT", fallback); val != fallback {
			t.Errorf("Expected %d, got %d", fallback, val)
		}
	})

	t.Run("Existing key", func(t *testing.T) {
		fallback := 123
		value := 42

		// Set an environment variable
		os.Setenv("TEST_INT", fmt.Sprintf("%d", value))
		defer os.Unsetenv("TEST_INT")

		// Check if GetInt returns the correct integer value
		if val := GetInt("TEST_INT", fallback); val != value {
			t.Errorf("Expected %d, got %d", value, val)
		}
	})
}

// TestGetBool tests the GetBool function
// It verifies that the function returns the correct boolean value for an existing key
// and the default boolean value for a non-existent key.
func TestGetBool(t *testing.T) {
	t.Run("Non-existent key", func(t *testing.T) {
		fallback := true

		// Check if GetBool returns the default boolean value
		if val := GetBool("NON_EXISTENT", fallback); val != fallback {
			t.Errorf("Expected %t, got %t", fallback, val)
		}
	})

	t.Run("Unexpected type", func(t *testing.T) {
		fallback := true
		value := "not_a_bool"

		// Set an environment variable to an unexpected type
		os.Setenv("TEST_BOOL", value)
		defer os.Unsetenv("TEST_BOOL")

		// Check if GetBool returns the default value
		if val := GetBool("TEST_BOOL", fallback); val != fallback {
			t.Errorf("Expected %t, got %t", fallback, val)
		}
	})

	t.Run("Existing key", func(t *testing.T) {
		fallback := true
		value := false

		// Set an environment variable
		os.Setenv("TEST_BOOL", fmt.Sprintf("%t", value))
		defer os.Unsetenv("TEST_BOOL")

		// Check if GetBool returns the correct boolean value
		if val := GetBool("TEST_BOOL", fallback); val != value {
			t.Errorf("Expected %t, got %t", value, val)
		}
	})
}

// TestGetDuration tests the GetDuration function
// It verifies that the function returns the correct duration value for an existing key
// and the default duration value for a non-existent key.
func TestGetDuration(t *testing.T) {
	t.Run("Non-existent key", func(t *testing.T) {
		fallback := 5 * time.Second

		// Check if GetDuration returns the default duration value
		if val := GetDuration("NON_EXISTENT", fallback); val != fallback {
			t.Errorf("Expected %s, got %s", fallback, val)
		}
	})

	t.Run("Unexpected type", func(t *testing.T) {
		fallback := 5 * time.Second
		value := "not_a_duration"

		// Set an environment variable to an unexpected type
		os.Setenv("TEST_DURATION", value)
		defer os.Unsetenv("TEST_DURATION")

		// Check if GetDuration returns the default value
		if val := GetDuration("TEST_DURATION", fallback); val != fallback {
			t.Errorf("Expected %s, got %s", fallback, val)
		}
	})

	t.Run("Existing key", func(t *testing.T) {
		fallback := 5 * time.Second
		value := 10 * time.Second

		// Set an environment variable
		os.Setenv("TEST_DURATION", value.String())
		defer os.Unsetenv("TEST_DURATION")

		// Check if GetDuration returns the correct duration value
		if val := GetDuration("TEST_DURATION", fallback); val != value {
			t.Errorf("Expected %s, got %s", value, val)
		}
	})
}
