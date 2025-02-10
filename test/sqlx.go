package test

import (
	"github.com/jmoiron/sqlx"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

type Sqlx struct {
	Mock sqlxmock.Sqlmock
	DB   *sqlx.DB
}

// CreateFullTestSqlx Create a full test suite
// Please clean test suite after use (defer CleanTestSqlx())
func (s *Sqlx) CreateFullTestSqlx(t *testing.T) {
	var err error
	s.DB, s.Mock, err = sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
}

// CleanTestSqlx Clean the test suite
func (s *Sqlx) CleanTestSqlx() {
	s.DB.Close()
}

func (s *Sqlx) MockQuery(rows *sqlxmock.Rows) (*sqlx.Rows, error) {
	query := "SELECT"
	s.Mock.ExpectQuery(query).WillReturnRows(rows)
	return s.DB.Queryx(query)
}
