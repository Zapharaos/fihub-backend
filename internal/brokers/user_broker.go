package brokers

import "github.com/google/uuid"

// UserBroker represents a user broker entity in the system
type UserBroker struct {
	UserID   uuid.UUID `json:"user_id"`
	BrokerID uuid.UUID `json:"broker_id"`
}
