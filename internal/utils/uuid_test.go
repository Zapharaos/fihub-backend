package utils

import "testing"

func Test_StringsToUUIDs(t *testing.T) {
	// Test with valid UUID strings
	strings := []string{"550e8400-e29b-41d4-a716-446655440000", "550e8400-e29b-41d4-a716-446655440001"}
	uuids, err := StringsToUUIDs(strings)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(uuids) != 2 {
		t.Fatalf("expected 2 UUIDs, got %d", len(uuids))
	}

	// Test with invalid UUID string
	strings = []string{"invalid-uuid-string"}
	_, err = StringsToUUIDs(strings)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
