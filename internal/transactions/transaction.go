package transactions

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/google/uuid"
	"strings"
	"time"
)

type TransactionType string

// Declare constants of type TransactionType
const (
	BUY  TransactionType = "BUY"
	SELL TransactionType = "SELL"
)

// TransactionInput represents a transaction entity in the system
type TransactionInput struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	BrokerID  uuid.UUID       `json:"broker_id"`
	Date      time.Time       `json:"date"`
	Type      TransactionType `json:"transaction_type"`
	Asset     string          `json:"asset"`
	Quantity  float64         `json:"quantity"`
	Price     float64         `json:"price"`
	PriceUnit float64         `json:"price_unit"`
	Fee       float64         `json:"fee"`
}

// Transaction represents a transaction entity in the system
type Transaction struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Broker    brokers.Broker  `json:"broker"`
	Date      time.Time       `json:"date"`
	Type      TransactionType `json:"transaction_type"`
	Asset     string          `json:"asset"`
	Quantity  float64         `json:"quantity"`
	Price     float64         `json:"price"`
	PriceUnit float64         `json:"price_unit"`
	Fee       float64         `json:"fee"`
}

// IsValid checks if a TransactionType is valid and
func (t TransactionType) IsValid() (bool, error) {

	// Check if the TransactionType is valid (through uppercase string comparison)
	switch TransactionType(strings.ToUpper(string(t))) {
	case BUY:
	case SELL:
		return true, nil
	}

	return false, errors.New("type-invalid")
}

// IsValid checks if a TransactionInput is valid and has no missing mandatory PGFields
// * BrokerID must not be empty
// * Date must not be empty
// * Date must not be in the future
// * Type must be valid (see TransactionType)
// * Asset must not be empty
// * Quantity must be positive
// * Price must be positive
// * PriceUnit must be positive
// * Fee must not be negative
func (t *TransactionInput) IsValid() (bool, error) {
	// Broker
	if t.BrokerID == uuid.Nil {
		return false, errors.New("broker-required")
	}

	// Date
	if t.Date.IsZero() {
		return false, errors.New("date-required")
	}
	if !t.Date.Before(time.Now()) {
		return false, errors.New("date-future")
	}

	// Transaction Type
	if ok, err := t.Type.IsValid(); !ok {
		return false, err
	}

	// Asset
	if t.Asset == "" {
		return false, errors.New("asset-required")
	}

	// Quantity
	if t.Quantity <= 0 {
		return false, errors.New("quantity-invalid")
	}

	// Price
	if t.Price <= 0 {
		return false, errors.New("price-invalid")
	}

	// PriceUnit
	if t.PriceUnit <= 0 {
		return false, errors.New("price-invalid")
	}

	// Fee
	if t.Fee < 0 {
		return false, errors.New("fee-invalid")
	}

	return true, nil
}
