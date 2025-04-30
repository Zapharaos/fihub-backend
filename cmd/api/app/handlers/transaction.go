package handlers

import (
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/internal/mappers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
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
// @Param 				transaction body 	models.TransactionInput true 	"transaction (json)"
// @Security 			Bearer
// @Success 			200 {object} 		models.Transaction 		"transaction"
// @Failure 			400 {object} 		render.ErrorResponse 			"Bad PasswordRequest"
// @Failure 			401 {string} 		string 							"Permission denied"
// @Failure 			500 {object} 		render.ErrorResponse 			"Internal Server Error"
// @Router /api/v1/transaction [post]
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

	// Verify BrokerUser existence
	_, err = clients.C().Broker().GetBrokerUser(r.Context(), &protogen.GetBrokerUserRequest{
		UserId:   user.ID.String(),
		BrokerId: transactionInput.BrokerID.String(),
	})
	if err != nil {
		zap.L().Error("Get BrokerUser", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map TransactionInput to gRPC ValidateTransactionRequest
	transactionRequest := &protogen.CreateTransactionRequest{
		UserId:          user.ID.String(),
		BrokerId:        transactionInput.BrokerID.String(),
		Date:            timestamppb.New(transactionInput.Date),
		TransactionType: mappers.TransactionTypeToProto(transactionInput.Type),
		Asset:           transactionInput.Asset,
		Quantity:        transactionInput.Quantity,
		Price:           transactionInput.Price,
		Fee:             transactionInput.Fee,
	}

	// Create the transaction
	response, err := clients.C().Transaction().CreateTransaction(r.Context(), transactionRequest)
	if err != nil {
		zap.L().Error("Create transaction", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map gRPC response to Transaction
	transaction := mappers.TransactionFromProto(response.Transaction)

	// Retrieve broker object
	responseBroker, err := clients.C().Broker().GetBroker(r.Context(), &protogen.GetBrokerRequest{
		Id: transaction.Broker.ID.String(),
	})
	if err != nil {
		zap.L().Error("Get broker", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Put the broker object into the transaction
	transaction.Broker = mappers.BrokerFromProto(responseBroker.Broker)
	render.JSON(w, r, transaction)
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
// @Success 		200 {object} 	models.Transaction 	"transaction"
// @Failure 		400 {object} 	render.ErrorResponse 		"Bad PasswordRequest"
// @Failure 		401 {string} 	string 						"Permission denied"
// @Failure 		404 {object} 	render.ErrorResponse 		"Not Found"
// @Failure 		500 {object} 	render.ErrorResponse 		"Internal Server Error"
// @Router /api/v1/transaction/{id} [get]
func GetTransaction(w http.ResponseWriter, r *http.Request) {

	// Retrieve transactionID
	transactionID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Retrieve the transaction
	response, err := clients.C().Transaction().GetTransaction(r.Context(), &protogen.GetTransactionRequest{
		TransactionId: transactionID.String(),
	})
	if err != nil {
		zap.L().Error("Get transaction", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Get the authenticated user from the context
	user, ok := U().GetUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Map gRPC response to Transaction
	transaction := mappers.TransactionFromProto(response.Transaction)

	// Check if the transaction belongs to the user
	if transaction.UserID != user.ID {
		zap.L().Warn("Transaction does not belong to user", zap.String("uuid", transactionID.String()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Retrieve broker object
	responseBroker, err := clients.C().Broker().GetBroker(r.Context(), &protogen.GetBrokerRequest{
		Id: transaction.Broker.ID.String(),
	})
	if err != nil {
		zap.L().Error("Get broker", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Put the broker object in the transaction
	transaction.Broker = mappers.BrokerFromProto(responseBroker.Broker)
	render.JSON(w, r, transaction)
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
// @Param 			transaction body models.Transaction true "transaction (json)"
// @Security 		Bearer
// @Success 		200 {object} 	models.Transaction 	"transaction"
// @Failure 		400 {object} 	render.ErrorResponse 		"Bad PasswordRequest"
// @Failure 		401 {string} 	string 						"Permission denied"
// @Failure			404	{object}	render.ErrorResponse		"Not Found"
// @Failure 		500 {object} 	render.ErrorResponse 		"Internal Server Error"
// @Router /api/v1/transaction/{id} [put]
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

	// Verify BrokerUser existence
	_, err = clients.C().Broker().GetBrokerUser(r.Context(), &protogen.GetBrokerUserRequest{
		UserId:   user.ID.String(),
		BrokerId: transactionInput.BrokerID.String(),
	})
	if err != nil {
		zap.L().Error("Get BrokerUser", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map TransactionInput to gRPC ValidateTransactionRequest
	transactionRequest := &protogen.UpdateTransactionRequest{
		TransactionId:   transactionID.String(),
		UserId:          user.ID.String(),
		BrokerId:        transactionInput.BrokerID.String(),
		Date:            timestamppb.New(transactionInput.Date),
		TransactionType: mappers.TransactionTypeToProto(transactionInput.Type),
		Asset:           transactionInput.Asset,
		Quantity:        transactionInput.Quantity,
		Price:           transactionInput.Price,
		Fee:             transactionInput.Fee,
	}

	// Create the transaction
	response, err := clients.C().Transaction().UpdateTransaction(r.Context(), transactionRequest)
	if err != nil {
		zap.L().Error("Update transaction", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map gRPC response to Transaction
	transaction := mappers.TransactionFromProto(response.Transaction)

	// Retrieve broker object
	responseBroker, err := clients.C().Broker().GetBroker(r.Context(), &protogen.GetBrokerRequest{
		Id: transaction.Broker.ID.String(),
	})
	if err != nil {
		zap.L().Error("Get broker", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Put the broker object in the transaction
	transaction.Broker = mappers.BrokerFromProto(responseBroker.Broker)
	render.JSON(w, r, transaction)
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
// @Router /api/v1/transaction/{id} [delete]
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

	// Retrieve the transaction
	_, err := clients.C().Transaction().DeleteTransaction(r.Context(), &protogen.DeleteTransactionRequest{
		TransactionId: transactionID.String(),
		UserId:        user.ID.String(),
	})
	if err != nil {
		zap.L().Error("Delete transaction", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.OK(w, r)
}

// ListTransactions godoc
//
// @Id 				GetTransactions
//
// @Summary 		List all transactions
// @Description 	Gets a list of all transactions.
// @Tags 			Transactions
// @Produce 		json
// @Security 		Bearer
// @Success 		200 {array} 	models.Transaction 	"List of transactions"
// @Failure 		401 {string} 	string 						"Permission denied"
// @Failure 		500 {object} 	render.ErrorResponse 		"Internal Server Error"
// @Router /api/v1/transaction [get]
func ListTransactions(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user from the context
	user, ok := U().GetUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// List transactions
	response, err := clients.C().Transaction().ListTransactions(r.Context(), &protogen.ListTransactionsRequest{
		UserId: user.ID.String(),
	})
	if err != nil {
		zap.L().Error("List transactions", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Retrieve broker objects
	responseBrokers, err := clients.C().Broker().ListBrokers(r.Context(), &protogen.ListBrokersRequest{
		EnabledOnly: false,
	})
	if err != nil {
		zap.L().Error("List brokers", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Create a map of brokers indexed by broker ID for faster lookup
	brokersMap := make(map[string]models.Broker)
	for _, b := range responseBrokers.Brokers {
		broker := mappers.BrokerFromProto(b)
		brokersMap[broker.ID.String()] = broker
	}

	// Map gRPC response to Transaction array
	t := make([]models.Transaction, len(response.Transactions))
	for i, protogenTransaction := range response.Transactions {
		transaction := mappers.TransactionFromProto(protogenTransaction)
		transaction.Broker = brokersMap[transaction.Broker.ID.String()]
		t[i] = transaction
	}

	render.JSON(w, r, t)
}
