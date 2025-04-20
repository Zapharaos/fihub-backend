package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/models"
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
	var userBrokerInput models.BrokerUserInput
	err := json.NewDecoder(r.Body).Decode(&userBrokerInput)
	if err != nil {
		zap.L().Warn("BrokerUser json decode", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Validate user
	if ok, err := userBrokerInput.IsValid(); !ok {
		zap.L().Warn("BrokerUser is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Convert to BrokerUser
	userBroker := userBrokerInput.ToUser()
	userBroker.UserID = user.ID

	// Retrieve broker to check if it exists
	broker, exists, err := brokers.R().B().Get(userBroker.Broker.ID)
	if err != nil {
		zap.L().Error("Get broker", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		zap.L().Warn("broker not found",
			zap.String("UserID", userBroker.UserID.String()),
			zap.String("BrokerID", userBroker.Broker.ID.String()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if broker is disabled
	if broker.Disabled {
		zap.L().Warn("Broker is disabled", zap.String("BrokerID", userBroker.Broker.ID.String()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Verify userBroker existence
	exists, err = brokers.R().U().Exists(userBroker)
	if err != nil {
		zap.L().Error("Check userBroker exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		zap.L().Warn("BrokerUser already exists", zap.String("BrokerID", userBroker.Broker.ID.String()))
		render.BadRequest(w, r, fmt.Errorf("broker-used"))
		return
	}

	// Create userBroker
	err = brokers.R().U().Create(userBroker)
	if err != nil {
		zap.L().Error("PostUserBroker.Create", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get userBrockers back from database
	userBrockers, err := brokers.R().U().GetAll(user.ID)
	if err != nil {
		zap.L().Error("Cannot get user brokers", zap.String("uuid", user.ID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(userBrockers) == 0 {
		zap.L().Error("BrokerUser broker not found after creation", zap.String("uuid", user.ID.String()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, userBrockers)
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

	// Build userBroker
	userBroker := models.BrokerUser{
		UserID: user.ID,
		Broker: models.Broker{ID: brokerID},
	}

	// Verify userBroker existence
	exists, err := brokers.R().U().Exists(userBroker)
	if err != nil {
		zap.L().Error("Check userBroker exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		zap.L().Warn("BrokerUser not found",
			zap.String("UserID", user.ID.String()),
			zap.String("BrokerID", brokerID.String()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Remove broker
	err = brokers.R().U().Delete(userBroker)
	if err != nil {
		zap.L().Error("Cannot remove user broker", zap.String("uuid", user.ID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.OK(w, r)
}

// GetUserBrokers godoc
//
//	@Id				GetUserBrokers
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
func GetUserBrokers(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user from the context
	user, ok := U().GetUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get userBrockers back from database
	userBrockers, err := brokers.R().U().GetAll(user.ID)
	if err != nil {
		zap.L().Error("Cannot get user brokers", zap.String("uuid", user.ID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, userBrockers)
}
