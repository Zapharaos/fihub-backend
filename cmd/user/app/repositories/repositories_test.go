package repositories_test

import (
	"github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewRepository tests the NewRepository function
// It verifies that the repositories are correctly assigned.
func TestNewRepository(t *testing.T) {

	// Replace with mocks repositories
	mockUserRepository := &mocks.UserRepository{}
	mockRoleRepository := &mocks.UserRoleRepository{}
	mockPermissionRepository := &mocks.UserPermissionRepository{}

	// Create a new repository
	repo := repositories.NewRepository(mockUserRepository, mockRoleRepository, mockPermissionRepository)

	// Verify that the repositories are correctly assigned
	assert.Equal(t, mockUserRepository, repo.U())
}

// TestReplaceGlobals tests the ReplaceGlobals function
// It verifies that the global repository can be replaced and restored correctly.
func TestReplaceGlobals(t *testing.T) {
	// Replace with mocks repositories
	mockUserRepository := &mocks.UserRepository{}
	mockRoleRepository := &mocks.UserRoleRepository{}
	mockPermissionRepository := &mocks.UserPermissionRepository{}
	mockRepository := repositories.NewRepository(mockUserRepository, mockRoleRepository, mockPermissionRepository)

	// Replace the global repository with a mocks repository
	restore := repositories.ReplaceGlobals(mockRepository)

	// Verify that the global repository instance has been replaced
	assert.Equal(t, mockRepository, repositories.R())

	// Restore the global repository instance
	restore()

	// Verify that the global repository instance has been restored
	assert.NotEqual(t, mockRepository, repositories.R())
}
