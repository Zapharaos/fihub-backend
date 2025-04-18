package models

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestUserBrokerInput_IsValid tests the IsValid method of the BrokerUserInput struct
func TestUserBrokerInput_IsValid(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string          // Test case name
		input    BrokerUserInput // BrokerUserInput instance to test
		expected bool            // Expected result
		err      error           // Expected error
	}{
		{
			name: "valid input",
			input: BrokerUserInput{
				BrokerID: uuid.New().String(),
			},
			expected: true,
			err:      nil,
		},
		{
			name: "invalid BrokerID",
			input: BrokerUserInput{
				BrokerID: "invalid-uuid",
			},
			expected: false,
			err:      errBrokerIdRequired,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := tt.input.IsValid()
			assert.Equal(t, tt.expected, valid)
			assert.Equal(t, tt.err, err)
		})
	}
}

// TestUserBrokerInput_ToUserBroker tests the ToUser method of the BrokerUserInput struct
func TestUserBrokerInput_ToUserBroker(t *testing.T) {
	brokerID := uuid.New()
	input := BrokerUserInput{
		BrokerID: brokerID.String(),
	}

	expected := BrokerUser{
		Broker: Broker{ID: brokerID},
	}

	result := input.ToUser()
	assert.Equal(t, expected, result)
}
