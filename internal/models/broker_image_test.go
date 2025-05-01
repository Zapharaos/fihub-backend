package models

import (
	"github.com/Zapharaos/fihub-backend/gen"
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

// TestBrokerImage_ToProtogenBrokerImage tests the ToProtogenBrokerImage method
func TestBrokerImage_ToProtogenBrokerImage(t *testing.T) {
	// Create a test BrokerImage
	id := uuid.New()
	brokerID := uuid.New()
	name := "Test Image"
	data := []byte{1, 2, 3, 4, 5}

	brokerImage := BrokerImage{
		ID:       id,
		BrokerID: brokerID,
		Name:     name,
		Data:     data,
	}

	// Convert to gen
	protoBrokerImage := brokerImage.ToProtogenBrokerImage()

	// Verify conversion was correct
	assert.Equal(t, id.String(), protoBrokerImage.Id)
	assert.Equal(t, brokerID.String(), protoBrokerImage.BrokerId)
	assert.Equal(t, name, protoBrokerImage.Name)
	assert.Equal(t, data, protoBrokerImage.Data)
}

// TestFromProtogenBrokerImage tests the FromProtogenBrokerImage function
func TestFromProtogenBrokerImage(t *testing.T) {
	// Create a test gen.BrokerImage
	id := uuid.New().String()
	brokerID := uuid.New().String()
	name := "Test Protogen Image"
	data := []byte{5, 4, 3, 2, 1}

	protoBrokerImage := &protogen.BrokerImage{
		Id:       id,
		BrokerId: brokerID,
		Name:     name,
		Data:     data,
	}

	// Convert from gen
	brokerImage := FromProtogenBrokerImage(protoBrokerImage)

	// Verify conversion was correct
	assert.Equal(t, id, brokerImage.ID.String())
	assert.Equal(t, brokerID, brokerImage.BrokerID.String())
	assert.Equal(t, name, brokerImage.Name)
	assert.Equal(t, data, brokerImage.Data)
}
