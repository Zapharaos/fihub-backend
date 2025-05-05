package mappers

import (
	"github.com/Zapharaos/fihub-backend/gen/go/brokerpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test_BrokerToProto tests the BrokerToProto method
func Test_BrokerToProto(t *testing.T) {
	// Create test broker
	id := uuid.New()
	imageID := uuid.New()
	broker := models.Broker{
		ID:       id,
		Name:     "Test Broker",
		ImageID:  uuid.NullUUID{UUID: imageID, Valid: true},
		Disabled: true,
	}

	// Convert to gen broker
	protogenBroker := BrokerToProto(broker)

	// Assert values were correctly converted
	assert.Equal(t, id.String(), protogenBroker.Id)
	assert.Equal(t, "Test Broker", protogenBroker.Name)
	assert.Equal(t, imageID.String(), protogenBroker.ImageId)
	assert.Equal(t, true, protogenBroker.Disabled)
}

// Test_BrokerFromProto tests the BrokerFromProto function
func Test_BrokerFromProto(t *testing.T) {
	// Create test IDs
	id := uuid.New()
	imageID := uuid.New()

	// Create a gen broker
	protogenBroker := &brokerpb.Broker{
		Id:       id.String(),
		Name:     "Test Broker",
		ImageId:  imageID.String(),
		Disabled: true,
	}

	// Convert to model broker
	broker := BrokerFromProto(protogenBroker)

	// Assert values were correctly converted
	assert.Equal(t, id, broker.ID)
	assert.Equal(t, "Test Broker", broker.Name)
	assert.Equal(t, true, broker.ImageID.Valid) // Note: ImageID.Valid is false in conversion
	assert.Equal(t, true, broker.Disabled)
}

// Test_BrokersToProto tests the BrokersToProto function
func Test_BrokersToProto(t *testing.T) {
	// Create test brokers
	brokers := []models.Broker{
		{
			ID:       uuid.New(),
			Name:     "Test Broker 1",
			ImageID:  uuid.NullUUID{UUID: uuid.New(), Valid: true},
			Disabled: true,
		},
		{
			ID:       uuid.New(),
			Name:     "Test Broker 2",
			ImageID:  uuid.NullUUID{UUID: uuid.New(), Valid: true},
			Disabled: false,
		},
	}

	// Convert to gen brokers
	protogenBrokers := BrokersToProto(brokers)

	// Assert values were correctly converted
	assert.Equal(t, 2, len(protogenBrokers))
	assert.Equal(t, brokers[0].ID.String(), protogenBrokers[0].Id)
	assert.Equal(t, "Test Broker 1", protogenBrokers[0].Name)
	assert.Equal(t, brokers[0].ImageID.UUID.String(), protogenBrokers[0].ImageId)
	assert.Equal(t, true, protogenBrokers[0].Disabled)

	assert.Equal(t, brokers[1].ID.String(), protogenBrokers[1].Id)
	assert.Equal(t, "Test Broker 2", protogenBrokers[1].Name)
	assert.Equal(t, brokers[1].ImageID.UUID.String(), protogenBrokers[1].ImageId)
	assert.Equal(t, false, protogenBrokers[1].Disabled)
}

// Test_BrokersFromProto tests the BrokersFromProto function
func Test_BrokersFromProto(t *testing.T) {
	// Create test IDs
	id1 := uuid.New()
	id2 := uuid.New()

	// Create a slice of gen brokers
	protogenBrokers := []*brokerpb.Broker{
		{
			Id:       id1.String(),
			Name:     "Test Broker 1",
			ImageId:  id2.String(),
			Disabled: true,
		},
		{
			Id:       id1.String(),
			Name:     "Test Broker 2",
			ImageId:  "bad-uuid",
			Disabled: false,
		},
	}

	// Convert to model brokers
	brokers := BrokersFromProto(protogenBrokers)

	// Assert values were correctly converted
	assert.Equal(t, 2, len(brokers))
	assert.Equal(t, id1, brokers[0].ID)
	assert.Equal(t, "Test Broker 1", brokers[0].Name)
	assert.NotEqual(t, uuid.Nil, brokers[0].ImageID.UUID)
	assert.Equal(t, true, brokers[0].ImageID.Valid)
	assert.Equal(t, true, brokers[0].Disabled)

	assert.Equal(t, id1, brokers[1].ID)
	assert.Equal(t, "Test Broker 2", brokers[1].Name)
	assert.Equal(t, uuid.Nil, brokers[1].ImageID.UUID)
	assert.Equal(t, false, brokers[1].ImageID.Valid)
	assert.Equal(t, false, brokers[1].Disabled)
}

// Test_BrokerImageToProto tests the BrokerImageToProto method
func Test_BrokerImageToProto(t *testing.T) {
	// Create a test BrokerImage
	id := uuid.New()
	brokerID := uuid.New()
	name := "Test Image"
	data := []byte{1, 2, 3, 4, 5}

	brokerImage := models.BrokerImage{
		ID:       id,
		BrokerID: brokerID,
		Name:     name,
		Data:     data,
	}

	// Convert to gen
	protoBrokerImage := BrokerImageToProto(brokerImage)

	// Verify conversion was correct
	assert.Equal(t, id.String(), protoBrokerImage.Id)
	assert.Equal(t, brokerID.String(), protoBrokerImage.BrokerId)
	assert.Equal(t, name, protoBrokerImage.Name)
	assert.Equal(t, data, protoBrokerImage.Data)
}

// Test_BrokerImageFromProto tests the BrokerImageFromProto function
func Test_BrokerImageFromProto(t *testing.T) {
	// Create a test gen.BrokerImage
	id := uuid.New().String()
	brokerID := uuid.New().String()
	name := "Test Protogen Image"
	data := []byte{5, 4, 3, 2, 1}

	protoBrokerImage := &brokerpb.BrokerImage{
		Id:       id,
		BrokerId: brokerID,
		Name:     name,
		Data:     data,
	}

	// Convert from gen
	brokerImage := BrokerImageFromProto(protoBrokerImage)

	// Verify conversion was correct
	assert.Equal(t, id, brokerImage.ID.String())
	assert.Equal(t, brokerID, brokerImage.BrokerID.String())
	assert.Equal(t, name, brokerImage.Name)
	assert.Equal(t, data, brokerImage.Data)
}

// Test_BrokerUserToProto tests the BrokerUserToProto function
func Test_BrokerUserToProto(t *testing.T) {
	userID := uuid.New()
	brokerID := uuid.New()
	brokerName := "Test Broker"
	disabled := false

	brokerUser := models.BrokerUser{
		UserID: userID,
		Broker: models.Broker{
			ID:       brokerID,
			Name:     brokerName,
			Disabled: disabled,
		},
	}

	result := BrokerUserToProto(brokerUser)

	assert.Equal(t, userID.String(), result.UserId)
	assert.Equal(t, brokerID.String(), result.Broker.Id)
	assert.Equal(t, brokerName, result.Broker.Name)
	assert.Equal(t, disabled, result.Broker.Disabled)
}

// Test_BrokerUserFromProto tests the BrokerUserFromProto function
func Test_BrokerUserFromProto(t *testing.T) {
	userID := uuid.New()
	brokerID := uuid.New()
	brokerName := "Test Broker"
	disabled := false

	protoBrokerUser := &brokerpb.BrokerUser{
		UserId: userID.String(),
		Broker: &brokerpb.Broker{
			Id:       brokerID.String(),
			Name:     brokerName,
			ImageId:  uuid.Nil.String(),
			Disabled: disabled,
		},
	}

	result := BrokerUserFromProto(protoBrokerUser)

	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, brokerID, result.Broker.ID)
	assert.Equal(t, brokerName, result.Broker.Name)
	assert.Equal(t, disabled, result.Broker.Disabled)
}

// Test_BrokerUsersToProto tests the BrokerUsersToProto function
func Test_BrokerUsersToProto(t *testing.T) {
	// Create test broker users
	users := []models.BrokerUser{
		{
			UserID: uuid.New(),
			Broker: models.Broker{
				ID:       uuid.New(),
				Name:     "Test Broker 1",
				ImageID:  uuid.NullUUID{UUID: uuid.New(), Valid: true},
				Disabled: true,
			},
		},
		{
			UserID: uuid.New(),
			Broker: models.Broker{
				ID:       uuid.New(),
				Name:     "Test Broker 2",
				ImageID:  uuid.NullUUID{UUID: uuid.New(), Valid: true},
				Disabled: false,
			},
		},
	}

	// Convert to gen broker users
	protogenUsers := BrokerUsersToProto(users)

	// Assert values were correctly converted
	assert.Equal(t, 2, len(protogenUsers))
	assert.Equal(t, users[0].UserID.String(), protogenUsers[0].UserId)
	assert.Equal(t, "Test Broker 1", protogenUsers[0].Broker.Name)
	assert.Equal(t, users[0].Broker.ImageID.UUID.String(), protogenUsers[0].Broker.ImageId)
	assert.Equal(t, true, protogenUsers[0].Broker.Disabled)

	assert.Equal(t, users[1].UserID.String(), protogenUsers[1].UserId)
	assert.Equal(t, "Test Broker 2", protogenUsers[1].Broker.Name)
	assert.Equal(t, users[1].Broker.ImageID.UUID.String(), protogenUsers[1].Broker.ImageId)
	assert.Equal(t, false, protogenUsers[1].Broker.Disabled)
}

// Test_BrokerUsersFromProto tests the BrokerUsersFromProto function
func Test_BrokerUsersFromProto(t *testing.T) {
	// Create a slice of gen broker users
	protogenUsers := []*brokerpb.BrokerUser{
		{
			UserId: uuid.New().String(),
			Broker: &brokerpb.Broker{
				Id:       uuid.New().String(),
				Name:     "Test Broker 1",
				ImageId:  uuid.New().String(),
				Disabled: true,
			},
		},
		{
			UserId: uuid.New().String(),
			Broker: &brokerpb.Broker{
				Id:       uuid.New().String(),
				Name:     "Test Broker 2",
				ImageId:  uuid.New().String(),
				Disabled: false,
			},
		},
	}

	// Convert to model broker users
	users := BrokerUsersFromProto(protogenUsers)

	// Assert values were correctly converted
	assert.Equal(t, 2, len(users))
	assert.Equal(t, protogenUsers[0].UserId, users[0].UserID.String())
	assert.Equal(t, "Test Broker 1", users[0].Broker.Name)
	assert.Equal(t, true, users[0].Broker.Disabled)

	assert.Equal(t, protogenUsers[1].UserId, users[1].UserID.String())
	assert.Equal(t, "Test Broker 2", users[1].Broker.Name)
	assert.Equal(t, false, users[1].Broker.Disabled)
}
