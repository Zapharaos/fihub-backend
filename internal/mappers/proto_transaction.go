package mappers

import (
	"github.com/Zapharaos/fihub-backend/gen/go/transactionpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TransactionTypeToProto converts a models.TransactionType to a transactionpb.TransactionType
func TransactionTypeToProto(t models.TransactionType) transactionpb.TransactionType {
	switch t {
	case models.BUY:
		return transactionpb.TransactionType_BUY
	case models.SELL:
		return transactionpb.TransactionType_SELL
	default:
		return transactionpb.TransactionType_TRANSACTION_TYPE_UNSPECIFIED
	}
}

// TransactionTypeFromProto converts a transactionpb.TransactionType to a models.TransactionType
func TransactionTypeFromProto(t transactionpb.TransactionType) models.TransactionType {
	switch t {
	case transactionpb.TransactionType_BUY:
		return models.BUY
	case transactionpb.TransactionType_SELL:
		return models.SELL
	default:
		return ""
	}
}

// TransactionToProto converts a models.Transaction to a transactionpb.Transaction
func TransactionToProto(t models.Transaction) *transactionpb.Transaction {
	return &transactionpb.Transaction{
		Id:              t.ID.String(),
		UserId:          t.UserID.String(),
		BrokerId:        t.Broker.ID.String(),
		Date:            timestamppb.New(t.Date),
		TransactionType: TransactionTypeToProto(t.Type),
		Asset:           t.Asset,
		Quantity:        t.Quantity,
		Price:           t.Price,
		PriceUnit:       t.PriceUnit,
		Fee:             t.Fee,
	}
}

// TransactionFromProto converts a transactionpb.Transaction to a models.Transaction
func TransactionFromProto(t *transactionpb.Transaction) models.Transaction {
	return models.Transaction{
		ID:     uuid.MustParse(t.GetId()),
		UserID: uuid.MustParse(t.GetUserId()),
		Broker: models.Broker{
			ID: uuid.MustParse(t.GetBrokerId()),
		},
		Date:      t.GetDate().AsTime(),
		Type:      TransactionTypeFromProto(t.GetTransactionType()),
		Asset:     t.GetAsset(),
		Quantity:  t.GetQuantity(),
		Price:     t.GetPrice(),
		PriceUnit: t.GetPriceUnit(),
		Fee:       t.GetFee(),
	}
}

// TransactionsToProto converts a slice of models.Transaction to a slice of transactionpb.Transaction
func TransactionsToProto(transactions []models.Transaction) []*transactionpb.Transaction {
	protoTransactions := make([]*transactionpb.Transaction, len(transactions))
	for i, transaction := range transactions {
		protoTransactions[i] = TransactionToProto(transaction)
	}
	return protoTransactions
}

// TransactionsFromProto converts a slice of transactionpb.Transaction to a slice of models.Transaction
func TransactionsFromProto(transactions []*transactionpb.Transaction) []models.Transaction {
	protoTransactions := make([]models.Transaction, len(transactions))
	for i, transaction := range transactions {
		protoTransactions[i] = TransactionFromProto(transaction)
	}
	return protoTransactions
}
