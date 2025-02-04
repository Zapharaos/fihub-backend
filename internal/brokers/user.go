package brokers

import (
	"github.com/google/uuid"
)

// User represents a user's Broker entity in the system
type User struct {
	UserID uuid.UUID `json:"-"`
	Broker Broker    `json:"broker"`
}

// UserInput represents a user Broker entity received by the system
type UserInput struct {
	UserID   string `json:"user_id"`
	BrokerID string `json:"broker_id"`
}

// IsValid checks if a UserInput is valid and has no missing mandatory PGFields
// * BrokerID must be valid
func (u *UserInput) IsValid() (bool, error) {
	if _, err := uuid.Parse(u.BrokerID); err != nil {
		return false, errBrokerIdRequired
	}
	return true, nil
}

// ToUser Returns a User struct from a UserInput struct
func (u *UserInput) ToUser() User {

	// BrokerID
	brokerID, _ := uuid.Parse(u.BrokerID)

	return User{
		Broker: Broker{ID: brokerID},
	}
}
