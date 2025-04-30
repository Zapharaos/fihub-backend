package mappers

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TransactionTypeToProto converts a models.TransactionType to a protogen.TransactionType
func TransactionTypeToProto(t models.TransactionType) protogen.TransactionType {
	switch t {
	case models.BUY:
		return protogen.TransactionType_BUY
	case models.SELL:
		return protogen.TransactionType_SELL
	default:
		return protogen.TransactionType_TRANSACTION_TYPE_UNSPECIFIED
	}
}

// TransactionTypeFromProto converts a protogen.TransactionType to a models.TransactionType
func TransactionTypeFromProto(t protogen.TransactionType) models.TransactionType {
	switch t {
	case protogen.TransactionType_BUY:
		return models.BUY
	case protogen.TransactionType_SELL:
		return models.SELL
	default:
		return ""
	}
}

// TransactionToProto converts a models.Transaction to a protogen.Transaction
func TransactionToProto(t models.Transaction) *protogen.Transaction {
	return &protogen.Transaction{
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

// TransactionFromProto converts a protogen.Transaction to a models.Transaction
func TransactionFromProto(t *protogen.Transaction) models.Transaction {
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

// TransactionsToProto converts a slice of models.Transaction to a slice of protogen.Transaction
func TransactionsToProto(transactions []models.Transaction) []*protogen.Transaction {
	protoTransactions := make([]*protogen.Transaction, len(transactions))
	for i, transaction := range transactions {
		protoTransactions[i] = TransactionToProto(transaction)
	}
	return protoTransactions
}

// TransactionsFromProto converts a slice of protogen.Transaction to a slice of models.Transaction
func TransactionsFromProto(transactions []*protogen.Transaction) []models.Transaction {
	protoTransactions := make([]models.Transaction, len(transactions))
	for i, transaction := range transactions {
		protoTransactions[i] = TransactionFromProto(transaction)
	}
	return protoTransactions
}
