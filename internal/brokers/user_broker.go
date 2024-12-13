package brokers

import (
	"errors"
	"github.com/google/uuid"
)

// UserBroker represents a user broker entity in the system
type UserBroker struct {
	UserID   uuid.UUID `json:"user_id"`
	BrokerID uuid.UUID `json:"broker_id"`
}

// UserBrokerInput represents a user broker entity received by the system
type UserBrokerInput struct {
	UserID   string `json:"user_id"`
	BrokerID string `json:"broker_id"`
}

// IsValid checks if a UserBrokerInput is valid and has no missing mandatory PGFields
// * BrokerID must be valid
func (u *UserBrokerInput) IsValid() (bool, error) {
	if _, err := uuid.Parse(u.BrokerID); err != nil {
		return false, errors.New("broker-required")
	}
	return true, nil
}

// ToUserBroker Returns a UserBroker struct
func (u *UserBrokerInput) ToUserBroker() UserBroker {

	// BrokerID
	brokerID, _ := uuid.Parse(u.BrokerID)

	return UserBroker{
		BrokerID: brokerID,
	}
}
