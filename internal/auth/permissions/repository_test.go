package permissions_test

import (
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestReplaceGlobals tests the ReplaceGlobals function
// It verifies that the global repository can be replaced and restored correctly.
func TestReplaceGlobals(t *testing.T) {
	// Replace the global repository with a mocks repository
	mockRepository := &mocks.PermissionsRepository{}
	restore := permissions.ReplaceGlobals(mockRepository)

	// Verify that the global repository instance has been replaced
	assert.Equal(t, mockRepository, permissions.R())

	// Restore the global repository instance
	restore()

	// Verify that the global repository instance has been restored
	assert.NotEqual(t, mockRepository, permissions.R())
}

// TestRepository tests the R function
// It verifies that the global repository can be accessed correctly.
func TestRepository(t *testing.T) {
	// Replace the global repository with a mocks repository
	mockRepository := &mocks.PermissionsRepository{}
	restore := permissions.ReplaceGlobals(mockRepository)
	defer restore()

	// Access the global repository
	assert.Equal(t, mockRepository, permissions.R())
}
