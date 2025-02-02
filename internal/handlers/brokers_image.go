package handlers

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/handlers/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
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
//	@Success		200	{object}	brokers.BrokerImage		"broker image"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers/{id}/image [post]
func CreateBrokerImage(w http.ResponseWriter, r *http.Request) {

	if !checkPermission(w, r, "admin.brokers.create") {
		return
	}

	// Get the broker ID
	brokerID, ok := parseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		zap.L().Warn("Form file", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}
	defer file.Close()

	// Read the file
	data, err := io.ReadAll(file)
	if err != nil {
		zap.L().Warn("Read file", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Check the MIME type
	mimeType := http.DetectContentType(data)
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		zap.L().Warn("Invalid MIME type", zap.String("mimeType", mimeType))
		render.BadRequest(w, r, errors.New("invalid-type"))
		return
	}

	// Create the broker image
	brokerImageInput := brokers.BrokerImage{
		ID:       uuid.New(),
		BrokerID: brokerID,
		Name:     header.Filename,
		Data:     data,
	}

	// Validate the broker image
	if ok, err = brokerImageInput.IsValid(); !ok {
		zap.L().Warn("Broker image is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Verify broker has no image
	ok, err = brokers.R().B().HasImage(brokerID)
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
	err = brokers.R().I().Create(brokerImageInput)
	if err != nil {
		zap.L().Error("PostBrokerImage.Create", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the broker image back from the database
	brokerImage, found, err := brokers.R().I().Get(brokerImageInput.ID)
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
	err = brokers.R().B().SetImage(brokerID, brokerImageInput.ID)
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
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	render.ErrorResponse	"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers/{id}/image/{image_id} [get]
func GetBrokerImage(w http.ResponseWriter, r *http.Request) {
	brokerID, imageID, ok := parseUUIDPair(w, r, "image_id")
	if !ok {
		return
	}

	// Verify imageBroker existence
	exists, err := brokers.R().I().Exists(brokerID, imageID)
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

	// Get the broker image
	brokerImage, found, err := brokers.R().I().Get(imageID)
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
//	@Success		200	{object}	brokers.BrokerImage		"broker image"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	render.ErrorResponse	"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers/{id}/image/{image_id} [put]
func UpdateBrokerImage(w http.ResponseWriter, r *http.Request) {

	if !checkPermission(w, r, "admin.brokers.update") {
		return
	}

	// Get the broker ID
	brokerID, imageID, ok := parseUUIDPair(w, r, "image_id")
	if !ok {
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		zap.L().Warn("Form file", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}
	defer file.Close()

	// Read the file
	data, err := io.ReadAll(file)
	if err != nil {
		zap.L().Warn("Read file", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Check the MIME type
	mimeType := http.DetectContentType(data)
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		zap.L().Warn("Invalid MIME type", zap.String("mimeType", mimeType))
		render.BadRequest(w, r, errors.New("invalid-type"))
		return
	}

	// Create the broker image
	brokerImageInput := brokers.BrokerImage{
		ID:       imageID,
		BrokerID: brokerID,
		Name:     header.Filename,
		Data:     data,
	}

	// Validate the broker image
	if ok, err = brokerImageInput.IsValid(); !ok {
		zap.L().Warn("Broker image is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Verify imageBroker existence
	exists, err := brokers.R().I().Exists(brokerID, imageID)
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
	err = brokers.R().I().Update(brokerImageInput)
	if err != nil {
		zap.L().Error("UpdateBrokerImage.Update", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the broker image back from the database
	brokerImage, found, err := brokers.R().I().Get(brokerImageInput.ID)
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
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	render.ErrorResponse	"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/brokers/{id}/image/{image_id} [delete]
func DeleteBrokerImage(w http.ResponseWriter, r *http.Request) {

	if !checkPermission(w, r, "admin.brokers.delete") {
		return
	}

	// Get the broker ID
	brokerID, imageID, ok := parseUUIDPair(w, r, "image_id")
	if !ok {
		return
	}

	// Verify imageBroker existence
	exists, err := brokers.R().I().Exists(brokerID, imageID)
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
	err = brokers.R().I().Delete(imageID)
	if err != nil {
		zap.L().Error("DeleteBrokerImage.Delete", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.OK(w, r)
}
