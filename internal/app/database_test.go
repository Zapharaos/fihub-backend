package app

import (
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/users"
	"github.com/Zapharaos/fihub-backend/internal/users/password"
	"github.com/Zapharaos/fihub-backend/internal/users/permissions"
	"github.com/Zapharaos/fihub-backend/internal/users/roles"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

// TestInit tests the initDatabase function to ensure that it correctly initializes the database.
// This test only verifies the database.
// Any further function calls within initDatabase are tested by their respective tests.
func TestInitDatabase(t *testing.T) {
	// Create a full test suite
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Call initDatabase function
	InitDatabase()

	// Assertions to verify Database initialization
	assert.NotNil(t, database.DB())
}

// TestInitPostgres tests the initPostgres function to ensure that it correctly initializes the repositories.
func TestInitPostgres(t *testing.T) {
	// Simulate a successful connection
	sqlxMock, _, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer sqlxMock.Close()

	// Call initPostgres function with the mock connection
	InitPostgres(sqlxMock)

	// Assertions to verify repositories initialization
	assert.NotNil(t, users.R(), "Users repository should be initialized")
	assert.NotNil(t, password.R(), "Password repository should be initialized")
	assert.NotNil(t, roles.R(), "Roles repository should be initialized")
	assert.NotNil(t, permissions.R(), "Permissions repository should be initialized")
	assert.NotNil(t, brokers.R(), "Brokers repository should be initialized")
}
