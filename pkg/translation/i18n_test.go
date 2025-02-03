package translation

import (
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

var defaultLang = language.English

func TestNewI18nService(t *testing.T) {
	t.Run("Without translation files", func(t *testing.T) {
		assert.Panics(t, func() {
			NewI18nService(defaultLang)
		})
	})

	t.Run("With translation files", func(t *testing.T) {
		ts := test.TestSuite{}
		_ = ts.CreateConfigTranslationsFullTestSuite(t)
		defer ts.CleanTestSuite(t)

		service := NewI18nService(defaultLang)
		assert.NotNil(t, service)
	})
}

func TestI18nService_Localizer(t *testing.T) {
	ts := test.TestSuite{}
	_ = ts.CreateConfigTranslationsFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	service := NewI18nService(defaultLang)

	t.Run("Retrieve default language localizer", func(t *testing.T) {
		localizer, err := service.Localizer(defaultLang)
		assert.NoError(t, err)
		assert.NotNil(t, localizer)
	})

	t.Run("Retrieve non-existing language localizer", func(t *testing.T) {
		_, err := service.Localizer(language.Spanish)
		assert.Error(t, err)
	})
}

func TestI18nService_Message(t *testing.T) {
	ts := test.TestSuite{}
	_ = ts.CreateConfigTranslationsFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	service := NewI18nService(defaultLang)
	localizer, _ := service.Localizer(defaultLang)

	t.Run("Retrieve existing message", func(t *testing.T) {
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
		message := &Message{
			ID: "nonexistent",
		}
		result := service.Message(localizer, message)
		assert.Equal(t, "", result)
	})
}
