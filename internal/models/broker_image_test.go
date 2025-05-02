package models

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TestBrokerImage_IsValid tests the IsValid method of the BrokerImage struct
func TestBrokerImage_IsValid(t *testing.T) {
	// Define valid values
	validUUID := uuid.New()
	validName := "Valid BrokerImage Name"
	validData := []byte{1, 2, 3}

	// Define test cases
	tests := []struct {
		name     string      // Test case name
		image    BrokerImage // BrokerImage instance to test
		expected bool        // Expected result
		err      error       // Expected error
	}{
		{
			name: "valid image",
			image: BrokerImage{
				BrokerID: validUUID,
				Name:     validName,
				Data:     validData,
			},
			expected: true,
			err:      nil,
		},
		{
			name: "missing broker ID",
			image: BrokerImage{
				BrokerID: uuid.Nil,
				Name:     validName,
				Data:     validData,
			},
			expected: false,
			err:      errBrokerIdRequired,
		},
		{
			name: "invalid image name (too short)",
			image: BrokerImage{
				BrokerID: validUUID,
				Name:     strings.Repeat("a", ImageNameMinLength-1),
				Data:     validData,
			},
			expected: false,
			err:      errImageNameInvalid,
		},
		{
			name: "invalid image name (too long)",
			image: BrokerImage{
				BrokerID: uuid.New(),
				Name:     strings.Repeat("a", ImageNameMaxLength+1),
				Data:     validData,
			},
			expected: false,
			err:      errImageNameInvalid,
		},
		{
			name: "missing image data",
			image: BrokerImage{
				BrokerID: uuid.New(),
				Name:     validName,
				Data:     []byte{},
			},
			expected: false,
			err:      errImageDataRequired,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := tt.image.IsValid()
			assert.Equal(t, tt.expected, valid)
			assert.Equal(t, tt.err, err)
		})
	}
}
