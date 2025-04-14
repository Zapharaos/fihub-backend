package transaction

import (
	"context"
	"github.com/Zapharaos/fihub-backend/protogen/transaction"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service is the implementation of the TransactionService interface.
type Service struct {
	transaction.UnimplementedTransactionServiceServer
}

// CreateTransaction implements the CreateTransaction RPC method.
func (s *Service) CreateTransaction(ctx context.Context, req *transaction.CreateTransactionRequest) (*transaction.CreateTransactionResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &transaction.CreateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetBrokerId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetBrokerId()), zap.Error(err))
		return &transaction.CreateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Construct the transaction input object
	transactionInput := TransactionInput{
		UserID:    userID,
		BrokerID:  brokerID,
		Date:      req.GetDate().AsTime(),
		Type:      TransactionType(req.GetTransactionType()),
		Asset:     req.GetAsset(),
		Quantity:  req.GetQuantity(),
		Price:     req.GetPrice(),
		Fee:       req.GetFee(),
		PriceUnit: req.GetPrice() / req.GetQuantity(),
	}

	// Validate the transaction input
	_, validationErr := transactionInput.IsValid()
	if validationErr != nil {
		// Log the validation error and return an invalid response
		zap.L().Error("Transaction validation failed", zap.Error(validationErr))
		return &transaction.CreateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.InvalidArgument, validationErr.Error())
	}

	// Create the transaction
	transactionID, err := R().Create(transactionInput)
	if err != nil {
		zap.L().Error("Create transaction", zap.Error(err))
		return &transaction.CreateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.Internal, "Failed to create transaction")
	}

	// Get transaction back from database
	t, ok, err := R().Get(transactionID)
	if err != nil {
		zap.L().Error("Cannot get transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
		return &transaction.CreateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.Internal, "Failed to get transaction")
	}
	if !ok {
		zap.L().Error("Transaction not found after creation", zap.String("uuid", transactionID.String()))
		return &transaction.CreateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.NotFound, "Transaction not found")
	}

	// Implement the logic to create a transaction
	// This is just a placeholder implementation
	return &transaction.CreateTransactionResponse{
		Transaction: t.ToGenTransaction(),
	}, nil
}

func (s *Service) GetTransaction(ctx context.Context, req *transaction.GetTransactionRequest) (*transaction.GetTransactionResponse, error) {
	// Implement the logic to get a transaction
	// This is just a placeholder implementation
	return &transaction.GetTransactionResponse{
		Transaction: &transaction.Transaction{},
	}, nil
}

func (s *Service) ListTransactions(ctx context.Context, req *transaction.ListTransactionsRequest) (*transaction.ListTransactionsResponse, error) {
	// Implement the logic to list transactions
	// This is just a placeholder implementation
	return &transaction.ListTransactionsResponse{
		Transactions: []*transaction.Transaction{},
	}, nil
}

func (s *Service) UpdateTransaction(ctx context.Context, req *transaction.UpdateTransactionRequest) (*transaction.UpdateTransactionResponse, error) {
	// Implement the logic to update a transaction
	// This is just a placeholder implementation
	return &transaction.UpdateTransactionResponse{
		Transaction: &transaction.Transaction{},
	}, nil
}

func (s *Service) DeleteTransaction(ctx context.Context, req *transaction.DeleteTransactionRequest) (*transaction.DeleteTransactionResponse, error) {
	// Implement the logic to delete a transaction
	// This is just a placeholder implementation
	return &transaction.DeleteTransactionResponse{}, nil
}
