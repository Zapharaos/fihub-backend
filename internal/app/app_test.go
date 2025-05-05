package app

import "testing"

// TestRecoverPanic tests the RecoverPanic function.
func TestRecoverPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("RecoverPanic did not re-panic as expected")
		}
	}()

	// Simulate a panic and recover it
	defer RecoverPanic()
	panic("test panic")
}
