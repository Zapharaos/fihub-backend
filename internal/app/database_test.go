package app

import (
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestConnectPostgres tests the ConnectPostgres function to ensure that it correctly initializes the database.
// This test only verifies the database.
func TestConnectPostgres(t *testing.T) {
	// Create a full test suite
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	ConnectPostgres()

	// Assertions to verify Database initialization
	assert.NotNil(t, database.DB())
}
