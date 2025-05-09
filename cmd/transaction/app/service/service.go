package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/transaction/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/transactionpb"
	"github.com/Zapharaos/fihub-backend/internal/mappers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service is the implementation of the TransactionService interface.
type Service struct {
	transactionpb.UnimplementedTransactionServiceServer
}

// CreateTransaction implements the CreateTransaction RPC method.
func (s *Service) CreateTransaction(ctx context.Context, req *transactionpb.CreateTransactionRequest) (*transactionpb.CreateTransactionResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &transactionpb.CreateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetBrokerId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetBrokerId()), zap.Error(err))
		return &transactionpb.CreateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Construct the transaction input object
	transactionInput := models.TransactionInput{
		UserID:    userID,
		BrokerID:  brokerID,
		Date:      req.GetDate().AsTime(),
		Type:      mappers.TransactionTypeFromProto(req.GetTransactionType()),
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
		return &transactionpb.CreateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.InvalidArgument, validationErr.Error())
	}

	// Create the transaction
	transactionID, err := repositories.R().Create(transactionInput)
	if err != nil {
		zap.L().Error("Create transaction", zap.Error(err))
		return &transactionpb.CreateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.Internal, "Failed to create transaction")
	}

	// Get transaction back from database
	t, ok, err := repositories.R().Get(transactionID)
	if err != nil {
		zap.L().Error("Cannot get transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
		return &transactionpb.CreateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.Internal, "Failed to get transaction")
	}
	if !ok {
		zap.L().Error("Transaction not found after creation", zap.String("uuid", transactionID.String()))
		return &transactionpb.CreateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.NotFound, "Transaction not found")
	}

	// Return the created transaction
	return &transactionpb.CreateTransactionResponse{
		Transaction: mappers.TransactionToProto(t),
	}, nil
}

// GetTransaction implements the GetTransaction RPC method.
func (s *Service) GetTransaction(ctx context.Context, req *transactionpb.GetTransactionRequest) (*transactionpb.GetTransactionResponse, error) {
	transactionID, err := uuid.Parse(req.GetTransactionId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid transaction ID", zap.String("transaction_id", req.GetTransactionId()), zap.Error(err))
		return &transactionpb.GetTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.InvalidArgument, "Invalid transaction ID")
	}

	// Get transaction
	t, ok, err := repositories.R().Get(transactionID)
	if err != nil {
		zap.L().Error("Cannot get transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
		return &transactionpb.GetTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.Internal, "Failed to get transaction")
	}
	if !ok {
		zap.L().Error("Transaction not found", zap.String("uuid", transactionID.String()))
		return &transactionpb.GetTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.NotFound, "Transaction not found")
	}

	// Return the transaction
	return &transactionpb.GetTransactionResponse{
		Transaction: mappers.TransactionToProto(t),
	}, nil
}

// ListTransactions implements the ListTransactions RPC method.
func (s *Service) ListTransactions(ctx context.Context, req *transactionpb.ListTransactionsRequest) (*transactionpb.ListTransactionsResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &transactionpb.ListTransactionsResponse{
			Transactions: nil,
		}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Get all transactions
	t, err := repositories.R().GetAll(userID)
	if err != nil {
		zap.L().Error("Cannot get transactions", zap.String("uuid", userID.String()), zap.Error(err))
		return &transactionpb.ListTransactionsResponse{
			Transactions: nil,
		}, status.Error(codes.Internal, "Failed to get transaction")
	}

	// Convert transactions to gRPC format
	return &transactionpb.ListTransactionsResponse{
		Transactions: mappers.TransactionsToProto(t),
	}, nil
}

// UpdateTransaction implements the UpdateTransaction RPC method.
func (s *Service) UpdateTransaction(ctx context.Context, req *transactionpb.UpdateTransactionRequest) (*transactionpb.UpdateTransactionResponse, error) {

	// Parse the transaction ID from the request
	transactionID, err := uuid.Parse(req.GetTransactionId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid transaction ID", zap.String("transaction_id", req.GetUserId()), zap.Error(err))
		return &transactionpb.UpdateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.InvalidArgument, "Invalid transaction ID")
	}

	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &transactionpb.UpdateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetBrokerId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetBrokerId()), zap.Error(err))
		return &transactionpb.UpdateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Construct the transaction input object
	transactionInput := models.TransactionInput{
		ID:        transactionID,
		UserID:    userID,
		BrokerID:  brokerID,
		Date:      req.GetDate().AsTime(),
		Type:      mappers.TransactionTypeFromProto(req.GetTransactionType()),
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
		return &transactionpb.UpdateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.InvalidArgument, validationErr.Error())
	}

	// Verify that the transaction belongs to the user
	oldTransaction, ok, err := repositories.R().Get(transactionID)
	if err != nil {
		zap.L().Error("Cannot get transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
		return &transactionpb.UpdateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.Internal, "Failed to get transaction")
	}
	if !ok {
		zap.L().Error("Transaction not found after creation", zap.String("uuid", transactionID.String()))
		return &transactionpb.UpdateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.NotFound, "Transaction not found")
	}
	if oldTransaction.UserID != userID {
		zap.L().Warn("Transaction does not belong to user", zap.String("uuid", transactionID.String()))
		return &transactionpb.UpdateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.PermissionDenied, "Transaction does not belong to user")
	}

	// Update the transaction
	err = repositories.R().Update(transactionInput)
	if err != nil {
		zap.L().Error("Cannot update transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
		return &transactionpb.UpdateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.Internal, "Failed to update transaction")
	}

	// Get transaction back from database
	t, ok, err := repositories.R().Get(transactionID)
	if err != nil {
		zap.L().Error("Cannot get transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
		return &transactionpb.UpdateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.Internal, "Failed to get transaction")
	}
	if !ok {
		zap.L().Error("Transaction not found after update", zap.String("uuid", transactionID.String()))
		return &transactionpb.UpdateTransactionResponse{
			Transaction: nil,
		}, status.Error(codes.NotFound, "Transaction not found")
	}

	// Return the created transaction
	return &transactionpb.UpdateTransactionResponse{
		Transaction: mappers.TransactionToProto(t),
	}, nil
}

// DeleteTransaction implements the DeleteTransaction RPC method.
func (s *Service) DeleteTransaction(ctx context.Context, req *transactionpb.DeleteTransactionRequest) (*transactionpb.DeleteTransactionResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &transactionpb.DeleteTransactionResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse the transaction ID from the request
	transactionID, err := uuid.Parse(req.GetTransactionId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid transaction ID", zap.String("transaction_id", req.GetTransactionId()), zap.Error(err))
		return &transactionpb.DeleteTransactionResponse{}, status.Error(codes.InvalidArgument, "Invalid transaction ID")
	}

	// Verify that the transaction belongs to the user
	t, ok, err := repositories.R().Get(transactionID)
	if err != nil {
		zap.L().Error("Cannot get transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
		return &transactionpb.DeleteTransactionResponse{}, status.Error(codes.Internal, "Failed to get transaction")
	}
	if !ok {
		zap.L().Error("Transaction not found", zap.String("uuid", transactionID.String()))
		return &transactionpb.DeleteTransactionResponse{}, status.Error(codes.NotFound, "Transaction not found")
	}
	if t.UserID != userID {
		zap.L().Warn("Transaction does not belong to user", zap.String("uuid", transactionID.String()))
		return &transactionpb.DeleteTransactionResponse{}, status.Error(codes.PermissionDenied, "Transaction does not belong to user")
	}

	// Remove transaction
	err = repositories.R().Delete(models.Transaction{ID: transactionID, UserID: userID})
	if err != nil {
		zap.L().Error("Cannot remove transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
		return &transactionpb.DeleteTransactionResponse{}, status.Error(codes.Internal, "Failed to remove transaction")
	}

	// Return success response
	return &transactionpb.DeleteTransactionResponse{}, nil
}

// DeleteTransactionByBroker implements the DeleteTransactionByBroker RPC method.
func (s *Service) DeleteTransactionByBroker(ctx context.Context, req *transactionpb.DeleteTransactionByBrokerRequest) (*transactionpb.DeleteTransactionByBrokerResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &transactionpb.DeleteTransactionByBrokerResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetBrokerId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetBrokerId()), zap.Error(err))
		return &transactionpb.DeleteTransactionByBrokerResponse{}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Remove transactions by broker
	err = repositories.R().DeleteByBroker(models.Transaction{
		UserID: userID,
		Broker: models.Broker{
			ID: brokerID,
		},
	})
	if err != nil {
		zap.L().Error("Cannot remove broker related transactions", zap.String("uuid", brokerID.String()), zap.Error(err))
		return &transactionpb.DeleteTransactionByBrokerResponse{}, status.Error(codes.Internal, "Failed to remove broker related transactions")
	}

	// Return success response
	return &transactionpb.DeleteTransactionByBrokerResponse{}, nil
}
