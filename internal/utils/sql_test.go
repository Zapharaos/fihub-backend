package utils

import (
	"database/sql"
	"github.com/Zapharaos/fihub-backend/test/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCheckRowAffected tests the CheckRowAffected function
// It checks if the function correctly identifies the number of rows affected
func TestCheckRowAffected(t *testing.T) {
	tests := []struct {
		name     string
		result   sql.Result
		nbRows   int64
		expected error
	}{
		{
			name:     "Multiple rows affected",
			result:   mock.SQLResult{ExpectRowsAffected: 2},
			nbRows:   1,
			expected: ErrNoRowAffected,
		},
		{
			name:     "No rows affected",
			result:   mock.SQLResult{ExpectRowsAffected: 0},
			nbRows:   1,
			expected: ErrNoRowAffected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckRowAffected(tt.result, tt.nbRows)
			assert.Equal(t, tt.expected, err)
		})
	}
}

// TestScanFirst tests the ScanFirst function
// It checks if the function correctly scans only the first row
func TestScanFirst(t *testing.T) {
	tests := []struct {
		name     string
		rows     mock.RowScanner
		scan     func(rows RowScanner) (int, error)
		expected int
		found    bool
		err      error
	}{
		{
			name: "No rows",
			rows: mock.RowScanner{
				Rows: [][]any{},
			},
			scan: func(rows RowScanner) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: 0,
			found:    false,
			err:      nil,
		}, {
			name: "Single row",
			rows: mock.RowScanner{
				Rows: [][]any{
					{1},
				},
			},
			scan: func(rows RowScanner) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: 1,
			found:    true,
			err:      nil,
		},
		{
			name: "Multiple row",
			rows: mock.RowScanner{
				Rows: [][]any{
					{1}, {2}, {3},
				},
			},
			scan: func(rows RowScanner) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: 1,
			found:    true,
			err:      nil,
		},
		{
			name: "Scan error",
			rows: mock.RowScanner{
				Rows: [][]any{
					{"scan string into int"},
				},
			},
			scan: func(rows RowScanner) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: 0,
			found:    false,
			err:      mock.ErrRowScannerUnsupportedType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, found, err := ScanFirst(&tt.rows, tt.scan)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.found, found)
			assert.Equal(t, tt.err, err)
		})
	}
}

// TestScanAll tests the ScanAll function
// It checks if the function correctly scans all rows
func TestScanAll(t *testing.T) {
	tests := []struct {
		name     string
		rows     mock.RowScanner
		scan     func(rows RowScanner) (int, error)
		expected []int
		err      error
	}{
		{
			name: "No rows",
			rows: mock.RowScanner{
				Rows: [][]any{},
			},
			scan: func(rows RowScanner) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: []int{},
			err:      nil,
		}, {
			name: "Single row",
			rows: mock.RowScanner{
				Rows: [][]any{
					{1},
				},
			},
			scan: func(rows RowScanner) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: []int{1},
			err:      nil,
		},
		{
			name: "Multiple row",
			rows: mock.RowScanner{
				Rows: [][]any{
					{1}, {2}, {3},
				},
			},
			scan: func(rows RowScanner) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: []int{1, 2, 3},
			err:      nil,
		},
		{
			name: "Scan error",
			rows: mock.RowScanner{
				Rows: [][]any{
					{"scan string into int"},
				},
			},
			scan: func(rows RowScanner) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: []int{},
			err:      mock.ErrRowScannerUnsupportedType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ScanAll(&tt.rows, tt.scan)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.err, err)
		})
	}
}
