package utils

import (
	"database/sql"
	"errors"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
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
			name: "With error",
			result: sqlmock.NewErrorResult(
				errors.New("error"),
			),
			nbRows:   0,
			expected: errors.New("error"),
		},
		{
			name:     "No rows affected, unexpected",
			result:   sqlmock.NewResult(0, 0),
			nbRows:   1,
			expected: ErrNoRowAffected,
		},
		{
			name:     "Multiple rows affected, unexpected",
			result:   sqlmock.NewResult(0, 2),
			nbRows:   1,
			expected: ErrNoRowAffected,
		},
		{
			name:     "No rows affected",
			result:   sqlmock.NewResult(0, 0),
			nbRows:   0,
			expected: nil,
		},
		{
			name:     "Multiple rows affected",
			result:   sqlmock.NewResult(0, 2),
			nbRows:   2,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckRowAffected(tt.result, tt.nbRows)
			assert.Equal(t, tt.expected, err)
		})
	}
}

// TestScanString tests the ScanString function
func TestScanString(t *testing.T) {
	s := test.Sqlx{}
	s.CreateFullTestSqlx(t)
	defer s.CleanTestSqlx()

	tests := []struct {
		name     string
		rows     *sqlmock.Rows
		expected string
		err      bool
	}{
		{
			name:     "No rows",
			rows:     sqlmock.NewRows([]string{"id"}),
			expected: "",
			err:      true,
		},
		{
			name:     "Scan error",
			rows:     sqlmock.NewRows([]string{"id"}).AddRow(nil),
			expected: "",
			err:      true,
		},
		{
			name:     "Single row",
			rows:     sqlmock.NewRows([]string{"id"}).AddRow("1"),
			expected: "1",
			err:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := s.MockQuery(tt.rows)
			assert.NoError(t, err)

			if rows.Next() {
				result, err := ScanString(rows)
				assert.Equal(t, tt.expected, result)
				if tt.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

// TestScanFirst tests the ScanFirst function
// It checks if the function correctly scans only the first row
func TestScanFirst(t *testing.T) {

	s := test.Sqlx{}
	s.CreateFullTestSqlx(t)
	defer s.CleanTestSqlx()

	tests := []struct {
		name     string
		rows     *sqlmock.Rows
		scan     func(rows *sqlx.Rows) (int, error)
		expected int
		found    bool
		err      bool
	}{
		{
			name: "No rows",
			rows: sqlmock.NewRows([]string{"id"}),
			scan: func(rows *sqlx.Rows) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: 0,
			found:    false,
			err:      false,
		},
		{
			name: "Single row",
			rows: sqlmock.NewRows([]string{"id"}).AddRow(1),
			scan: func(rows *sqlx.Rows) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: 1,
			found:    true,
			err:      false,
		},
		{
			name: "Multiple row",
			rows: sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2),
			scan: func(rows *sqlx.Rows) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: 1,
			found:    true,
			err:      false,
		},
		{
			name: "Scan error",
			rows: sqlmock.NewRows([]string{"id"}).AddRow("invalid"),
			scan: func(rows *sqlx.Rows) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: 0,
			found:    false,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := s.MockQuery(tt.rows)
			assert.NoError(t, err)

			result, found, err := ScanFirst(rows, tt.scan)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.found, found)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestScanAll tests the ScanAll function
// It checks if the function correctly scans all rows
func TestScanAll(t *testing.T) {
	s := test.Sqlx{}
	s.CreateFullTestSqlx(t)
	defer s.CleanTestSqlx()

	tests := []struct {
		name     string
		rows     *sqlmock.Rows
		scan     func(rows *sqlx.Rows) (int, error)
		expected []int
		err      bool
	}{
		{
			name: "No rows",
			rows: sqlmock.NewRows([]string{"id"}),
			scan: func(rows *sqlx.Rows) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: []int{},
			err:      false,
		}, {
			name: "Single row",
			rows: sqlmock.NewRows([]string{"id"}).AddRow(1),
			scan: func(rows *sqlx.Rows) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: []int{1},
			err:      false,
		},
		{
			name: "Multiple row",
			rows: sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2).AddRow(3),
			scan: func(rows *sqlx.Rows) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: []int{1, 2, 3},
			err:      false,
		},
		{
			name: "Scan error",
			rows: sqlmock.NewRows([]string{"id"}).AddRow("invalid"),
			scan: func(rows *sqlx.Rows) (int, error) {
				var value int
				err := rows.Scan(&value)
				return value, err
			},
			expected: []int{},
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := s.MockQuery(tt.rows)
			assert.NoError(t, err)

			result, err := ScanAll(rows, tt.scan)
			assert.Equal(t, tt.expected, result)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestScanFirstStruct(t *testing.T) {
	s := test.Sqlx{}
	s.CreateFullTestSqlx(t)
	defer s.CleanTestSqlx()

	type TestStruct struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
		Skip string
	}

	tests := []struct {
		name     string
		rows     *sqlmock.Rows
		expected TestStruct
		found    bool
		err      bool
	}{
		{
			name:     "Scan error",
			rows:     sqlmock.NewRows([]string{"id", "name"}).AddRow("invalid", nil),
			expected: TestStruct{},
			found:    false,
			err:      true,
		},
		{
			name:     "No rows",
			rows:     sqlmock.NewRows([]string{"id", "name"}),
			expected: TestStruct{},
			found:    false,
			err:      false,
		},
		{
			name:     "Single row",
			rows:     sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test"),
			expected: TestStruct{ID: 1, Name: "test"},
			found:    true,
			err:      false,
		},
		{
			name:     "Multiple row",
			rows:     sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test1").AddRow(2, "test2"),
			expected: TestStruct{ID: 1, Name: "test1"},
			found:    true,
			err:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := s.MockQuery(tt.rows)
			assert.NoError(t, err)

			result, found, err := ScanFirstStruct[TestStruct](rows)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.found, found)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestScanAllStruct(t *testing.T) {
	s := test.Sqlx{}
	s.CreateFullTestSqlx(t)
	defer s.CleanTestSqlx()

	type TestStruct struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
		Skip string
	}

	tests := []struct {
		name     string
		rows     *sqlmock.Rows
		expected []TestStruct
		err      bool
	}{
		{
			name:     "Scan error",
			rows:     sqlmock.NewRows([]string{"id", "name"}).AddRow("invalid", nil),
			expected: []TestStruct{},
			err:      true,
		},
		{
			name:     "No rows",
			rows:     sqlmock.NewRows([]string{"id", "name"}),
			expected: []TestStruct(nil),
			err:      false,
		},
		{
			name:     "Single row",
			rows:     sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test"),
			expected: []TestStruct{{ID: 1, Name: "test"}},
			err:      false,
		},
		{
			name:     "Multiple row",
			rows:     sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test1").AddRow(2, "test2"),
			expected: []TestStruct{{ID: 1, Name: "test1"}, {ID: 2, Name: "test2"}},
			err:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := s.MockQuery(tt.rows)
			assert.NoError(t, err)

			result, err := ScanAllStruct[TestStruct](rows)
			assert.Equal(t, tt.expected, result)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
