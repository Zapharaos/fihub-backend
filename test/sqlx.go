package test

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

type Sqlx struct {
	mock sqlmock.Sqlmock
	sql  *sql.DB
	sqlx *sqlx.DB
}

// CreateFullTestSqlx Create a full test suite
// Please clean test suite after use (defer CleanTestSqlx())
func (s *Sqlx) CreateFullTestSqlx(t *testing.T) {
	// Create a new sqlmock
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	s.mock = mock
	s.sql = sqlDB
	s.sqlx = sqlx.NewDb(sqlDB, "sqlmock")
}

// CleanTestSqlx Clean the test suite
func (s *Sqlx) CleanTestSqlx() {
	s.sql.Close()
	s.sqlx.Close()
}

func (s *Sqlx) MockQuery(rows *sqlmock.Rows) (*sqlx.Rows, error) {
	query := "SELECT"
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	return s.sqlx.Queryx(query)
}
