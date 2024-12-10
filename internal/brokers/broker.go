package brokers

import (
	"github.com/google/uuid"
)

// Broker represents a broker entity in the system
type Broker struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
