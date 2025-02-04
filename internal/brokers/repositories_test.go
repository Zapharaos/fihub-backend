package brokers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewRepository tests the NewRepository function
// It verifies that the repositories are correctly assigned.
func TestNewRepository(t *testing.T) {

	// Replace with mock repositories
	mockBrokerRepository := &MockBrokerRepository{}
	mockUserRepository := &MockUserBrokerRepository{}
	mockImageRepository := &MockImageRepository{}

	// Create a new repository
	repo := NewRepository(mockBrokerRepository, mockUserRepository, mockImageRepository)

	// Verify that the repositories are correctly assigned
	assert.Equal(t, mockBrokerRepository, repo.B())
	assert.Equal(t, mockUserRepository, repo.U())
	assert.Equal(t, mockImageRepository, repo.I())
}

// TestReplaceGlobals tests the ReplaceGlobals function
// It verifies that the global repository can be replaced and restored correctly.
func TestReplaceGlobals(t *testing.T) {
	// Replace with mock repositories
	mockBrokerRepository := &MockBrokerRepository{}
	mockUserRepository := &MockUserBrokerRepository{}
	mockImageRepository := &MockImageRepository{}
	mockRepository := NewRepository(mockBrokerRepository, mockUserRepository, mockImageRepository)

	// Replace the global repository with a mock repository
	restore := ReplaceGlobals(mockRepository)
	
	// Verify that the global repository instance has been replaced
	assert.Equal(t, mockRepository, R())

	// Restore the global repository instance
	restore()

	// Verify that the global repository instance has been restored
	assert.NotEqual(t, mockRepository, R())
}
