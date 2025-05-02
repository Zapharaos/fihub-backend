package models

import (
	"errors"
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
	ID       uuid.UUID `json:"id" db:"id"`
	BrokerID uuid.UUID `json:"broker_id" db:"broker_id"`
	Name     string    `json:"name" db:"name"`
	Data     []byte    `json:"-" db:"data"`
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
