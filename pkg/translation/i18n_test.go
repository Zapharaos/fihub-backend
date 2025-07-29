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
		// Expect an error when translation files are missing
		_, err := NewI18nService(defaultLang, "nonexistent/path", "active")
		assert.Error(t, err)
	})

	t.Run("With translation files", func(t *testing.T) {
		// Setup test suite with translation files
		ts := test.TestSuite{}
		_ = ts.CreateConfigTranslationsFullTestSuite(t)
		defer ts.CleanTestSuite(t)

		// Expect no error and a non-nil service instance
		service, err := NewI18nService(defaultLang, "config/translations", "active")
		assert.NoError(t, err)
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
	service, err := NewI18nService(defaultLang, "config/translations", "active")
	assert.NoError(t, err)

	t.Run("Retrieve default language localizer", func(t *testing.T) {
		// Expect no error and a non-nil localizer for the default language
		localizer, found, err := service.Localizer(defaultLang)
		assert.NoError(t, err)
		assert.True(t, found)
		assert.NotNil(t, localizer)
	})

	t.Run("Retrieve non-existing language localizer", func(t *testing.T) {
		// Expect no error but found=false when trying to retrieve a localizer for a non-existing language
		// It should return the default language localizer instead
		localizer, found, err := service.Localizer(language.Spanish)
		assert.NoError(t, err)
		assert.False(t, found)      // Spanish not found, so it returns default
		assert.NotNil(t, localizer) // Should still return the default localizer
	})
}

// TestI18nService_Translate tests the translation of messages from I18nService.
func TestI18nService_Translate(t *testing.T) {
	// Setup test suite with translation files
	ts := test.TestSuite{}
	_ = ts.CreateConfigTranslationsFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Create a new I18nService instance
	service, err := NewI18nService(defaultLang, "config/translations", "active")
	assert.NoError(t, err)

	localizer, found, err := service.Localizer(defaultLang)
	assert.NoError(t, err)
	assert.True(t, found)

	t.Run("Translate existing message", func(t *testing.T) {
		// Define a message with ID "hello" and template data
		message := &Message{
			ID: "hello",
			Data: map[string]interface{}{
				"name": "World",
			},
		}
		// Assuming "hello" message is defined in active.en.toml
		result, success, err := service.Translate(localizer, message)
		assert.NoError(t, err)
		assert.True(t, success)
		assert.Equal(t, "Hello, World!", result)
	})

	t.Run("Translate non-existing message", func(t *testing.T) {
		// Define a message with a non-existing ID
		message := &Message{
			ID: "nonexistent",
		}
		// Expect an error when trying to translate a non-existing message
		result, success, err := service.Translate(localizer, message)
		assert.Error(t, err)
		assert.False(t, success)
		assert.Empty(t, result)
	})

	t.Run("Translate with nil message", func(t *testing.T) {
		// Expect an error when message is nil
		result, success, err := service.Translate(localizer, nil)
		assert.Error(t, err)
		assert.False(t, success)
		assert.Empty(t, result)
	})

	t.Run("Translate with empty message ID", func(t *testing.T) {
		// Define a message with empty ID
		message := &Message{
			ID: "",
		}
		// Expect an error when message ID is empty
		result, success, err := service.Translate(localizer, message)
		assert.Error(t, err)
		assert.False(t, success)
		assert.Empty(t, result)
	})
}

// TestI18nService_MustTranslate tests the MustTranslate method from I18nService.
func TestI18nService_MustTranslate(t *testing.T) {
	// Setup test suite with translation files
	ts := test.TestSuite{}
	_ = ts.CreateConfigTranslationsFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Create a new I18nService instance
	service, err := NewI18nService(defaultLang, "config/translations", "active")
	assert.NoError(t, err)

	localizer, found, err := service.Localizer(defaultLang)
	assert.NoError(t, err)
	assert.True(t, found)

	t.Run("MustTranslate existing message", func(t *testing.T) {
		// Define a message with ID "hello" and template data
		message := &Message{
			ID: "hello",
			Data: map[string]interface{}{
				"name": "World",
			},
		}
		// Should not panic and return the correct translation
		result := service.MustTranslate(localizer, message)
		assert.Equal(t, "Hello, World!", result)
	})

	t.Run("MustTranslate non-existing message panics", func(t *testing.T) {
		// Define a message with a non-existing ID
		message := &Message{
			ID: "nonexistent",
		}
		// Expect a panic when trying to translate a non-existing message
		assert.Panics(t, func() {
			service.MustTranslate(localizer, message)
		})
	})
}
