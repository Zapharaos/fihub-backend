package app

import (
	"github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/password"
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
	assert.NotNil(t, repositories.R(), "Users repository should be initialized")
	assert.NotNil(t, password.R(), "Password repository should be initialized")
	assert.NotNil(t, repositories.R(), "Roles repository should be initialized")
	assert.NotNil(t, repositories.R(), "Permissions repository should be initialized")
}
