package models

import (
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestBroker_IsValid tests the IsValid method of the Broker struct
func TestBroker_IsValid(t *testing.T) {
	// Define valid values
	validUUID := uuid.New()
	validName := "Valid Broker Name"

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
				ID:   validUUID,
				Name: validName,
			},
			expected: true,
			err:      nil,
		},
		{
			name: "invalid broker with empty name",
			broker: Broker{
				ID:   validUUID,
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

// TestBroker_ToProtogenBroker tests the ToProtogenBroker method of the Broker struct
func TestBroker_ToProtogenBroker(t *testing.T) {
	// Create test broker
	id := uuid.New()
	imageID := uuid.New()
	broker := Broker{
		ID:       id,
		Name:     "Test Broker",
		ImageID:  uuid.NullUUID{UUID: imageID, Valid: true},
		Disabled: true,
	}

	// Convert to protogen broker
	protogenBroker := broker.ToProtogenBroker()

	// Assert values were correctly converted
	assert.Equal(t, id.String(), protogenBroker.Id)
	assert.Equal(t, "Test Broker", protogenBroker.Name)
	assert.Equal(t, imageID.String(), protogenBroker.ImageId)
	assert.Equal(t, true, protogenBroker.Disabled)
}

// TestFromProtogenBroker tests the FromProtogenBroker function
func TestFromProtogenBroker(t *testing.T) {
	// Create test IDs
	id := uuid.New()
	imageID := uuid.New()

	// Create a protogen broker
	protogenBroker := &protogen.Broker{
		Id:       id.String(),
		Name:     "Test Broker",
		ImageId:  imageID.String(),
		Disabled: true,
	}

	// Convert to model broker
	broker := FromProtogenBroker(protogenBroker)

	// Assert values were correctly converted
	assert.Equal(t, id, broker.ID)
	assert.Equal(t, "Test Broker", broker.Name)
	assert.Equal(t, false, broker.ImageID.Valid) // Note: ImageID.Valid is false in conversion
	assert.Equal(t, true, broker.Disabled)
}
