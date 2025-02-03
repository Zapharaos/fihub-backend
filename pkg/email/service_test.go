package email

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// EmailService is a mock implementation of the email.Service interface
type EmailService struct {
	SendError error
}

// NewEmailService creates a new instance of the email.Service
func NewEmailService(sendError error) Service {
	s := EmailService{
		SendError: sendError,
	}
	var service Service = &s
	return service
}

func (e EmailService) Send(emailTo, subject, plainTextContent, htmlContent string) error {
	return e.SendError
}

// TestTranslationReplaceGlobals tests the ReplaceGlobals function
// It verifies that the global service can be replaced and restored correctly.
func TestEmailReplaceGlobals(t *testing.T) {
	// Replace the global service with a mock service
	mockService := &EmailService{}
	restore := ReplaceGlobals(mockService)

	// Ensure the global service is replaced
	assert.Equal(t, mockService, S())

	// Restore the previous global service
	restore()
	assert.NotEqual(t, mockService, S())
}

// TestTranslationS tests the S function
// It verifies that the global service can be accessed correctly.
func TestEmailS(t *testing.T) {
	// Replace the global service with a mock service
	mockService := &EmailService{}
	restore := ReplaceGlobals(mockService)
	defer restore()

	// Access the global service
	service := S()
	assert.Equal(t, mockService, service)
}
