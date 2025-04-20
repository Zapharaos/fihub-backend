package models

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
)

var (
	errBrokerIdRequired   = errors.New("broker-required")
	errBrokerNameRequired = errors.New("name-required")
)

// Broker represents a broker entity in the system
type Broker struct {
	ID       uuid.UUID     `json:"id"`
	Name     string        `json:"name"`
	ImageID  uuid.NullUUID `json:"image_id" swaggertype:"string"`
	Disabled bool          `json:"disabled"`
}

// IsValid checks if the Broker is valid
func (b Broker) IsValid() (bool, error) {
	if b.Name == "" {
		return false, errBrokerNameRequired
	}
	return true, nil
}

// ToProtogenBroker converts a Broker to a protogen.Broker
func (b Broker) ToProtogenBroker() *protogen.Broker {
	return &protogen.Broker{
		Id:       b.ID.String(),
		Name:     b.Name,
		ImageId:  b.ImageID.UUID.String(),
		Disabled: b.Disabled,
	}
}

// FromProtogenBroker converts a protogen.Broker to a Broker
func FromProtogenBroker(b *protogen.Broker) Broker {
	return Broker{
		ID:       uuid.MustParse(b.GetId()),
		Name:     b.GetName(),
		ImageID:  uuid.NullUUID{},
		Disabled: b.GetDisabled(),
	}
}
