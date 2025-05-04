package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

type MockSql struct {
	ConnectDB    *sqlx.DB
	ConnectError error
}

// Connect simulates a database connection and returns the mocks connection and error.
func (m MockSql) Connect() (*sqlx.DB, error) {
	return m.ConnectDB, m.ConnectError
}

// TestNewPostgresDB tests the NewPostgresDB function with different connection scenarios.
func TestNewPostgresDB(t *testing.T) {
	t.Run("Connection to fail", func(t *testing.T) {
		// Simulate a failed connection
		mock := MockSql{
			ConnectDB:    nil,
			ConnectError: fmt.Errorf("connection failed"),
		}
		db := NewPostgresDB(mock)
		assert.Nil(t, db.DB)
	})

	t.Run("Connection to succeed", func(t *testing.T) {
		// Simulate a successful connection
		sqlxMock, _, err := sqlmock.Newx()
		assert.NoError(t, err)
		defer sqlxMock.Close()

		mock := MockSql{
			ConnectDB:    sqlxMock,
			ConnectError: nil,
		}
		db := NewPostgresDB(mock)
		assert.NotNil(t, db)
	})
}

// TestPostgresDB_IsHealthy tests the IsHealthy function with different connection states.
func TestPostgresDB_IsHealthy(t *testing.T) {
	// Test case for nil DB
	t.Run("Nil DB", func(t *testing.T) {
		db := PostgresDB{DB: nil}
		assert.False(t, db.IsHealthy())
	})

	// Test case for failed ping
	t.Run("Failed ping", func(t *testing.T) {
		sqlxMock, mock, err := sqlmock.Newx(sqlmock.MonitorPingsOption(true))
		assert.NoError(t, err)
		defer sqlxMock.Close()

		// Setup expectation for PingContext to fail
		mock.ExpectPing().WillReturnError(fmt.Errorf("ping failed"))

		db := PostgresDB{DB: sqlxMock}
		assert.False(t, db.IsHealthy())
	})

	// Test case for successful ping
	t.Run("Successful ping", func(t *testing.T) {
		sqlxMock, mock, err := sqlmock.Newx()
		assert.NoError(t, err)
		defer sqlxMock.Close()

		// Setup expectation for PingContext to succeed
		mock.ExpectPing()

		db := PostgresDB{DB: sqlxMock}
		assert.True(t, db.IsHealthy())
	})
}
