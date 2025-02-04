package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhashkevych/go-sqlxmock"
)

// TestNewDatabases tests the creation of NewDatabases non-nil instances.
func TestNewDatabases(t *testing.T) {
	postgres, _, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer postgres.Close()

	databases := NewDatabases(postgres)
	assert.NotNil(t, databases.Postgres())
}

// TestReplaceGlobals tests the ReplaceGlobals function.
// Verifies that the ReplaceGlobals function correctly replaces the
// global database instance and restores it after the test.
func TestReplaceGlobals(t *testing.T) {
	db, _, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	databases := NewDatabases(db)
	restore := ReplaceGlobals(databases)
	assert.Equal(t, databases, DB())

	restore()
	assert.NotEqual(t, databases, DB())
}
