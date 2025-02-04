package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/handlers/render"
	"github.com/Zapharaos/fihub-backend/internal/transactions"
	"go.uber.org/zap"
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
// @Param 				transaction body 	transactions.TransactionInput true 	"transaction (json)"
// @Security 			Bearer
// @Success 			200 {object} 		transactions.Transaction 		"transaction"
// @Failure 			400 {object} 		render.ErrorResponse 			"Bad Request"
// @Failure 			401 {string} 		string 							"Permission denied"
// @Failure 			500 {object} 		render.ErrorResponse 			"Internal Server Error"
// @Router /api/v1/transactions [post]
func CreateTransaction(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user from the context
	user, ok := getUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse request body
	var transactionInput transactions.TransactionInput
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

	// Set the user ID + PriceUnit
	transactionInput.UserID = user.ID
	transactionInput.PriceUnit = transactionInput.Price / transactionInput.Quantity

	// check if userBroker exists
	exists, err := brokers.R().U().Exists(brokers.User{UserID: user.ID, Broker: brokers.Broker{ID: transactionInput.BrokerID}})
	if err != nil {
		zap.L().Error("Check userBroker exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		zap.L().Warn("User not found")
		render.BadRequest(w, r, fmt.Errorf("broker-invalid"))
		return
	}

	// Create the transaction using the transactions.r() repository
	transactionID, err := transactions.R().Create(transactionInput)
	if err != nil {
		zap.L().Error("Create transaction", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get transaction back from database
	transaction, ok, err := transactions.R().Get(transactionID)
	if err != nil {
		zap.L().Error("Cannot get transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		zap.L().Error("Transaction not found after creation", zap.String("uuid", transactionID.String()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
// @Success 		200 {object} 	transactions.Transaction 	"transaction"
// @Failure 		400 {object} 	render.ErrorResponse 		"Bad Request"
// @Failure 		401 {string} 	string 						"Permission denied"
// @Failure 		404 {object} 	render.ErrorResponse 		"Not Found"
// @Failure 		500 {object} 	render.ErrorResponse 		"Internal Server Error"
// @Router /api/v1/transactions/{id} [get]
func GetTransaction(w http.ResponseWriter, r *http.Request) {

	// Retrieve transactionID
	transactionID, ok := parseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Get transaction
	transaction, ok, err := transactions.R().Get(transactionID)
	if err != nil {
		zap.L().Error("Cannot get transaction", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		zap.L().Warn("Transaction not found", zap.String("uuid", transactionID.String()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Verify that the transaction belongs to the user
	user, ok := getUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if transaction.UserID != user.ID {
		zap.L().Warn("Transaction does not belong to user", zap.String("uuid", transactionID.String()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
// @Param 			transaction body transactions.Transaction true "transaction (json)"
// @Security 		Bearer
// @Success 		200 {object} 	transactions.Transaction 	"transaction"
// @Failure 		400 {object} 	render.ErrorResponse 		"Bad Request"
// @Failure 		401 {string} 	string 						"Permission denied"
// @Failure			404	{object}	render.ErrorResponse		"Not Found"
// @Failure 		500 {object} 	render.ErrorResponse 		"Internal Server Error"
// @Router /api/v1/transactions/{id} [put]
func UpdateTransaction(w http.ResponseWriter, r *http.Request) {

	// Retrieve transactionID
	transactionID, ok := parseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Get the authenticated user from the context
	user, ok := getUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse request body
	var transactionInput transactions.TransactionInput
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

	// Set the transaction ID + PriceUnit
	transactionInput.ID = transactionID
	transactionInput.PriceUnit = transactionInput.Price / transactionInput.Quantity

	// Verify that the transaction belongs to the user
	oldTransaction, ok, err := transactions.R().Get(transactionID)
	if err != nil {
		zap.L().Error("Cannot get transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		zap.L().Warn("Transaction not found", zap.String("uuid", transactionID.String()))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if oldTransaction.UserID != user.ID {
		zap.L().Warn("Transaction does not belong to user", zap.String("uuid", transactionID.String()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Verify that the userBroker exists
	exists, err := brokers.R().U().Exists(brokers.User{UserID: user.ID, Broker: brokers.Broker{ID: transactionInput.BrokerID}})
	if err != nil {
		zap.L().Error("Check userBroker exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		zap.L().Warn("User not found")
		render.BadRequest(w, r, fmt.Errorf("broker-invalid"))
		return
	}

	// Update the transaction using the transactions.r() repository
	err = transactions.R().Update(transactionInput)
	if err != nil {
		zap.L().Error("Update transaction", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get transaction back from database
	transaction, ok, err := transactions.R().Get(transactionID)
	if err != nil {
		zap.L().Error("Cannot get transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		zap.L().Error("Transaction not found after update", zap.String("uuid", transactionID.String()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
// @Failure 		400 {object} 	render.ErrorResponse 	"Bad Request"
// @Failure 		401 {string} 	string 					"Permission denied"
// @Failure			404	{object}	render.ErrorResponse	"Not Found"
// @Failure 		500 {object} 	render.ErrorResponse 	"Internal Server Error"
// @Router /api/v1/transactions/{id} [delete]
func DeleteTransaction(w http.ResponseWriter, r *http.Request) {

	// Retrieve transactionID
	transactionID, ok := parseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Get the authenticated user from the context
	user, ok := getUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Verify that the transaction belongs to the user
	transaction, ok, err := transactions.R().Get(transactionID)
	if err != nil {
		zap.L().Error("Cannot get transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		zap.L().Warn("Transaction not found", zap.String("uuid", transactionID.String()))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if transaction.UserID != user.ID {
		zap.L().Warn("Transaction does not belong to user", zap.String("uuid", transactionID.String()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Remove transaction
	err = transactions.R().Delete(transactions.Transaction{ID: transactionID, UserID: user.ID})
	if err != nil {
		zap.L().Error("Cannot remove transaction", zap.String("uuid", transactionID.String()), zap.Error(err))
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
	user, ok := getUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get all transactions
	userTransactions, err := transactions.R().GetAll(user.ID)
	if err != nil {
		zap.L().Error("Cannot get transactions", zap.String("uuid", user.ID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, userTransactions)
}
