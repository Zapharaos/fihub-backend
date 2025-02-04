package transactions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestReplaceGlobals tests the ReplaceGlobals function
// It verifies that the global repository can be replaced and restored correctly.
func TestReplaceGlobals(t *testing.T) {
	// Replace the global repository with a mock repository
	mockRepository := &MockRepository{}
	restore := ReplaceGlobals(mockRepository)

	// Ensure the global service is replaced
	assert.Equal(t, mockRepository, R())

	// Restore the previous global service
	restore()
	assert.NotEqual(t, mockRepository, R())
}

// TestRepository tests the R function
// It verifies that the global repository can be accessed correctly.
func TestRepository(t *testing.T) {
	// Replace the global repository with a mock repository
	mockRepository := &MockRepository{}
	restore := ReplaceGlobals(mockRepository)
	defer restore()

	// Access the global repository
	repository := R()
	assert.Equal(t, mockRepository, repository)
}
