package utils

import (
	"testing"
)

// TestRandStringWithCharset tests the RandStringWithCharset function
// It generates a random string of a specified length and charset
// It checks if the length of the generated string is correct
// It also verifies that all characters in the result are from the specified charset
func TestRandStringWithCharset(t *testing.T) {
	length := 10
	charset := "abcd"
	result := RandStringWithCharset(length, charset)
	if len(result) != length {
		t.Errorf("Expected length %d, but got %d", length, len(result))
	}
	for _, char := range result {
		if !contains(charset, char) {
			t.Errorf("Unexpected character %c in result", char)
		}
	}
}

// TestRandString tests the RandString function
// It generates a random string of a specified length using the default charset
// It checks if the length of the generated string is correct
// It also verifies that all characters in the result are from the default charset
func TestRandString(t *testing.T) {
	length := 10
	result := RandString(length)
	if len(result) != length {
		t.Errorf("Expected length %d, but got %d", length, len(result))
	}
	for _, char := range result {
		if !contains(charset, char) {
			t.Errorf("Unexpected character %c in result", char)
		}
	}
}

// TestRandDigitString tests the RandDigitString function
// It generates a random string of a specified length using the digits charset
// It checks if the length of the generated string is correct
// It also verifies that all characters in the result are from the digits charset
func TestRandDigitString(t *testing.T) {
	length := 10
	result := RandDigitString(length)
	if len(result) != length {
		t.Errorf("Expected length %d, but got %d", length, len(result))
	}
	for _, char := range result {
		if !contains(digitsCharset, char) {
			t.Errorf("Unexpected character %c in result", char)
		}
	}
}

func contains(charset string, char rune) bool {
	for _, c := range charset {
		if c == char {
			return true
		}
	}
	return false
}
