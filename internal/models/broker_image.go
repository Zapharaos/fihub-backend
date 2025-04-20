package models

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
)

var (
	errImageNameInvalid  = errors.New("name-invalid")
	errImageDataRequired = errors.New("data-required")
)

const (
	ImageNameMinLength = 3
	ImageNameMaxLength = 100
)

// BrokerImage represents an image associated with a Broker
type BrokerImage struct {
	ID       uuid.UUID `json:"id"`
	BrokerID uuid.UUID `json:"broker_id"`
	Name     string    `json:"name"`
	Data     []byte    `json:"-"`
}

// IsValid checks if the BrokerImage is valid
func (i BrokerImage) IsValid() (bool, error) {
	if i.BrokerID == uuid.Nil {
		return false, errBrokerIdRequired
	}
	if len(i.Name) < ImageNameMinLength {
		return false, errImageNameInvalid
	}
	if len(i.Name) > ImageNameMaxLength {
		return false, errImageNameInvalid
	}
	if len(i.Data) == 0 {
		return false, errImageDataRequired
	}
	return true, nil
}

// ToProtogenBrokerImage converts a BrokerImage to a protogen.BrokerImage
func (i BrokerImage) ToProtogenBrokerImage() *protogen.BrokerImage {
	return &protogen.BrokerImage{
		Id:       i.ID.String(),
		BrokerId: i.BrokerID.String(),
		Name:     i.Name,
		Data:     i.Data,
	}
}

// FromProtogenBrokerImage converts a protogen.BrokerImage to a BrokerImage
func FromProtogenBrokerImage(i *protogen.BrokerImage) BrokerImage {
	return BrokerImage{
		ID:       uuid.MustParse(i.GetId()),
		BrokerID: uuid.MustParse(i.GetBrokerId()),
		Name:     i.GetName(),
		Data:     i.GetData(),
	}
}
