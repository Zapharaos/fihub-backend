package brokers

import (
	"github.com/google/uuid"
)

// BrokerImage represents an image associated with a broker
type BrokerImage struct {
	ID       uuid.UUID `json:"id"`
	BrokerID uuid.UUID `json:"broker_id"`
	Name     string    `json:"name"`
	Data     []byte    `json:"-"`
}

// IsValid checks if the image is valid
func (i BrokerImage) IsValid() (bool, error) {
	if i.BrokerID == uuid.Nil {
		return false, nil
	}
	if len(i.Name) < 3 {
		return false, nil
	}
	if len(i.Name) > 100 {
		return false, nil
	}
	if len(i.Data) == 0 {
		return false, nil
	}
	return true, nil
}