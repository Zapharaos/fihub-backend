package handlers

import (
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

// CreateBrokerImage godoc
//
//	@Id				CreateBrokerImage
//
//	@Summary		Create a new broker image
//	@Description	Create a new broker image. (Permission: <b>admin.brokers.create</b>)
//	@Tags			BrokerImages
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		string	true	"broker ID"
//	@Param			file	formData	file	true	"image file"
//	@Security		Bearer
//	@Success		200	{object}	models.BrokerImage		"broker image"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers/{id}/image [post]
func CreateBrokerImage(w http.ResponseWriter, r *http.Request) {

	if !U().CheckPermission(w, r, "admin.brokers.create") {
		return
	}

	// Get the broker ID
	brokerID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Parse the multipart form
	data, name, ok := U().ReadImage(w, r)
	if !ok {
		return
	}

	// Create the broker image
	brokerImageInput := models.BrokerImage{
		ID:       uuid.New(),
		BrokerID: brokerID,
		Name:     name,
		Data:     data,
	}

	// Validate the broker image
	if ok, err := brokerImageInput.IsValid(); !ok {
		zap.L().Warn("Broker image is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Verify broker has no image
	ok, err := repositories.R().B().HasImage(brokerID)
	if err != nil {
		zap.L().Error("HasImageBroker", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ok {
		zap.L().Warn("Broker already has an image", zap.String("broker_id", brokerID.String()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create the broker image
	err = repositories.R().I().Create(brokerImageInput)
	if err != nil {
		zap.L().Error("PostBrokerImage.Create", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the broker image back from the database
	brokerImage, found, err := repositories.R().I().Get(brokerImageInput.ID)
	if err != nil {
		zap.L().Error("Cannot get broker image", zap.String("uuid", brokerImageInput.ID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		zap.L().Error("Broker image not found after create", zap.String("uuid", brokerImageInput.ID.String()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the broker image
	err = repositories.R().B().SetImage(brokerID, brokerImageInput.ID)
	if err != nil {
		zap.L().Error("SetImageBroker.Create", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the broker image
	render.JSON(w, r, brokerImage)
}

// GetBrokerImage godoc
//
//	@Id				GetBrokerImage
//
//	@Summary		Get a broker image
//	@Description	Get a broker image.
//	@Tags			BrokerImages
//	@Accept			json
//	@Produce		image/jpeg
//	@Produce		image/png
//	@Param			id			path	string	true	"broker ID"
//	@Param			image_id	path	string	true	"image ID"
//	@Security		Bearer
//	@Success		200	{file}		file					"image"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		404	{object}	render.ErrorResponse	"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers/{id}/image/{image_id} [get]
func GetBrokerImage(w http.ResponseWriter, r *http.Request) {
	imageID, ok := U().ParseParamUUID(w, r, "image_id")
	if !ok {
		return
	}

	// Get the broker image
	brokerImage, found, err := repositories.R().I().Get(imageID)
	if err != nil {
		zap.L().Error("Cannot get brokerImage", zap.String("image_id", imageID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		zap.L().Warn("BrokerImage does not exist", zap.String("image_id", imageID.String()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Set headers
	w.Header().Set("Content-Disposition", "inline; filename="+brokerImage.Name)

	// Write the image
	_, err = w.Write(brokerImage.Data)
	if err != nil {
		zap.L().Error("Cannot write image", zap.Error(err))
		render.Error(w, r, err, "Write image")
		return
	}
}

// UpdateBrokerImage godoc
//
//	@Id				UpdateBrokerImage
//
//	@Summary		Update a broker image
//	@Description	Update a broker image. (Permission: <b>admin.brokers.update</b>)
//	@Tags			BrokerImages
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id			path		string	true	"broker ID"
//	@Param			image_id	path		string	true	"image ID"
//	@Param			file		formData	file	true	"image file"
//	@Security		Bearer
//	@Success		200	{object}	models.BrokerImage		"broker image"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		404	{object}	render.ErrorResponse	"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers/{id}/image/{image_id} [put]
func UpdateBrokerImage(w http.ResponseWriter, r *http.Request) {

	if !U().CheckPermission(w, r, "admin.brokers.update") {
		return
	}

	// Get the broker ID
	brokerID, imageID, ok := U().ParseUUIDPair(w, r, "image_id")
	if !ok {
		return
	}

	// Parse the multipart form
	data, name, ok := U().ReadImage(w, r)
	if !ok {
		return
	}

	// Create the broker image
	brokerImageInput := models.BrokerImage{
		ID:       imageID,
		BrokerID: brokerID,
		Name:     name,
		Data:     data,
	}

	// Validate the broker image
	if ok, err := brokerImageInput.IsValid(); !ok {
		zap.L().Warn("Broker image is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Verify imageBroker existence
	exists, err := repositories.R().I().Exists(brokerID, imageID)
	if err != nil {
		zap.L().Error("Check imageBroker exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		zap.L().Warn("ImageBroker not found", zap.String("broker_id", brokerID.String()), zap.String("image_id", imageID.String()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Update the broker image
	err = repositories.R().I().Update(brokerImageInput)
	if err != nil {
		zap.L().Error("UpdateBrokerImage.Update", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the broker image back from the database
	brokerImage, found, err := repositories.R().I().Get(brokerImageInput.ID)
	if err != nil {
		zap.L().Error("Cannot get broker image", zap.String("uuid", brokerImageInput.ID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		zap.L().Error("Broker image not found after create", zap.String("uuid", brokerImageInput.ID.String()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the broker image
	render.JSON(w, r, brokerImage)
}

// DeleteBrokerImage godoc
//
//	@Id				DeleteBrokerImage
//
//	@Summary		Delete a broker image
//	@Description	Delete a broker image. (Permission: <b>admin.brokers.delete</b>)
//	@Tags			BrokerImages
//	@Produce		json
//	@Param			id			path	string	true	"broker ID"
//	@Param			image_id	path	string	true	"image ID"
//	@Security		Bearer
//	@Success		200	{object}	string					"OK"`
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		404	{object}	render.ErrorResponse	"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers/{id}/image/{image_id} [delete]
func DeleteBrokerImage(w http.ResponseWriter, r *http.Request) {

	if !U().CheckPermission(w, r, "admin.brokers.delete") {
		return
	}

	// Get the broker ID
	brokerID, imageID, ok := U().ParseUUIDPair(w, r, "image_id")
	if !ok {
		return
	}

	// Verify imageBroker existence
	exists, err := repositories.R().I().Exists(brokerID, imageID)
	if err != nil {
		zap.L().Error("Check imageBroker exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		zap.L().Warn("ImageBroker not found", zap.String("broker_id", brokerID.String()), zap.String("image_id", imageID.String()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Delete the broker image
	err = repositories.R().I().Delete(imageID)
	if err != nil {
		zap.L().Error("DeleteBrokerImage.Delete", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.OK(w, r)
}
