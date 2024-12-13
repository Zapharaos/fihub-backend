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

// ToUserBroker checks if a UserBrokerInput is valid and converts it into UserBroker type
// * BrokerID must not be of uuid.UUID type
func (ubi *UserBrokerInput) ToUserBroker() (UserBroker, bool, error) {

	// BrokerID : Validate
	brokerID, err := uuid.Parse(ubi.BrokerID)
	if err != nil {
		return UserBroker{}, false, errors.New("broker-required")
	}

	// Valid
	userBroker := UserBroker{
		BrokerID: brokerID,
	}
	return userBroker, true, nil
}
