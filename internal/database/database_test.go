package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDB tests that the DB function returns the global database instance
func TestDB(t *testing.T) {
	// Create a simple test database instance
	testDB := Databases{
		postgres: PostgresDB{},
		redis:    RedisDB{},
	}

	// Replace the global instance
	restore := ReplaceGlobals(testDB)
	defer restore()

	// Verify DB() returns the global instance
	assert.Equal(t, &testDB, DB())
}

// TestDatabaseAccessors tests the accessor methods
func TestDatabaseAccessors(t *testing.T) {
	// Setup simple test database
	pg := PostgresDB{}
	rd := RedisDB{}
	db := Databases{
		postgres: pg,
		redis:    rd,
	}

	// Test accessors
	assert.Equal(t, pg, db.Postgres())
	assert.Equal(t, rd, db.Redis())
}

// TestDatabaseSetters tests the setter methods
func TestDatabaseSetters(t *testing.T) {
	// Initial database
	db := Databases{
		postgres: PostgresDB{},
		redis:    RedisDB{},
	}

	// New instances to set
	newPg := PostgresDB{}
	newRd := RedisDB{}

	// Set new instances
	db.SetPostgres(newPg)
	db.SetRedis(newRd)

	// Verify the values were set
	assert.Equal(t, newPg, db.postgres)
	assert.Equal(t, newRd, db.redis)
}

// TestCloseAll tests the CloseAll method
func TestCloseAll(t *testing.T) {
	// Setup simple test database
	db := Databases{
		postgres: PostgresDB{},
		redis:    RedisDB{},
	}

	// Should not panic
	assert.NotPanics(t, db.CloseAll)
}
