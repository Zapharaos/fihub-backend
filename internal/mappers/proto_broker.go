package mappers

import (
	"github.com/Zapharaos/fihub-backend/gen/go/brokerpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
)

// BrokerToProto converts a models.Broker to a brokerpb.Broker
func BrokerToProto(broker models.Broker) *brokerpb.Broker {
	return &brokerpb.Broker{
		Id:       broker.ID.String(),
		Name:     broker.Name,
		ImageId:  broker.ImageID.UUID.String(),
		Disabled: broker.Disabled,
	}
}

// BrokerFromProto converts a brokerpb.Broker to a models.Broker
func BrokerFromProto(broker *brokerpb.Broker) models.Broker {
	imageId, err := uuid.Parse(broker.GetImageId())
	if err != nil {
		imageId = uuid.Nil
	}

	return models.Broker{
		ID:   uuid.MustParse(broker.GetId()),
		Name: broker.GetName(),
		ImageID: uuid.NullUUID{
			UUID:  imageId,
			Valid: imageId != uuid.Nil,
		},
		Disabled: broker.GetDisabled(),
	}
}

// BrokersToProto converts a slice of models.Broker to a slice of brokerpb.Broker
func BrokersToProto(brokers []models.Broker) []*brokerpb.Broker {
	protoBrokers := make([]*brokerpb.Broker, len(brokers))
	for i, broker := range brokers {
		protoBrokers[i] = BrokerToProto(broker)
	}
	return protoBrokers
}

// BrokersFromProto converts a slice of brokerpb.Broker to a slice of models.Broker
func BrokersFromProto(brokers []*brokerpb.Broker) []models.Broker {
	protoBrokers := make([]models.Broker, len(brokers))
	for i, broker := range brokers {
		protoBrokers[i] = BrokerFromProto(broker)
	}
	return protoBrokers
}

// BrokerImageToProto converts a models.BrokerImage to a brokerpb.BrokerImage
func BrokerImageToProto(image models.BrokerImage) *brokerpb.BrokerImage {
	return &brokerpb.BrokerImage{
		Id:       image.ID.String(),
		BrokerId: image.BrokerID.String(),
		Name:     image.Name,
		Data:     image.Data,
	}
}

// BrokerImageFromProto converts a brokerpb.BrokerImage to a models.BrokerImage
func BrokerImageFromProto(image *brokerpb.BrokerImage) models.BrokerImage {
	return models.BrokerImage{
		ID:       uuid.MustParse(image.GetId()),
		BrokerID: uuid.MustParse(image.GetBrokerId()),
		Name:     image.GetName(),
		Data:     image.GetData(),
	}
}

// BrokerUserToProto converts a models.BrokerUser to a brokerpb.BrokerUser
func BrokerUserToProto(user models.BrokerUser) *brokerpb.BrokerUser {
	return &brokerpb.BrokerUser{
		UserId: user.UserID.String(),
		Broker: BrokerToProto(user.Broker),
	}
}

// BrokerUserFromProto converts a brokerpb.BrokerUser to a models.BrokerUser
func BrokerUserFromProto(user *brokerpb.BrokerUser) models.BrokerUser {
	return models.BrokerUser{
		UserID: uuid.MustParse(user.GetUserId()),
		Broker: BrokerFromProto(user.GetBroker()),
	}
}

// BrokerUsersToProto converts a slice of models.BrokerUser to a slice of brokerpb.BrokerUser
func BrokerUsersToProto(users []models.BrokerUser) []*brokerpb.BrokerUser {
	protoUsers := make([]*brokerpb.BrokerUser, len(users))
	for i, user := range users {
		protoUsers[i] = BrokerUserToProto(user)
	}
	return protoUsers
}

// BrokerUsersFromProto converts a slice of brokerpb.BrokerUser to a slice of models.BrokerUser
func BrokerUsersFromProto(users []*brokerpb.BrokerUser) []models.BrokerUser {
	protoUsers := make([]models.BrokerUser, len(users))
	for i, user := range users {
		protoUsers[i] = BrokerUserFromProto(user)
	}
	return protoUsers
}
