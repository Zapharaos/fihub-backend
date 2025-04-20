package models

import (
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
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

// TestToGenTransactionType tests the ToGenTransactionType method of TransactionType
func TestToGenTransactionType(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    TransactionType
		expected protogen.TransactionType
	}{
		{"BUY to gen", BUY, protogen.TransactionType_BUY},
		{"SELL to gen", SELL, protogen.TransactionType_SELL},
		{"Invalid to gen", TransactionType("INVALID"), protogen.TransactionType_TRANSACTION_TYPE_UNSPECIFIED},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToGenTransactionType()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestToGenTransaction tests the ToGenTransaction method of Transaction
func TestToGenTransaction(t *testing.T) {
	// Create test UUIDs
	id := uuid.New()
	userId := uuid.New()
	brokerId := uuid.New()
	testDate := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	// Define test case
	tr := Transaction{
		ID:        id,
		UserID:    userId,
		Broker:    Broker{ID: brokerId},
		Date:      testDate,
		Type:      BUY,
		Asset:     "asset",
		Quantity:  10.5,
		Price:     150.75,
		PriceUnit: 14.36,
		Fee:       1.99,
	}

	// Convert to gen transaction
	result := tr.ToGenTransaction()

	// Assert results
	assert.Equal(t, id.String(), result.Id)
	assert.Equal(t, userId.String(), result.UserId)
	assert.Equal(t, brokerId.String(), result.BrokerId)
	assert.Equal(t, testDate.Unix(), result.Date.AsTime().Unix())
	assert.Equal(t, protogen.TransactionType_BUY, result.TransactionType)
	assert.Equal(t, "asset", result.Asset)
	assert.Equal(t, 10.5, result.Quantity)
	assert.Equal(t, 150.75, result.Price)
	assert.Equal(t, 14.36, result.PriceUnit)
	assert.Equal(t, 1.99, result.Fee)
}

// TestFromGenTransactionType tests the FromGenTransactionType function
func TestFromGenTransactionType(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    protogen.TransactionType
		expected TransactionType
	}{
		{"BUY from gen", protogen.TransactionType_BUY, BUY},
		{"SELL from gen", protogen.TransactionType_SELL, SELL},
		{"Unspecified from gen", protogen.TransactionType_TRANSACTION_TYPE_UNSPECIFIED, TransactionType("")},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromGenTransactionType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFromGenTransaction tests the FromGenTransaction function
func TestFromGenTransaction(t *testing.T) {
	// Create test UUIDs as strings
	idStr := uuid.New().String()
	userIdStr := uuid.New().String()
	brokerIdStr := uuid.New().String()
	id, _ := uuid.Parse(idStr)
	userId, _ := uuid.Parse(userIdStr)
	brokerId, _ := uuid.Parse(brokerIdStr)

	testDate := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	// Create a gen transaction
	genTransaction := &protogen.Transaction{
		Id:              idStr,
		UserId:          userIdStr,
		BrokerId:        brokerIdStr,
		Date:            timestamppb.New(testDate),
		TransactionType: protogen.TransactionType_SELL,
		Asset:           "TSLA",
		Quantity:        5.25,
		Price:           200.50,
		PriceUnit:       38.19,
		Fee:             2.75,
	}

	// Convert from gen transaction
	result := FromGenTransaction(genTransaction)

	// Assert results
	assert.Equal(t, id, result.ID)
	assert.Equal(t, userId, result.UserID)
	assert.Equal(t, brokerId, result.Broker.ID)
	assert.Equal(t, testDate.Unix(), result.Date.Unix())
	assert.Equal(t, SELL, result.Type)
	assert.Equal(t, "TSLA", result.Asset)
	assert.Equal(t, 5.25, result.Quantity)
	assert.Equal(t, 200.50, result.Price)
	assert.Equal(t, 38.19, result.PriceUnit)
	assert.Equal(t, 2.75, result.Fee)

	// Test with invalid UUIDs
	invalidGenTransaction := &protogen.Transaction{
		Id:       "invalid-uuid",
		UserId:   userIdStr,
		BrokerId: brokerIdStr,
	}

	invalidResult := FromGenTransaction(invalidGenTransaction)
	assert.Equal(t, Transaction{}, invalidResult)
}
