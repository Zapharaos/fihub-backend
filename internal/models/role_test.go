package models

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TestRole_IsValid tests the IsValid method of the Role struct
func TestRole_IsValid(t *testing.T) {
	// Define valid values
	validUUID := uuid.New()
	validName := strings.Repeat("a", NameMinLength)

	// Define test cases
	tests := []struct {
		name     string
		role     Role
		expected bool
		err      error
	}{
		{
			name:     "valid role",
			role:     Role{Id: validUUID, Name: validName},
			expected: true,
			err:      nil,
		},
		{
			name:     "empty name",
			role:     Role{Id: validUUID, Name: ""},
			expected: false,
			err:      ErrNameRequired,
		},
		{
			name:     "name too short",
			role:     Role{Id: validUUID, Name: strings.Repeat("a", NameMinLength-1)},
			expected: false,
			err:      ErrNameInvalid,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := tt.role.IsValid()
			assert.Equal(t, tt.expected, valid)
			assert.Equal(t, tt.err, err)
		})
	}
}

// TestRolesWithPermissions_GetUUIDs tests the GetUUIDs method of the RolesWithPermissions struct
func TestRolesWithPermissions_GetUUIDs(t *testing.T) {
	// Define valid values
	validUUID1 := uuid.New()
	validUUID2 := uuid.New()

	// Define test cases
	tests := []struct {
		name     string
		roles    RolesWithPermissions
		expected []uuid.UUID
	}{
		{
			name: "multiple roles",
			roles: RolesWithPermissions{
				{Role: Role{Id: validUUID1}},
				{Role: Role{Id: validUUID2}},
			},
			expected: []uuid.UUID{validUUID1, validUUID2},
		},
		{
			name:     "no roles",
			roles:    RolesWithPermissions{},
			expected: []uuid.UUID{},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuids := tt.roles.GetUUIDs()
			assert.Equal(t, tt.expected, uuids)
		})
	}
}
