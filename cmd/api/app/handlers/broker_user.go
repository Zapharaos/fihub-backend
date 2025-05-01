package handlers

import (
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/internal/mappers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"go.uber.org/zap"
	"net/http"
)

// CreateUserBroker godoc
//
//	@Id				CreateUserBroker
//
//	@Summary		Create a new user broker
//	@Description	Create a new user broker.
//	@Tags			Broker, User
//	@Accept			json
//	@Produce		json
//	@Param			userBroker	body	models.BrokerUserInput	true	"userBroker (json)"
//	@Security		Bearer
//	@Success		200	{array}		models.BrokerUser		"Updated list of user brokers"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/broker/user [post]
func CreateUserBroker(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context
	userID, ok := U().GetUserIDFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse request body
	var brokerUserInput models.BrokerUserInput
	err := json.NewDecoder(r.Body).Decode(&brokerUserInput)
	if err != nil {
		zap.L().Warn("BrokerUser json decode", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Map models.BrokerUserInput to gRPC protogen.CreateBrokerUserRequest
	brokerUserRequest := &protogen.CreateBrokerUserRequest{
		UserId:   userID,
		BrokerId: brokerUserInput.BrokerID,
	}

	// Create the BrokerUser
	response, err := clients.C().Broker().CreateBrokerUser(r.Context(), brokerUserRequest)
	if err != nil {
		zap.L().Error("Create BrokerUser", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map gRPC response to UserBrokers array
	userBrokers := make([]models.BrokerUser, len(response.UserBrokers))
	for i, protogenBrokerUser := range response.UserBrokers {
		userBrokers[i] = mappers.BrokerUserFromProto(protogenBrokerUser)
	}

	render.JSON(w, r, userBrokers)
}

// DeleteUserBroker godoc
//
//	@Id				DeleteUserBroker
//
//	@Summary		Delete a user broker
//	@Description	Delete a user broker.
//	@Tags			Broker, User
//	@Produce		json
//	@Param			id	path	string	true	"broker ID"
//	@Security		Bearer
//	@Success		200	{array}		string					"Status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/broker/{id}/user [delete]
func DeleteUserBroker(w http.ResponseWriter, r *http.Request) {
	// Retrieve brokerID
	brokerID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Get the authenticated user from the context
	userID, ok := U().GetUserIDFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Map models.BrokerUserInput to gRPC protogen.DeleteBrokerUserRequest
	brokerUserRequest := &protogen.DeleteBrokerUserRequest{
		UserId:   userID,
		BrokerId: brokerID.String(),
	}

	// Delete the BrokerUser
	_, err := clients.C().Broker().DeleteBrokerUser(r.Context(), brokerUserRequest)
	if err != nil {
		zap.L().Error("Delete BrokerUser", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Delete Transactions related to the BrokerUser
	_, err = clients.C().Transaction().DeleteTransactionByBroker(r.Context(), &protogen.DeleteTransactionByBrokerRequest{
		UserId:   userID,
		BrokerId: brokerID.String(),
	})
	if err != nil {
		zap.L().Error("Delete Transactions by BrokerUser", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.OK(w, r)
}

// ListUserBrokers godoc
//
//	@Id				ListUserBrokers
//
//	@Summary		Get all user's brokers
//	@Description	Gets a list of all user's brokers.
//	@Tags			Broker, User
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}		models.BrokerUser		"List of brokers"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/broker/user [get]
func ListUserBrokers(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user from the context
	userID, ok := U().GetUserIDFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Map props to gRPC protogen.GetUserBrokersRequest
	brokerUserRequest := &protogen.ListUserBrokersRequest{
		UserId: userID,
	}

	// Get the userBrokers
	response, err := clients.C().Broker().ListUserBrokers(r.Context(), brokerUserRequest)
	if err != nil {
		zap.L().Error("Get BrokerUsers", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map gRPC response to UserBrokers array
	userBrokers := make([]models.BrokerUser, len(response.UserBrokers))
	for i, protogenBrokerUser := range response.UserBrokers {
		userBrokers[i] = mappers.BrokerUserFromProto(protogenBrokerUser)
	}

	render.JSON(w, r, userBrokers)
}
