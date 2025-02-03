package env

import (
	"os"
	"testing"
	"time"
)

// TestGetString tests the GetString function
// It verifies that the function returns the correct value for an existing key
// and the default value for a non-existent key.
func TestGetString(t *testing.T) {
	t.Run("Existing key", func(t *testing.T) {
		// Set an environment variable
		os.Setenv("TEST_STRING", "value")
		defer os.Unsetenv("TEST_STRING")

		// Check if GetString returns the correct value
		if val := GetString("TEST_STRING", "default"); val != "value" {
			t.Errorf("Expected 'value', got '%s'", val)
		}
	})

	t.Run("Non-existent key", func(t *testing.T) {
		// Check if GetString returns the default value
		if val := GetString("NON_EXISTENT", "default"); val != "default" {
			t.Errorf("Expected 'default', got '%s'", val)
		}
	})
}

// TestGetInt tests the GetInt function
// It verifies that the function returns the correct integer value for an existing key
// and the default integer value for a non-existent key.
func TestGetInt(t *testing.T) {
	t.Run("Existing key", func(t *testing.T) {
		// Set an environment variable
		os.Setenv("TEST_INT", "42")
		defer os.Unsetenv("TEST_INT")

		// Check if GetInt returns the correct integer value
		if val := GetInt("TEST_INT", 0); val != 42 {
			t.Errorf("Expected 42, got %d", val)
		}
	})

	t.Run("Non-existent key", func(t *testing.T) {
		// Check if GetInt returns the default integer value
		if val := GetInt("NON_EXISTENT", 0); val != 0 {
			t.Errorf("Expected 0, got %d", val)
		}
	})
}

// TestGetBool tests the GetBool function
// It verifies that the function returns the correct boolean value for an existing key
// and the default boolean value for a non-existent key.
func TestGetBool(t *testing.T) {
	t.Run("Existing key", func(t *testing.T) {
		// Set an environment variable
		os.Setenv("TEST_BOOL", "true")
		defer os.Unsetenv("TEST_BOOL")

		// Check if GetBool returns the correct boolean value
		if val := GetBool("TEST_BOOL", false); val != true {
			t.Errorf("Expected true, got %v", val)
		}
	})

	t.Run("Non-existent key", func(t *testing.T) {
		// Check if GetBool returns the default boolean value
		if val := GetBool("NON_EXISTENT", false); val != false {
			t.Errorf("Expected false, got %v", val)
		}
	})
}

// TestGetDuration tests the GetDuration function
// It verifies that the function returns the correct duration value for an existing key
// and the default duration value for a non-existent key.
func TestGetDuration(t *testing.T) {
	t.Run("Existing key", func(t *testing.T) {
		// Set an environment variable
		os.Setenv("TEST_DURATION", "1h")
		defer os.Unsetenv("TEST_DURATION")

		// Check if GetDuration returns the correct duration value
		expected := time.Hour
		if val := GetDuration("TEST_DURATION", 0); val != expected {
			t.Errorf("Expected %v, got %v", expected, val)
		}
	})

	t.Run("Non-existent key", func(t *testing.T) {
		// Check if GetDuration returns the default duration value
		if val := GetDuration("NON_EXISTENT", 0); val != 0 {
			t.Errorf("Expected 0, got %v", val)
		}
	})
}
