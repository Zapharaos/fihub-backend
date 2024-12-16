package handlers

import (
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/handlers/render"
	"go.uber.org/zap"
	"net/http"
)

// CreateBroker godoc
//
//	@Id				CreateBroker
//
//	@Summary		Create a new broker
//	@Description	Create a new broker.
//	@Tags			Brokers
//	@Accept			json
//	@Produce		json
//	@Param			broker	body	brokers.Broker	true	"broker (json)"
//	@Security		Bearer
//	@Success		200	{object}	brokers.Broker			"broker"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers [post]
func CreateBroker(w http.ResponseWriter, r *http.Request) {

	// TODO : permissions

	// Parse request body
	var broker brokers.Broker
	err := json.NewDecoder(r.Body).Decode(&broker)
	if err != nil {
		zap.L().Warn("Broker json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate the broker
	if valid, err := broker.IsValid(); !valid {
		zap.L().Warn("Broker is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Create the broker
	brokerID, err := brokers.R().B().Create(broker)
	if err != nil {
		zap.L().Warn("Create broker", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the broker from the database
	broker, found, err := brokers.R().B().Get(brokerID)
	if err != nil {
		zap.L().Error("Cannot get broker", zap.String("uuid", brokerID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		zap.L().Error("Broker not found after creation", zap.String("uuid", brokerID.String()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, broker)
}

// GetBroker godoc
//
//	@Id				GetBroker
//
//	@Summary		Get a broker
//	@Description	Gets a broker.
//	@Tags			Brokers
//	@Produce		json
//	@Param			id	path	string	true	"broker id"
//	@Security		Bearer
//	@Success		200	{object}	brokers.Broker			"broker"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	render.ErrorResponse	"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers/{id} [get]
func GetBroker(w http.ResponseWriter, r *http.Request) {

	// Retrieve brokerID
	brokerID, ok := ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Get the broker from the database
	broker, found, err := brokers.R().B().Get(brokerID)
	if err != nil {
		zap.L().Error("Cannot get broker", zap.String("uuid", brokerID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		zap.L().Error("Broker not found", zap.String("uuid", brokerID.String()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	render.JSON(w, r, broker)
}

// UpdateBroker godoc
//
//	@Id				UpdateBroker
//
//	@Summary		Update a broker
//	@Description	Updates a broker.
//	@Tags			Brokers
//	@Accept			json
//	@Produce		json
//	@Param			id			path	string					true	"broker ID"
//	@Param			broker		body	brokers.Broker			true	"broker (json)"
//	@Security		Bearer
//	@Success		200	{object}	brokers.Broker			"broker"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers/{id} [put]
func UpdateBroker(w http.ResponseWriter, r *http.Request) {

	// TODO : permissions

	// Retrieve brokerID
	brokerID, ok := ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Parse request body
	var broker brokers.Broker
	err := json.NewDecoder(r.Body).Decode(&broker)
	if err != nil {
		zap.L().Warn("Broker json decode", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	broker.ID = brokerID

	// Validate the broker
	if valid, err := broker.IsValid(); !valid {
		zap.L().Warn("Broker is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Update the broker
	err = brokers.R().B().Update(broker)
	if err != nil {
		zap.L().Warn("Update broker", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the broker from the database
	broker, found, err := brokers.R().B().Get(brokerID)
	if err != nil {
		zap.L().Error("Cannot get broker", zap.String("uuid", brokerID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		zap.L().Error("Broker not found after update", zap.String("uuid", brokerID.String()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, broker)
}

// DeleteBroker godoc
//
//	@Id				DeleteBroker
//
//	@Summary		Delete a broker
//	@Description	Deletes a broker.
//	@Tags			Brokers
//	@Produce		json
//	@Param			id	path	string	true	"broker ID"
//	@Security		Bearer
//	@Success		200	{object}	string					"Status OK"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers/{id} [delete]
func DeleteBroker(w http.ResponseWriter, r *http.Request) {

	// TODO : permissions

	// Retrieve brokerID
	brokerID, ok := ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Delete the broker
	err := brokers.R().B().Delete(brokerID)
	if err != nil {
		zap.L().Warn("Delete broker", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.OK(w, r)
}

// GetBrokers godoc
//
//	@Id				GetBrokers
//
//	@Summary		Get all brokers
//	@Description	Gets a list of all brokers.
//	@Tags			Brokers
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}		brokers.Broker			"list of brokers"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers [get]
func GetBrokers(w http.ResponseWriter, r *http.Request) {

	result, err := brokers.R().B().GetAll()
	if err != nil {
		render.Error(w, r, err, "Get brokers")
		return
	}

	render.JSON(w, r, result)
}
