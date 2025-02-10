package brokers_test

import (
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewRepository tests the NewRepository function
// It verifies that the repositories are correctly assigned.
func TestNewRepository(t *testing.T) {

	// Replace with mocks repositories
	mockBrokerRepository := &mocks.BrokerRepository{}
	mockUserRepository := &mocks.BrokerUserRepository{}
	mockImageRepository := &mocks.BrokerImageRepository{}

	// Create a new repository
	repo := brokers.NewRepository(mockBrokerRepository, mockUserRepository, mockImageRepository)

	// Verify that the repositories are correctly assigned
	assert.Equal(t, mockBrokerRepository, repo.B())
	assert.Equal(t, mockUserRepository, repo.U())
	assert.Equal(t, mockImageRepository, repo.I())
}

// TestReplaceGlobals tests the ReplaceGlobals function
// It verifies that the global repository can be replaced and restored correctly.
func TestReplaceGlobals(t *testing.T) {
	// Replace with mocks repositories
	mockBrokerRepository := &mocks.BrokerRepository{}
	mockUserRepository := &mocks.BrokerUserRepository{}
	mockImageRepository := &mocks.BrokerImageRepository{}
	mockRepository := brokers.NewRepository(mockBrokerRepository, mockUserRepository, mockImageRepository)

	// Replace the global repository with a mocks repository
	restore := brokers.ReplaceGlobals(mockRepository)

	// Verify that the global repository instance has been replaced
	assert.Equal(t, mockRepository, brokers.R())

	// Restore the global repository instance
	restore()

	// Verify that the global repository instance has been restored
	assert.NotEqual(t, mockRepository, brokers.R())
}
