package transactions

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestTransactionTypeIsValid tests the IsValid method of TransactionType
func TestTransactionTypeIsValid(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string          // Test case name
		input    TransactionType // TransactionType instance to test
		expected bool            // Expected result
	}{
		{"Valid BUY", BUY, true},
		{"Valid SELL", SELL, true},
		{"Invalid Type", TransactionType("INVALID"), false},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, _ := tt.input.IsValid()
			assert.Equal(t, tt.expected, valid)
		})
	}
}

// TestTransactionInputIsValid tests the IsValid method of TransactionInput
func TestTransactionInputIsValid(t *testing.T) {
	// Define valid values
	validUUID := uuid.New()
	validDate := time.Now().Add(-time.Hour) // 1 hour in the past
	validTransactionType := BUY
	validAsset := "AAPL"
	validQuantity := 10.0
	validPrice := 150.0
	validFee := 1.0

	// Define test cases
	tests := []struct {
		name     string           // Test case name
		input    TransactionInput // TransactionInput instance to test
		expected bool             // Expected result
		error    error            // Expected error
	}{
		{
			"Valid TransactionInput",
			TransactionInput{
				ID:       validUUID,
				UserID:   validUUID,
				BrokerID: validUUID,
				Date:     validDate,
				Type:     validTransactionType,
				Asset:    validAsset,
				Quantity: validQuantity,
				Price:    validPrice,
				Fee:      validFee,
			},
			true,
			nil,
		},
		{
			"Invalid BrokerID",
			TransactionInput{
				ID:       validUUID,
				UserID:   validUUID,
				BrokerID: uuid.Nil,
				Date:     validDate,
				Type:     validTransactionType,
				Asset:    validAsset,
				Quantity: validQuantity,
				Price:    validPrice,
				Fee:      validFee,
			},
			false,
			errBrokerRequired,
		},
		{
			"Empty Date",
			TransactionInput{
				ID:       validUUID,
				UserID:   validUUID,
				BrokerID: validUUID,
				Date:     time.Time{},
				Type:     validTransactionType,
				Asset:    validAsset,
				Quantity: validQuantity,
				Price:    validPrice,
				Fee:      validFee,
			},
			false,
			errDateRequired,
		},
		{
			"Future Date",
			TransactionInput{
				ID:       validUUID,
				UserID:   validUUID,
				BrokerID: validUUID,
				Date:     time.Now().Add(time.Hour), // 1 hour in the future
				Type:     validTransactionType,
				Asset:    validAsset,
				Quantity: validQuantity,
				Price:    validPrice,
				Fee:      validFee,
			},
			false,
			errDateFuture,
		},
		{
			"Invalid Type",
			TransactionInput{
				ID:       validUUID,
				UserID:   validUUID,
				BrokerID: validUUID,
				Date:     validDate,
				Type:     TransactionType("INVALID"),
				Asset:    validAsset,
				Quantity: validQuantity,
				Price:    validPrice,
				Fee:      validFee,
			},
			false,
			errTypeInvalid,
		},
		{
			"Empty Asset",
			TransactionInput{
				ID:       validUUID,
				UserID:   validUUID,
				BrokerID: validUUID,
				Date:     validDate,
				Type:     validTransactionType,
				Asset:    "",
				Quantity: validQuantity,
				Price:    validPrice,
				Fee:      validFee,
			},
			false,
			errAssetRequired,
		},
		{
			"Negative Quantity",
			TransactionInput{
				ID:       validUUID,
				UserID:   validUUID,
				BrokerID: validUUID,
				Date:     validDate,
				Type:     validTransactionType,
				Asset:    validAsset,
				Quantity: -10,
				Price:    validPrice,
				Fee:      validFee,
			},
			false,
			errQuantityInvalid,
		},
		{
			"Negative Price",
			TransactionInput{
				ID:       validUUID,
				UserID:   validUUID,
				BrokerID: validUUID,
				Date:     validDate,
				Type:     validTransactionType,
				Asset:    validAsset,
				Quantity: validQuantity,
				Price:    -150,
				Fee:      validFee,
			},
			false,
			errPriceInvalid,
		},
		{
			"Negative Fee",
			TransactionInput{
				ID:       validUUID,
				UserID:   validUUID,
				BrokerID: validUUID,
				Date:     validDate,
				Type:     validTransactionType,
				Asset:    validAsset,
				Quantity: validQuantity,
				Price:    validPrice,
				Fee:      -1,
			},
			false,
			errFeeInvalid,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, _ := tt.input.IsValid()
			assert.Equal(t, tt.expected, valid)
		})
	}
}
