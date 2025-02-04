package brokers

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestBroker_IsValid tests the IsValid method of the Broker struct
func TestBroker_IsValid(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string // Test case name
		broker   Broker // Broker instance to test
		expected bool   // Expected result
		err      error  // Expected error
	}{
		{
			name: "valid broker",
			broker: Broker{
				ID:   uuid.New(),
				Name: "Valid Broker",
			},
			expected: true,
			err:      nil,
		},
		{
			name: "invalid broker with empty name",
			broker: Broker{
				ID:   uuid.New(),
				Name: "",
			},
			expected: false,
			err:      errBrokerNameRequired,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.broker.IsValid()
			assert.Equal(t, tt.expected, got)
			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
