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
	defer restore()

	// Access the global repository
	repository := R()
	assert.Equal(t, mockRepository, repository)
}
