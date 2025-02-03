package env

import (
	"os"
	"testing"
	"time"
)

func TestGetString(t *testing.T) {
	t.Run("Existing key", func(t *testing.T) {
		os.Setenv("TEST_STRING", "value")
		defer os.Unsetenv("TEST_STRING")

		if val := GetString("TEST_STRING", "default"); val != "value" {
			t.Errorf("Expected 'value', got '%s'", val)
		}
	})

	t.Run("Non-existent key", func(t *testing.T) {
		if val := GetString("NON_EXISTENT", "default"); val != "default" {
			t.Errorf("Expected 'default', got '%s'", val)
		}
	})
}

func TestGetInt(t *testing.T) {
	t.Run("Existing key", func(t *testing.T) {
		os.Setenv("TEST_INT", "42")
		defer os.Unsetenv("TEST_INT")

		if val := GetInt("TEST_INT", 0); val != 42 {
			t.Errorf("Expected 42, got %d", val)
		}
	})

	t.Run("Non-existent key", func(t *testing.T) {
		if val := GetInt("NON_EXISTENT", 0); val != 0 {
			t.Errorf("Expected 0, got %d", val)
		}
	})
}

func TestGetBool(t *testing.T) {
	t.Run("Existing key", func(t *testing.T) {
		os.Setenv("TEST_BOOL", "true")
		defer os.Unsetenv("TEST_BOOL")

		if val := GetBool("TEST_BOOL", false); val != true {
			t.Errorf("Expected true, got %v", val)
		}
	})

	t.Run("Non-existent key", func(t *testing.T) {
		if val := GetBool("NON_EXISTENT", false); val != false {
			t.Errorf("Expected false, got %v", val)
		}
	})
}

func TestGetDuration(t *testing.T) {
	t.Run("Existing key", func(t *testing.T) {
		os.Setenv("TEST_DURATION", "1h")
		defer os.Unsetenv("TEST_DURATION")

		expected := time.Hour
		if val := GetDuration("TEST_DURATION", 0); val != expected {
			t.Errorf("Expected %v, got %v", expected, val)
		}
	})

	t.Run("Non-existent key", func(t *testing.T) {
		if val := GetDuration("NON_EXISTENT", 0); val != 0 {
			t.Errorf("Expected 0, got %v", val)
		}
	})
}
