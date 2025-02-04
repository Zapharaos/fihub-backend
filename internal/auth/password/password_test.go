package password

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestInitRequest tests the InitRequest function
func TestInitRequest(t *testing.T) {
	userID := uuid.New()
	request, duration := InitRequest(userID)

	assert.NotNil(t, request.ID)
	assert.Equal(t, userID, request.UserID)
	assert.NotEmpty(t, request.Token)
	assert.WithinDuration(t, time.Now().Add(duration), request.ExpiresAt, time.Second)
	assert.WithinDuration(t, time.Now(), request.CreatedAt, time.Second)
}
