package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewOtpTemplate tests the NewOtpTemplate function
func TestNewOtpTemplate(t *testing.T) {
	data := OtpData{
		OTP:         "123456",
		Greeting:    "Hello",
		MainContent: "Your OTP is below",
		DoNotShare:  "Do not share this OTP with anyone.",
	}

	template := NewOtpTemplate(data)

	assert.Equal(t, "otp", template.Name)
	assert.Equal(t, otpHtml, template.ContentRaw)
	assert.Equal(t, data, template.Data)
}
