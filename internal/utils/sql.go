package utils

import (
	"database/sql"
	"errors"
	"go.uber.org/zap"
)

// RowScanner is an interface that abstracts the sqlx.Rows functionality.
type RowScanner interface {
	Next() bool
	Scan(dest ...interface{}) error
}

var ErrNoRowAffected = errors.New("no row affected (or multiple row affected) instead of 1 row")

// CheckRowAffected checks if the number of rows affected by a query is the expected number
func CheckRowAffected(result sql.Result, nbRows int64) error {
	i, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if i != nbRows {
		return ErrNoRowAffected
	}
	return nil
}

// ScanFirst scans the first row of a sql.Rows and returns the result
func ScanFirst[T any](rows RowScanner, scan func(rows RowScanner) (T, error)) (T, bool, error) {
	if rows.Next() {
		obj, err := scan(rows)
		return obj, err == nil, err
	}
	var a T
	return a, false, nil
}

// ScanAll scans all the rows of the given rows and returns a slice of DataSource
func ScanAll[T any](rows RowScanner, scan func(rows RowScanner) (T, error)) ([]T, error) {
	objs := make([]T, 0)
	for rows.Next() {
		obj, err := scan(rows)
		if err != nil {
			zap.L().Warn("scan error", zap.Error(err))
			return []T{}, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}
