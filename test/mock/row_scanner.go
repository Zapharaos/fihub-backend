package mock

import "errors"

// RowScanner is a mock implementation of the RowScanner interface for testing purposes
type RowScanner struct {
	Rows    [][]any
	Current int
}

var ErrRowScannerUnsupportedType = errors.New("unsupported type")
var ErrRowScannerNoMoreRows = errors.New("no more rows")

// Next advances to the next row
func (m *RowScanner) Next() bool {
	if m.Current < len(m.Rows) {
		m.Current++
		return true
	}
	return false
}

// Scan scans the current row into the destination
func (m *RowScanner) Scan(dest ...any) error {
	i := m.Current - 1
	if i < len(m.Rows) {
		for j, v := range m.Rows[i] {
			switch s := v.(type) {
			case int:
				switch d := dest[j].(type) {
				case *int:
					*d = s
					return nil
				default:
					return ErrRowScannerUnsupportedType
				}
			default:
				return ErrRowScannerUnsupportedType
			}
		}
	}
	return ErrRowScannerNoMoreRows
}
