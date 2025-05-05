package mappers

import (
	"github.com/Zapharaos/fihub-backend/gen/go/transactionpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

// TestToGenTransactionType tests the TransactionTypeToProto method of TransactionType
func Test_TransactionTypeToProto(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    models.TransactionType
		expected transactionpb.TransactionType
	}{
		{"BUY to gen", models.BUY, transactionpb.TransactionType_BUY},
		{"SELL to gen", models.SELL, transactionpb.TransactionType_SELL},
		{"Invalid to gen", models.TransactionType("INVALID"), transactionpb.TransactionType_TRANSACTION_TYPE_UNSPECIFIED},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TransactionTypeToProto(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test_TransactionTypeFromProto tests the TransactionTypeFromProto function
func Test_TransactionTypeFromProto(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    transactionpb.TransactionType
		expected models.TransactionType
	}{
		{"BUY from gen", transactionpb.TransactionType_BUY, models.BUY},
		{"SELL from gen", transactionpb.TransactionType_SELL, models.SELL},
		{"Unspecified from gen", transactionpb.TransactionType_TRANSACTION_TYPE_UNSPECIFIED, models.TransactionType("")},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TransactionTypeFromProto(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test_TransactionToProto tests the TransactionToProto method
func Test_TransactionToProto(t *testing.T) {
	// Create test UUIDs
	id := uuid.New()
	userId := uuid.New()
	brokerId := uuid.New()
	testDate := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	// Define test case
	tr := models.Transaction{
		ID:        id,
		UserID:    userId,
		Broker:    models.Broker{ID: brokerId},
		Date:      testDate,
		Type:      models.BUY,
		Asset:     "asset",
		Quantity:  10.5,
		Price:     150.75,
		PriceUnit: 14.36,
		Fee:       1.99,
	}

	// Convert to gen transaction
	result := TransactionToProto(tr)

	// Assert results
	assert.Equal(t, id.String(), result.Id)
	assert.Equal(t, userId.String(), result.UserId)
	assert.Equal(t, brokerId.String(), result.BrokerId)
	assert.Equal(t, testDate.Unix(), result.Date.AsTime().Unix())
	assert.Equal(t, transactionpb.TransactionType_BUY, result.TransactionType)
	assert.Equal(t, "asset", result.Asset)
	assert.Equal(t, 10.5, result.Quantity)
	assert.Equal(t, 150.75, result.Price)
	assert.Equal(t, 14.36, result.PriceUnit)
	assert.Equal(t, 1.99, result.Fee)
}

// Test_TransactionFromProto tests the TransactionFromProto function
func Test_TransactionFromProto(t *testing.T) {
	// Create test UUIDs as strings
	idStr := uuid.New().String()
	userIdStr := uuid.New().String()
	brokerIdStr := uuid.New().String()
	id, _ := uuid.Parse(idStr)
	userId, _ := uuid.Parse(userIdStr)
	brokerId, _ := uuid.Parse(brokerIdStr)

	testDate := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	// Create a gen transaction
	genTransaction := &transactionpb.Transaction{
		Id:              idStr,
		UserId:          userIdStr,
		BrokerId:        brokerIdStr,
		Date:            timestamppb.New(testDate),
		TransactionType: transactionpb.TransactionType_SELL,
		Asset:           "TSLA",
		Quantity:        5.25,
		Price:           200.50,
		PriceUnit:       38.19,
		Fee:             2.75,
	}

	// Convert from gen transaction
	result := TransactionFromProto(genTransaction)

	// Assert results
	assert.Equal(t, id, result.ID)
	assert.Equal(t, userId, result.UserID)
	assert.Equal(t, brokerId, result.Broker.ID)
	assert.Equal(t, testDate.Unix(), result.Date.Unix())
	assert.Equal(t, models.SELL, result.Type)
	assert.Equal(t, "TSLA", result.Asset)
	assert.Equal(t, 5.25, result.Quantity)
	assert.Equal(t, 200.50, result.Price)
	assert.Equal(t, 38.19, result.PriceUnit)
	assert.Equal(t, 2.75, result.Fee)
}

// Test_TransactionsToProto tests the TransactionsToProto function
func Test_TransactionsToProto(t *testing.T) {
	// Create test UUIDs
	id := uuid.New()
	userId := uuid.New()
	brokerId := uuid.New()
	testDate := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	// Create a slice of transactions
	transactions := []models.Transaction{
		{
			ID:        id,
			UserID:    userId,
			Broker:    models.Broker{ID: brokerId},
			Date:      testDate,
			Type:      models.BUY,
			Asset:     "asset1",
			Quantity:  10.5,
			Price:     150.75,
			PriceUnit: 14.36,
			Fee:       1.99,
		},
	}

	// Convert to gen transactions
	result := TransactionsToProto(transactions)

	// Assert results
	assert.Equal(t, id.String(), result[0].Id)
	assert.Equal(t, userId.String(), result[0].UserId)
	assert.Equal(t, brokerId.String(), result[0].BrokerId)
	assert.Equal(t, testDate.Unix(), result[0].Date.AsTime().Unix())
	assert.Equal(t, transactionpb.TransactionType_BUY, result[0].TransactionType)
	assert.Equal(t, "asset1", result[0].Asset)
	assert.Equal(t, 10.5, result[0].Quantity)
	assert.Equal(t, 150.75, result[0].Price)
	assert.Equal(t, 14.36, result[0].PriceUnit)
	assert.Equal(t, 1.99, result[0].Fee)
}

// Test_TransactionsFromProto tests the TransactionsFromProto function
func Test_TransactionsFromProto(t *testing.T) {
	// Create test UUIDs
	id := uuid.New()
	userId := uuid.New()
	brokerId := uuid.New()
	testDate := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	// Create a slice of gen transactions
	genTransactions := []*transactionpb.Transaction{
		{
			Id:              id.String(),
			UserId:          userId.String(),
			BrokerId:        brokerId.String(),
			Date:            timestamppb.New(testDate),
			TransactionType: transactionpb.TransactionType_BUY,
			Asset:           "asset1",
			Quantity:        10.5,
			Price:           150.75,
			PriceUnit:       14.36,
			Fee:             1.99,
		},
	}

	// Convert from gen transactions
	result := TransactionsFromProto(genTransactions)

	// Assert results
	assert.Equal(t, id, result[0].ID)
	assert.Equal(t, userId, result[0].UserID)
	assert.Equal(t, brokerId, result[0].Broker.ID)
	assert.Equal(t, testDate.Unix(), result[0].Date.Unix())
	assert.Equal(t, models.BUY, result[0].Type)
	assert.Equal(t, "asset1", result[0].Asset)
	assert.Equal(t, 10.5, result[0].Quantity)
	assert.Equal(t, 150.75, result[0].Price)
	assert.Equal(t, 14.36, result[0].PriceUnit)
	assert.Equal(t, 1.99, result[0].Fee)
}
