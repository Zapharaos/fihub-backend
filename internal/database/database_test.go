package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhashkevych/go-sqlxmock"
)

// TestNewDatabases tests the creation of NewDatabases non-nil instances.
func TestNewDatabases(t *testing.T) {
	// Mock a postgres database instance
	postgres, _, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer postgres.Close()

	// Create a new Databases instance
	databases := NewDatabases(PostgresDB{
		DB: postgres,
	})
	assert.NotNil(t, databases.Postgres().DB)
}

// TestReplaceGlobals tests the ReplaceGlobals function.
// Verifies that the ReplaceGlobals function correctly replaces the
// global database instance and restores it after the test.
func TestReplaceGlobals(t *testing.T) {
	// Mock a postgres database instance
	db, _, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	// Create a new Databases instance
	databases := NewDatabases(PostgresDB{
		DB: db,
	})
	restore := ReplaceGlobals(databases)

	// Verify that the global database instance has been replaced
	assert.Equal(t, databases, DB())

	// Restore the global database instance
	restore()

	// Verify that the global database instance has been restored
	assert.NotEqual(t, databases, DB())
}
