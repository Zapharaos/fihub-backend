package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewSqlDatabase tests the creation of a new SqlDatabase instance with given credentials.
func TestNewSqlDatabase(t *testing.T) {
	// Create a new SqlDatabase instance
	credentials := SqlCredentials{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpassword",
		DbName:   "testdb",
	}

	// Create a new SqlDatabase instance
	db := NewSqlDatabase(credentials)

	assert.NotNil(t, db, "Expected non-nil SqlDatabase instance")
}
