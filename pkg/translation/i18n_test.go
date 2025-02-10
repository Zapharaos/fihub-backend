package translation

import (
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

var defaultLang = language.English

// TestNewI18nService tests the creation of a new I18nService instance.
func TestNewI18nService(t *testing.T) {
	t.Run("Without translation files", func(t *testing.T) {
		// Expect a panic when translation files are missing
		assert.Panics(t, func() {
			NewI18nService(defaultLang)
		})
	})

	t.Run("With translation files", func(t *testing.T) {
		// Setup test suite with translation files
		ts := test.TestSuite{}
		_ = ts.CreateConfigTranslationsFullTestSuite(t)
		defer ts.CleanTestSuite(t)

		// Expect no panic and a non-nil service instance
		service := NewI18nService(defaultLang)
		assert.NotNil(t, service)
	})
}

// TestI18nService_Localizer tests the retrieval of localizers from I18nService.
func TestI18nService_Localizer(t *testing.T) {
	// Setup test suite with translation files
	ts := test.TestSuite{}
	_ = ts.CreateConfigTranslationsFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Create a new I18nService instance
	service := NewI18nService(defaultLang)

	t.Run("Retrieve default language localizer", func(t *testing.T) {
		// Expect no error and a non-nil localizer for the default language
		localizer, err := service.Localizer(defaultLang)
		assert.NoError(t, err)
		assert.NotNil(t, localizer)
	})

	t.Run("Retrieve non-existing language localizer", func(t *testing.T) {
		// Expect an error when trying to retrieve a localizer for a non-existing language
		_, err := service.Localizer(language.Spanish)
		assert.Error(t, err)
	})
}

// TestI18nService_Message tests the retrieval of localized messages from I18nService.
func TestI18nService_Message(t *testing.T) {
	// Setup test suite with translation files
	ts := test.TestSuite{}
	_ = ts.CreateConfigTranslationsFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Create a new I18nService instance
	service := NewI18nService(defaultLang)
	localizer, _ := service.Localizer(defaultLang)

	t.Run("Retrieve existing message", func(t *testing.T) {
		// Define a message with ID "hello" and template data
		message := &Message{
			ID: "hello",
			Data: map[string]interface{}{
				"name": "World",
			},
		}
		// Assuming "hello" message is defined in active.en.toml
		result := service.Message(localizer, message)
		assert.Equal(t, "Hello, World!", result)
	})

	t.Run("Retrieve non-existing message", func(t *testing.T) {
		// Define a message with a non-existing ID
		message := &Message{
			ID: "nonexistent",
		}
		// Expect the result to be an empty string
		result := service.Message(localizer, message)
		assert.Equal(t, "", result)
	})
}
