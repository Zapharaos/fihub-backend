package brokers

import (
	"errors"
	"github.com/google/uuid"
)

var (
	errBrokerNameRequired = errors.New("name-required")
)

// Broker represents a broker entity in the system
type Broker struct {
	ID       uuid.UUID     `json:"id"`
	Name     string        `json:"name"`
	ImageID  uuid.NullUUID `json:"image_id" swaggertype:"string"`
	Disabled bool          `json:"disabled"`
}

// IsValid checks if the broker is valid
func (b Broker) IsValid() (bool, error) {
	if b.Name == "" {
		return false, errBrokerNameRequired
	}
	return true, nil
}
