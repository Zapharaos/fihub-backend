package models

import (
	"errors"
	"github.com/google/uuid"
)

var (
	errBrokerIdRequired   = errors.New("broker-required")
	errBrokerNameRequired = errors.New("name-required")
)

// Broker represents a broker entity in the system
type Broker struct {
	ID       uuid.UUID     `json:"id" db:"id"`
	Name     string        `json:"name" db:"name"`
	ImageID  uuid.NullUUID `json:"image_id" db:"image_id" swaggertype:"string"`
	Disabled bool          `json:"disabled" db:"disabled"`
}

// IsValid checks if the Broker is valid
func (b Broker) IsValid() (bool, error) {
	if b.Name == "" {
		return false, errBrokerNameRequired
	}
	return true, nil
}
