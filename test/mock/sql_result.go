package mock

import (
	"errors"
)

// SQLResult is a mock implementation of sql.Result for testing purposes
type SQLResult struct {
	ExpectRowsAffected int64
}

func (m SQLResult) LastInsertId() (int64, error) {
	return 0, errors.New("not implemented")
}

func (m SQLResult) RowsAffected() (int64, error) {
	return m.ExpectRowsAffected, nil
}
