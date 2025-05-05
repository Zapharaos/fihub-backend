package handlers

import (
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/gen/go/brokerpb"
	"github.com/Zapharaos/fihub-backend/internal/mappers"
	"go.uber.org/zap"
	"net/http"
)

// CreateBrokerImage godoc
//
//	@Id				CreateBrokerImage
//
//	@Summary		Create a new broker image
//	@Description	Create a new broker image. (Permission: <b>admin.brokers.create</b>)
//	@Tags			Broker, Image
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		string	true	"broker ID"
//	@Param			file	formData	file	true	"image file"
//	@Security		Bearer
//	@Success		200	{object}	models.BrokerImage		"broker image"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/broker/{id}/image [post]
func CreateBrokerImage(w http.ResponseWriter, r *http.Request) {
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

	// Create gRPC gen.CreateBrokerImageRequest
	brokerUserRequest := &brokerpb.CreateBrokerImageRequest{
		BrokerId: brokerID.String(),
		Name:     name,
		Data:     data,
	}

	// Create the BrokerImage
	response, err := clients.C().Broker().CreateBrokerImage(r.Context(), brokerUserRequest)
	if err != nil {
		zap.L().Error("Create BrokerImage", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Return the broker image
	render.JSON(w, r, mappers.BrokerImageFromProto(response.Image))
}

// GetBrokerImage godoc
//
//	@Id				GetBrokerImage
//
//	@Summary		Get a broker image
//	@Description	Get a broker image.
//	@Tags			Broker, Image
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
//	@Router			/api/v1/broker/{id}/image/{image_id} [get]
func GetBrokerImage(w http.ResponseWriter, r *http.Request) {
	imageID, ok := U().ParseParamUUID(w, r, "image_id")
	if !ok {
		return
	}

	// Create gRPC gen.GetBrokerImageRequest
	brokerUserRequest := &brokerpb.GetBrokerImageRequest{
		ImageId: imageID.String(),
	}

	// Get the BrokerImage
	response, err := clients.C().Broker().GetBrokerImage(r.Context(), brokerUserRequest)
	if err != nil {
		zap.L().Error("Get BrokerImage", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Set headers
	w.Header().Set("Content-Disposition", "inline; filename="+response.Name)

	// Write the image
	_, err = w.Write(response.Data)
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
//	@Tags			Broker, Image
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
//	@Router			/api/v1/broker/{id}/image/{image_id} [put]
func UpdateBrokerImage(w http.ResponseWriter, r *http.Request) {
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

	// Create gRPC gen.UpdateBrokerImageRequest
	brokerUserRequest := &brokerpb.UpdateBrokerImageRequest{
		ImageId:  imageID.String(),
		BrokerId: brokerID.String(),
		Name:     name,
		Data:     data,
	}

	// Create the BrokerImage
	response, err := clients.C().Broker().UpdateBrokerImage(r.Context(), brokerUserRequest)
	if err != nil {
		zap.L().Error("Create BrokerImage", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Return the broker image
	render.JSON(w, r, mappers.BrokerImageFromProto(response.Image))
}

// DeleteBrokerImage godoc
//
//	@Id				DeleteBrokerImage
//
//	@Summary		Delete a broker image
//	@Description	Delete a broker image. (Permission: <b>admin.brokers.delete</b>)
//	@Tags			Broker, Image
//	@Produce		json
//	@Param			id			path	string	true	"broker ID"
//	@Param			image_id	path	string	true	"image ID"
//	@Security		Bearer
//	@Success		200	{object}	string					"OK"`
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		404	{object}	render.ErrorResponse	"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/broker/{id}/image/{image_id} [delete]
func DeleteBrokerImage(w http.ResponseWriter, r *http.Request) {
	// Get the broker ID
	brokerID, imageID, ok := U().ParseUUIDPair(w, r, "image_id")
	if !ok {
		return
	}

	// Create gRPC gen.DeleteBrokerImageRequest
	brokerUserRequest := &brokerpb.DeleteBrokerImageRequest{
		ImageId:  imageID.String(),
		BrokerId: brokerID.String(),
	}

	// Create the BrokerImage
	_, err := clients.C().Broker().DeleteBrokerImage(r.Context(), brokerUserRequest)
	if err != nil {
		zap.L().Error("Create BrokerImage", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.OK(w, r)
}
