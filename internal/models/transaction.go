package models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type TransactionType string

// Declare constants of type TransactionType
const (
	BUY  TransactionType = "BUY"
	SELL TransactionType = "SELL"
)

var (
	errBrokerRequired  = errors.New("broker-required")
	errDateRequired    = errors.New("date-required")
	errDateFuture      = errors.New("date-future")
	errTypeInvalid     = errors.New("type-invalid")
	errAssetRequired   = errors.New("asset-required")
	errQuantityInvalid = errors.New("quantity-invalid")
	errPriceInvalid    = errors.New("price-invalid")
	errFeeInvalid      = errors.New("fee-invalid")
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
	Broker    Broker          `json:"broker"`
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
	if t == BUY || t == SELL {
		return true, nil
	}
	return false, errTypeInvalid
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
		return false, errBrokerRequired
	}

	// Date
	if t.Date.IsZero() {
		return false, errDateRequired
	}
	if !t.Date.Before(time.Now()) {
		return false, errDateFuture
	}

	// Transaction Type
	if ok, err := t.Type.IsValid(); !ok {
		return false, err
	}

	// Asset
	if t.Asset == "" {
		return false, errAssetRequired
	}

	// Quantity
	if t.Quantity <= 0 {
		return false, errQuantityInvalid
	}

	// Price
	if t.Price <= 0 {
		return false, errPriceInvalid
	}

	// Fee
	if t.Fee < 0 {
		return false, errFeeInvalid
	}

	return true, nil
}
