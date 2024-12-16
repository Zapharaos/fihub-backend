package brokers

import (
	"github.com/google/uuid"
)

// Broker represents a broker entity in the system
type Broker struct {
	ID      uuid.UUID  `json:"id"`
	Name    string     `json:"name"`
	ImageID *uuid.UUID `json:"image_id"`
}

// IsValid checks if the broker is valid
func (b Broker) IsValid() (bool, error) {
	if b.Name == "" {
		return false, nil
	}
	return true, nil
}
