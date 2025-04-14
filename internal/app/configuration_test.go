package app

import (
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestInitConfiguration tests the InitConfiguration function.
func TestInitConfiguration(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Create a full test suite
		ts := test.TestSuite{}
		_ = ts.CreateFullTestSuite(t)
		defer ts.CleanTestSuite(t)

		// Assert that the configuration was loaded successfully
		assert.NoError(t, InitConfiguration("test"), "Configuration should load successfully")
		assert.Equal(t, "test", viper.GetString("APP_ENV"), "APP_ENV should be set to 'test'")
	})

	t.Run("Failure", func(t *testing.T) {
		err := InitConfiguration("nonexistent")

		// Assert that the configuration was not loaded
		assert.Error(t, err, "Configuration should fail to load")
	})
}
