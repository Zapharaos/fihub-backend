package handlers

import (
	"context"
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
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
//	@Tags			BrokerUser
//	@Accept			json
//	@Produce		json
//	@Param			userBroker	body	models.BrokerUserInput	true	"userBroker (json)"
//	@Security		Bearer
//	@Success		200	{array}		models.BrokerUser		"Updated list of user brokers"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/brokers [post]
func CreateUserBroker(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the context
	user, ok := U().GetUserFromContext(r)
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
		UserId:   user.ID.String(),
		BrokerId: brokerUserInput.BrokerID,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the BrokerUser
	response, err := clients.C().Broker().CreateBrokerUser(ctx, brokerUserRequest)
	if err != nil {
		zap.L().Error("Create BrokerUser", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map gRPC response to UserBrokers array
	userBrokers := make([]models.BrokerUser, len(response.UserBrokers))
	for i, protogenBrokerUser := range response.UserBrokers {
		userBrokers[i] = models.FromProtogenBrokerUser(protogenBrokerUser)
	}

	render.JSON(w, r, userBrokers)
}

// DeleteUserBroker godoc
//
//	@Id				DeleteUserBroker
//
//	@Summary		Delete a user broker
//	@Description	Delete a user broker.
//	@Tags			BrokerUser
//	@Produce		json
//	@Param			id	path	string	true	"broker ID"
//	@Security		Bearer
//	@Success		200	{array}		string					"Status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/brokers/{id} [delete]
func DeleteUserBroker(w http.ResponseWriter, r *http.Request) {
	// Retrieve brokerID
	brokerID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Get the authenticated user from the context
	user, ok := U().GetUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Map models.BrokerUserInput to gRPC protogen.DeleteBrokerUserRequest
	brokerUserRequest := &protogen.DeleteBrokerUserRequest{
		UserId:   user.ID.String(),
		BrokerId: brokerID.String(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Delete the BrokerUser
	_, err := clients.C().Broker().DeleteBrokerUser(ctx, brokerUserRequest)
	if err != nil {
		zap.L().Error("Delete BrokerUser", zap.Error(err))
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
//	@Tags			BrokerUser
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}		models.BrokerUser		"List of brokers"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/brokers [get]
func ListUserBrokers(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user from the context
	user, ok := U().GetUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Map props to gRPC protogen.GetUserBrokersRequest
	brokerUserRequest := &protogen.ListUserBrokersRequest{
		UserId: user.ID.String(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get the userBrokers
	response, err := clients.C().Broker().ListUserBrokers(ctx, brokerUserRequest)
	if err != nil {
		zap.L().Error("Get BrokerUsers", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map gRPC response to UserBrokers array
	userBrokers := make([]models.BrokerUser, len(response.UserBrokers))
	for i, protogenBrokerUser := range response.UserBrokers {
		userBrokers[i] = models.FromProtogenBrokerUser(protogenBrokerUser)
	}

	render.JSON(w, r, userBrokers)
}
