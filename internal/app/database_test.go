package app

import (
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestInitPostgres tests the InitPostgres function to ensure that it correctly initializes the database.
// This test only verifies the database.
func TestInitPostgres(t *testing.T) {
	// Create a full test suite
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	InitPostgres()

	// Assertions to verify Database initialization
	assert.NotNil(t, database.DB())
	assert.NotNil(t, database.DB().Postgres())
}

// TestInitRedis tests the InitRedis function to ensure that it correctly initializes the database.
// This test only verifies the database.
func TestInitRedis(t *testing.T) {
	// Create a full test suite
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	InitRedis()

	// Assertions to verify Database initialization
	assert.NotNil(t, database.DB())
	assert.NotNil(t, database.DB().Redis())
}
