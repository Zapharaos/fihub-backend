package models

import (
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
)

// BrokerUser represents a user's Broker entity in the system
type BrokerUser struct {
	UserID uuid.UUID `json:"-"`
	Broker Broker    `json:"broker"`
}

// BrokerUserInput represents a user Broker entity received by the system
type BrokerUserInput struct {
	UserID   string `json:"user_id"`
	BrokerID string `json:"broker_id"`
}

// IsValid checks if a BrokerUserInput is valid and has no missing mandatory PGFields
// * BrokerID must be valid
func (u *BrokerUserInput) IsValid() (bool, error) {
	if _, err := uuid.Parse(u.BrokerID); err != nil {
		return false, errBrokerIdRequired
	}
	return true, nil
}

// ToUser Returns a BrokerUser struct from a BrokerUserInput struct
func (u *BrokerUserInput) ToUser() BrokerUser {
	userID, _ := uuid.Parse(u.UserID)
	brokerID, _ := uuid.Parse(u.BrokerID)

	return BrokerUser{
		UserID: userID,
		Broker: Broker{ID: brokerID},
	}
}

// ToProtogenBrokerUser converts a BrokerUser to a protogen.BrokerUser
func (ub *BrokerUser) ToProtogenBrokerUser() *protogen.BrokerUser {
	return &protogen.BrokerUser{
		UserId: ub.UserID.String(),
		Broker: ub.Broker.ToProtogenBroker(),
	}
}

// FromProtogenBrokerUser converts a protogen.BrokerUser to a BrokerUser
func FromProtogenBrokerUser(ub *protogen.BrokerUser) BrokerUser {
	return BrokerUser{
		UserID: uuid.MustParse(ub.GetUserId()),
		Broker: FromProtogenBroker(ub.Broker),
	}
}
