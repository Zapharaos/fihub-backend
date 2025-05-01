package models

import (
	"github.com/Zapharaos/fihub-backend/gen"
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

// TestBrokerUser_ToProtogenBrokerUser tests the ToProtogenBrokerUser method of the BrokerUser struct
func TestBrokerUser_ToProtogenBrokerUser(t *testing.T) {
	userID := uuid.New()
	brokerID := uuid.New()
	brokerName := "Test Broker"
	disabled := false

	brokerUser := BrokerUser{
		UserID: userID,
		Broker: Broker{
			ID:       brokerID,
			Name:     brokerName,
			Disabled: disabled,
		},
	}

	result := brokerUser.ToProtogenBrokerUser()

	assert.Equal(t, userID.String(), result.UserId)
	assert.Equal(t, brokerID.String(), result.Broker.Id)
	assert.Equal(t, brokerName, result.Broker.Name)
	assert.Equal(t, disabled, result.Broker.Disabled)
}

// TestFromProtogenBrokerUser tests the FromProtogenBrokerUser function
func TestFromProtogenBrokerUser(t *testing.T) {
	userID := uuid.New()
	brokerID := uuid.New()
	brokerName := "Test Broker"
	disabled := false

	protoBrokerUser := &protogen.BrokerUser{
		UserId: userID.String(),
		Broker: &protogen.Broker{
			Id:       brokerID.String(),
			Name:     brokerName,
			ImageId:  uuid.Nil.String(),
			Disabled: disabled,
		},
	}

	result := FromProtogenBrokerUser(protoBrokerUser)

	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, brokerID, result.Broker.ID)
	assert.Equal(t, brokerName, result.Broker.Name)
	assert.Equal(t, disabled, result.Broker.Disabled)
}
