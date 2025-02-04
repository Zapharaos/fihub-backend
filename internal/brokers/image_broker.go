package brokers

import (
	"errors"
	"github.com/google/uuid"
)

var (
	errBrokerIdRequired  = errors.New("broker-required")
	errImageNameInvalid  = errors.New("name-invalid")
	errImageDataRequired = errors.New("data-required")
)

// BrokerImage represents an image associated with a broker
type BrokerImage struct {
	ID       uuid.UUID `json:"id"`
	BrokerID uuid.UUID `json:"broker_id"`
	Name     string    `json:"name"`
	Data     []byte    `json:"-"`
}

const ImageNameMinLength = 3
const ImageNameMaxLength = 100

// IsValid checks if the image is valid
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
