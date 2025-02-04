package brokers

import (
	"errors"
	"github.com/google/uuid"
)

var (
	errImageNameInvalid  = errors.New("name-invalid")
	errImageDataRequired = errors.New("data-required")
)

// Image represents an image associated with a Broker
type Image struct {
	ID       uuid.UUID `json:"id"`
	BrokerID uuid.UUID `json:"broker_id"`
	Name     string    `json:"name"`
	Data     []byte    `json:"-"`
}

const ImageNameMinLength = 3
const ImageNameMaxLength = 100

// IsValid checks if the Image is valid
func (i Image) IsValid() (bool, error) {
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
