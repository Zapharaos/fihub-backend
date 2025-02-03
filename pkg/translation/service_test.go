package translation

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

// TranslationService is a mock implementation of the translation.Service interface
type TranslationService struct {
	Messages map[language.Tag]map[string]string
}

var ErrLocalizerNotFound = errors.New("localizer not found")

// NewTranslationService creates a new instance of TranslationService
func NewTranslationService(messages map[language.Tag]map[string]string) Service {
	s := TranslationService{
		Messages: messages,
	}
	var service Service = &s
	return service
}

// Localizer returns a mock localizer
func (t TranslationService) Localizer(language language.Tag) (interface{}, error) {
	localizer, found := t.Messages[language]
	if !found {
		return nil, ErrLocalizerNotFound
	}
	return localizer, nil
}

// Message returns a mock localized message
func (t TranslationService) Message(localizer interface{}, message *Message) string {
	// Verify that the localizer is of the correct type
	loc, ok := localizer.(map[string]string)
	if !ok {
		return ""
	}

	if msg, ok := loc[message.ID]; ok {
		return msg
	}
	return message.ID
}

// TestTranslationReplaceGlobals tests the ReplaceGlobals function
// It verifies that the global service can be replaced and restored correctly.
func TestTranslationReplaceGlobals(t *testing.T) {
	// Replace the global service with a mock service
	mockService := &TranslationService{}
	restore := ReplaceGlobals(mockService)

	// Ensure the global service is replaced
	assert.Equal(t, mockService, S())

	// Restore the previous global service
	restore()
	assert.NotEqual(t, mockService, S())
}

// TestTranslationS tests the S function
// It verifies that the global service can be accessed correctly.
func TestTranslationS(t *testing.T) {
	// Replace the global service with a mock service
	mockService := &TranslationService{}
	restore := ReplaceGlobals(mockService)
	defer restore()

	// Access the global service
	service := S()
	assert.Equal(t, mockService, service)
}
