package utils

import "github.com/google/uuid"

// StringsToUUIDs converts a slice of strings to a slice of UUIDs
// Returns the UUIDs and any errors encountered during conversion
func StringsToUUIDs(strings []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, len(strings))
	for i, s := range strings {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		uuids[i] = id
	}
	return uuids, nil
}
