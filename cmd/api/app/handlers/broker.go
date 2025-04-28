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

// CreateBroker godoc
//
//	@Id				CreateBroker
//
//	@Summary		Create a new broker
//	@Description	Create a new broker. (Permission: <b>admin.brokers.create</b>)
//	@Tags			Broker
//	@Accept			json
//	@Produce		json
//	@Param			broker	body	models.Broker	true	"broker (json)"
//	@Security		Bearer
//	@Success		200	{object}	models.Broker			"broker"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/broker [post]
func CreateBroker(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var broker models.Broker
	err := json.NewDecoder(r.Body).Decode(&broker)
	if err != nil {
		zap.L().Warn("Broker json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create gRPC protogen.CreateBrokerRequest
	brokerUserRequest := &protogen.CreateBrokerRequest{
		Name:     broker.Name,
		Disabled: broker.Disabled,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the Broker
	response, err := clients.C().Broker().CreateBroker(ctx, brokerUserRequest)
	if err != nil {
		zap.L().Error("Create Broker", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.JSON(w, r, models.FromProtogenBroker(response.Broker))
}

// GetBroker godoc
//
//	@Id				GetBroker
//
//	@Summary		Get a broker
//	@Description	Gets a broker.
//	@Tags			Broker
//	@Produce		json
//	@Param			id	path	string	true	"broker id"
//	@Security		Bearer
//	@Success		200	{object}	models.Broker			"broker"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		404	{object}	render.ErrorResponse	"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/broker/{id} [get]
func GetBroker(w http.ResponseWriter, r *http.Request) {
	// Retrieve brokerID
	brokerID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Create gRPC protogen.GetBrokerRequest
	brokerUserRequest := &protogen.GetBrokerRequest{
		Id: brokerID.String(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get the Broker
	response, err := clients.C().Broker().GetBroker(ctx, brokerUserRequest)
	if err != nil {
		zap.L().Error("Create Broker", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.JSON(w, r, models.FromProtogenBroker(response.Broker))
}

// UpdateBroker godoc
//
//	@Id				UpdateBroker
//
//	@Summary		Update a broker
//	@Description	Updates a broker. (Permission: <b>admin.brokers.update</b>)
//	@Tags			Broker
//	@Accept			json
//	@Produce		json
//	@Param			id			path	string					true	"broker ID"
//	@Param			broker		body	models.Broker			true	"broker (json)"
//	@Security		Bearer
//	@Success		200	{object}	models.Broker			"broker"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		404	{object}	render.ErrorResponse	"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/broker/{id} [put]
func UpdateBroker(w http.ResponseWriter, r *http.Request) {
	// Retrieve brokerID
	brokerID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Parse request body
	var broker models.Broker
	err := json.NewDecoder(r.Body).Decode(&broker)
	if err != nil {
		zap.L().Warn("Broker json decode", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Create gRPC protogen.UpdateBrokerRequest
	brokerUserRequest := &protogen.UpdateBrokerRequest{
		Id:       brokerID.String(),
		Name:     broker.Name,
		Disabled: broker.Disabled,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Update the Broker
	response, err := clients.C().Broker().UpdateBroker(ctx, brokerUserRequest)
	if err != nil {
		zap.L().Error("Update Broker", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.JSON(w, r, models.FromProtogenBroker(response.Broker))
}

// DeleteBroker godoc
//
//	@Id				DeleteBroker
//
//	@Summary		Delete a broker
//	@Description	Deletes a broker. (Permission: <b>admin.brokers.delete</b>)
//	@Tags			Broker
//	@Produce		json
//	@Param			id	path	string	true	"broker ID"
//	@Security		Bearer
//	@Success		200	{object}	string					"Status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		404	{object}	render.ErrorResponse	"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/broker/{id} [delete]
func DeleteBroker(w http.ResponseWriter, r *http.Request) {
	// Retrieve brokerID
	brokerID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Create gRPC protogen.DeleteBrokerRequest
	brokerUserRequest := &protogen.DeleteBrokerRequest{
		Id: brokerID.String(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Delete the Broker
	_, err := clients.C().Broker().DeleteBroker(ctx, brokerUserRequest)
	if err != nil {
		zap.L().Error("Delete Broker", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.OK(w, r)
}

// ListBrokers godoc
//
//	@Id				ListBrokers
//
//	@Summary		Get all brokers
//	@Description	Gets a list of all brokers.
//	@Tags			Broker
//	@Produce		json
//	@Param			enabled	query	string	false	"enabled only"
//	@Security		Bearer
//	@Success		200	{array}		models.Broker			"list of brokers"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/broker [get]
func ListBrokers(w http.ResponseWriter, r *http.Request) {
	enabled, ok := U().ParseParamBool(w, r, "enabled")
	if !ok {
		return
	}

	// Create gRPC protogen.ListBrokersRequest
	brokerUserRequest := &protogen.ListBrokersRequest{
		EnabledOnly: enabled,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// List the Broker
	response, err := clients.C().Broker().ListBrokers(ctx, brokerUserRequest)
	if err != nil {
		zap.L().Error("List Broker", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map gRPC response to Brokers array
	brokers := make([]models.Broker, len(response.Brokers))
	for i, protogenBroker := range response.Brokers {
		brokers[i] = models.FromProtogenBroker(protogenBroker)
	}

	render.JSON(w, r, brokers)
}
