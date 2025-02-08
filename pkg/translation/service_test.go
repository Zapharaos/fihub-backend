package translation

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

// TestTranslationReplaceGlobals tests the ReplaceGlobals function
// It verifies that the global service can be replaced and restored correctly.
func TestTranslationReplaceGlobals(t *testing.T) {
	// Mock
	ctrl := gomock.NewController(t)
	m := NewMockService(ctrl)
	defer ctrl.Finish()

	// Replace the global service with a mock service
	restore := ReplaceGlobals(m)

	// Ensure the global service is replaced
	assert.Equal(t, m, S())

	// Restore the previous global service
	restore()
	assert.NotEqual(t, m, S())
}

// TestTranslationS tests the S function
// It verifies that the global service can be accessed correctly.
func TestTranslationS(t *testing.T) {
	// Mock
	ctrl := gomock.NewController(t)
	m := NewMockService(ctrl)
	defer ctrl.Finish()

	// Replace the global service with a mock service
	restore := ReplaceGlobals(m)
	defer restore()

	// Access the global service
	service := S()
	assert.Equal(t, m, service)
}
