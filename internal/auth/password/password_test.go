package password

import (
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestInitRequest tests the InitRequest function
func TestInitRequest(t *testing.T) {
	// Mock the viper configuration
	viper.Set("OTP_DURATION", 15*time.Minute)
	viper.Set("OTP_LENGTH", 6)

	userID := uuid.New()
	request, duration := InitRequest(userID)

	assert.NotNil(t, request.ID)
	assert.Equal(t, userID, request.UserID)
	assert.NotEmpty(t, request.Token)
	assert.WithinDuration(t, time.Now().Add(duration), request.ExpiresAt, time.Second)
	assert.WithinDuration(t, time.Now(), request.CreatedAt, time.Second)
}
