package mappers

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
)

// BrokerToProto converts a models.Broker to a protogen.Broker
func BrokerToProto(broker models.Broker) *protogen.Broker {
	return &protogen.Broker{
		Id:       broker.ID.String(),
		Name:     broker.Name,
		ImageId:  broker.ImageID.UUID.String(),
		Disabled: broker.Disabled,
	}
}

// BrokerFromProto converts a protogen.Broker to a models.Broker
func BrokerFromProto(broker *protogen.Broker) models.Broker {
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

// BrokersToProto converts a slice of models.Broker to a slice of protogen.Broker
func BrokersToProto(brokers []models.Broker) []*protogen.Broker {
	protoBrokers := make([]*protogen.Broker, len(brokers))
	for i, broker := range brokers {
		protoBrokers[i] = BrokerToProto(broker)
	}
	return protoBrokers
}

// BrokersFromProto converts a slice of protogen.Broker to a slice of models.Broker
func BrokersFromProto(brokers []*protogen.Broker) []models.Broker {
	protoBrokers := make([]models.Broker, len(brokers))
	for i, broker := range brokers {
		protoBrokers[i] = BrokerFromProto(broker)
	}
	return protoBrokers
}

// BrokerImageToProto converts a models.BrokerImage to a protogen.BrokerImage
func BrokerImageToProto(image models.BrokerImage) *protogen.BrokerImage {
	return &protogen.BrokerImage{
		Id:       image.ID.String(),
		BrokerId: image.BrokerID.String(),
		Name:     image.Name,
		Data:     image.Data,
	}
}

// BrokerImageFromProto converts a protogen.BrokerImage to a models.BrokerImage
func BrokerImageFromProto(image *protogen.BrokerImage) models.BrokerImage {
	return models.BrokerImage{
		ID:       uuid.MustParse(image.GetId()),
		BrokerID: uuid.MustParse(image.GetBrokerId()),
		Name:     image.GetName(),
		Data:     image.GetData(),
	}
}

// BrokerUserToProto converts a models.BrokerUser to a protogen.BrokerUser
func BrokerUserToProto(user models.BrokerUser) *protogen.BrokerUser {
	return &protogen.BrokerUser{
		UserId: user.UserID.String(),
		Broker: BrokerToProto(user.Broker),
	}
}

// BrokerUserFromProto converts a protogen.BrokerUser to a models.BrokerUser
func BrokerUserFromProto(user *protogen.BrokerUser) models.BrokerUser {
	return models.BrokerUser{
		UserID: uuid.MustParse(user.GetUserId()),
		Broker: BrokerFromProto(user.GetBroker()),
	}
}

// BrokerUsersToProto converts a slice of models.BrokerUser to a slice of protogen.BrokerUser
func BrokerUsersToProto(users []models.BrokerUser) []*protogen.BrokerUser {
	protoUsers := make([]*protogen.BrokerUser, len(users))
	for i, user := range users {
		protoUsers[i] = BrokerUserToProto(user)
	}
	return protoUsers
}

// BrokerUsersFromProto converts a slice of protogen.BrokerUser to a slice of models.BrokerUser
func BrokerUsersFromProto(users []*protogen.BrokerUser) []models.BrokerUser {
	protoUsers := make([]models.BrokerUser, len(users))
	for i, user := range users {
		protoUsers[i] = BrokerUserFromProto(user)
	}
	return protoUsers
}
