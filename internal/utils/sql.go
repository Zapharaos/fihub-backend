package utils

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrNoRowAffected = errors.New("no row affected (or multiple row affected) instead of 1 row")
)

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

// ScanString scans a single string from the given sqlx.Rows
func ScanString(rows *sqlx.Rows) (string, error) {
	var result string
	if err := rows.Scan(&result); err != nil {
		return "", err
	}
	return result, nil
}

// ScanFirst scans the first row of a sql.Rows and returns the result
func ScanFirst[T any](rows *sqlx.Rows, scan func(rows *sqlx.Rows) (T, error)) (T, bool, error) {
	if rows.Next() {
		obj, err := scan(rows)
		return obj, err == nil, err
	}
	var a T
	return a, false, nil
}

// ScanAll scans all the rows of the given rows and returns a slice of T
func ScanAll[T any](rows *sqlx.Rows, scan func(rows *sqlx.Rows) (T, error)) ([]T, error) {
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

// ScanFirstStruct scans the first row of the given rows and returns a struct of type T
func ScanFirstStruct[T any](rows *sqlx.Rows) (T, bool, error) {
	var results []T
	err := sqlx.StructScan(rows, &results)
	if err != nil {
		var a T
		return a, false, err
	}

	if len(results) == 0 {
		var a T
		return a, false, nil
	}

	return results[0], true, nil
}

// ScanAllStruct scans all the rows of the given rows and returns a slice of T
func ScanAllStruct[T any](rows *sqlx.Rows) ([]T, error) {
	var results []T
	err := sqlx.StructScan(rows, &results)
	if err != nil {
		return []T{}, err
	}
	return results, nil
}
