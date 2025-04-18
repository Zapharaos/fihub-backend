package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	gentransaction "github.com/Zapharaos/fihub-backend/protogen/transaction"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"time"
)

// CreateTransaction 	godoc
//
// @Id 					CreateTransaction
//
// @Summary 			Create a new transaction
// @Description 		Create a new transaction.
// @Tags 				Transactions
// @Accept 				json
// @Produce 			json
// @Param 				transaction body 	transactions.TransactionInput true 	"transaction (json)"
// @Security 			Bearer
// @Success 			200 {object} 		transactions.Transaction 		"transaction"
// @Failure 			400 {object} 		render.ErrorResponse 			"Bad PasswordRequest"
// @Failure 			401 {string} 		string 							"Permission denied"
// @Failure 			500 {object} 		render.ErrorResponse 			"Internal Server Error"
// @Router /api/v1/transactions [post]
func CreateTransaction(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user from the context
	user, ok := U().GetUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse request body to TransactionInput
	var transactionInput models.TransactionInput
	err := json.NewDecoder(r.Body).Decode(&transactionInput)
	if err != nil {
		zap.L().Warn("Transaction json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if userBroker exists
	exists, err := brokers.R().U().Exists(models.BrokerUser{UserID: user.ID, Broker: models.Broker{ID: transactionInput.BrokerID}})
	if err != nil {
		zap.L().Error("Check userBroker exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		zap.L().Warn("BrokerUser not found")
		render.BadRequest(w, r, fmt.Errorf("broker-invalid"))
		return
	}

	// Map TransactionInput to gRPC ValidateTransactionRequest
	transactionRequest := &gentransaction.CreateTransactionRequest{
		UserId:          user.ID.String(),
		BrokerId:        transactionInput.BrokerID.String(),
		Date:            timestamppb.New(transactionInput.Date),
		TransactionType: transactionInput.Type.ToGenTransactionType(),
		Asset:           transactionInput.Asset,
		Quantity:        transactionInput.Quantity,
		Price:           transactionInput.Price,
		Fee:             transactionInput.Fee,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create the transaction
	response, err := clients.C().Transaction().CreateTransaction(ctx, transactionRequest)
	if err != nil {
		zap.L().Error("Create transaction", zap.Error(err))
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.InvalidArgument:
				render.BadRequest(w, r, err)
				return
			case codes.NotFound:
				w.WriteHeader(http.StatusNotFound)
				return
			case codes.Internal:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Map gRPC response to Transaction
	render.JSON(w, r, models.FromGenTransaction(response.Transaction))
}

// GetTransaction godoc
//
// @Id 				GetTransaction
//
// @Summary 		Get a transaction
// @Description 	Get a transaction.
// @Tags 			Transactions
// @Accept 			json
// @Produce 		json
// @Param 			id path 		string true 				"transaction id"
// @Security 		Bearer
// @Success 		200 {object} 	transactions.Transaction 	"transaction"
// @Failure 		400 {object} 	render.ErrorResponse 		"Bad PasswordRequest"
// @Failure 		401 {string} 	string 						"Permission denied"
// @Failure 		404 {object} 	render.ErrorResponse 		"Not Found"
// @Failure 		500 {object} 	render.ErrorResponse 		"Internal Server Error"
// @Router /api/v1/transactions/{id} [get]
func GetTransaction(w http.ResponseWriter, r *http.Request) {

	// Retrieve transactionID
	transactionID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve the transaction
	response, err := clients.C().Transaction().GetTransaction(ctx, &gentransaction.GetTransactionRequest{
		TransactionId: transactionID.String(),
	})
	if err != nil {
		zap.L().Error("Get transaction", zap.Error(err))
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.NotFound:
				w.WriteHeader(http.StatusNotFound)
				return
			case codes.Internal:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Map gRPC response to Transaction
	t := models.FromGenTransaction(response.Transaction)

	// Verify that the transaction belongs to the user
	user, ok := U().GetUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if t.UserID != user.ID {
		zap.L().Warn("Transaction does not belong to user", zap.String("uuid", transactionID.String()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	render.JSON(w, r, t)
}

// UpdateTransaction godoc
//
// @Id 				UpdateTransaction
//
// @Summary 		Update a transaction
// @Description 	Update a transaction.
// @Tags 			Transactions
// @Accept 			json
// @Produce 		json
// @Param 			id path 		string true 				"transaction ID"
// @Param 			transaction body transactions.Transaction true "transaction (json)"
// @Security 		Bearer
// @Success 		200 {object} 	transactions.Transaction 	"transaction"
// @Failure 		400 {object} 	render.ErrorResponse 		"Bad PasswordRequest"
// @Failure 		401 {string} 	string 						"Permission denied"
// @Failure			404	{object}	render.ErrorResponse		"Not Found"
// @Failure 		500 {object} 	render.ErrorResponse 		"Internal Server Error"
// @Router /api/v1/transactions/{id} [put]
func UpdateTransaction(w http.ResponseWriter, r *http.Request) {

	// Retrieve transactionID
	transactionID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Get the authenticated user from the context
	user, ok := U().GetUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse request body
	var transactionInput models.TransactionInput
	err := json.NewDecoder(r.Body).Decode(&transactionInput)
	if err != nil {
		zap.L().Warn("Transaction json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate the transaction
	if valid, err := transactionInput.IsValid(); !valid {
		zap.L().Warn("Transaction is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Verify that the userBroker exists
	exists, err := brokers.R().U().Exists(models.BrokerUser{UserID: user.ID, Broker: models.Broker{ID: transactionInput.BrokerID}})
	if err != nil {
		zap.L().Error("Check userBroker exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		zap.L().Warn("BrokerUser not found")
		render.BadRequest(w, r, fmt.Errorf("broker-invalid"))
		return
	}

	// Map TransactionInput to gRPC ValidateTransactionRequest
	transactionRequest := &gentransaction.UpdateTransactionRequest{
		TransactionId:   transactionID.String(),
		UserId:          user.ID.String(),
		BrokerId:        transactionInput.BrokerID.String(),
		Date:            timestamppb.New(transactionInput.Date),
		TransactionType: transactionInput.Type.ToGenTransactionType(),
		Asset:           transactionInput.Asset,
		Quantity:        transactionInput.Quantity,
		Price:           transactionInput.Price,
		Fee:             transactionInput.Fee,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create the transaction
	response, err := clients.C().Transaction().UpdateTransaction(ctx, transactionRequest)
	if err != nil {
		zap.L().Error("Update transaction", zap.Error(err))
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.InvalidArgument:
				render.BadRequest(w, r, err)
				return
			case codes.NotFound:
				w.WriteHeader(http.StatusNotFound)
				return
			case codes.Internal:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Map gRPC response to Transaction
	render.JSON(w, r, models.FromGenTransaction(response.Transaction))
}

// DeleteTransaction godoc
//
// @Id 				DeleteTransaction
//
// @Summary 		Delete a transaction
// @Description 	Delete a transaction.
// @Tags 			Transactions
// @Accept 			json
// @Produce 		json
// @Param 			id path 		string true 			"transaction ID"
// @Security 		Bearer
// @Success 		200 {array} 	string 					"Status OK"
// @Failure 		400 {object} 	render.ErrorResponse 	"Bad PasswordRequest"
// @Failure 		401 {string} 	string 					"Permission denied"
// @Failure			404	{object}	render.ErrorResponse	"Not Found"
// @Failure 		500 {object} 	render.ErrorResponse 	"Internal Server Error"
// @Router /api/v1/transactions/{id} [delete]
func DeleteTransaction(w http.ResponseWriter, r *http.Request) {

	// Retrieve transactionID
	transactionID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Get the authenticated user from the context
	user, ok := U().GetUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve the transaction
	_, err := clients.C().Transaction().DeleteTransaction(ctx, &gentransaction.DeleteTransactionRequest{
		TransactionId: transactionID.String(),
		UserId:        user.ID.String(),
	})
	if err != nil {
		zap.L().Error("Delete transaction", zap.Error(err))
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.NotFound:
				w.WriteHeader(http.StatusNotFound)
				return
			case codes.PermissionDenied:
				w.WriteHeader(http.StatusUnauthorized)
				return
			case codes.Internal:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.OK(w, r)
}

// GetTransactions godoc
//
// @Id 				GetTransactions
//
// @Summary 		Get all transactions
// @Description 	Gets a list of all transactions.
// @Tags 			Transactions
// @Produce 		json
// @Security 		Bearer
// @Success 		200 {array} 	transactions.Transaction 	"List of transactions"
// @Failure 		401 {string} 	string 						"Permission denied"
// @Failure 		500 {object} 	render.ErrorResponse 		"Internal Server Error"
// @Router /api/v1/transactions [get]
func GetTransactions(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user from the context
	user, ok := U().GetUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve the transaction
	response, err := clients.C().Transaction().ListTransactions(ctx, &gentransaction.ListTransactionsRequest{
		UserId: user.ID.String(),
	})
	if err != nil {
		zap.L().Error("List transactions", zap.Error(err))
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.NotFound:
				w.WriteHeader(http.StatusNotFound)
				return
			case codes.Internal:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Map gRPC response to Transaction array
	t := make([]models.Transaction, len(response.Transactions))
	for i, genTransaction := range response.Transactions {
		t[i] = models.FromGenTransaction(genTransaction)
	}

	render.JSON(w, r, t)
}
