package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/brokerpb"
	"github.com/Zapharaos/fihub-backend/internal/mappers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateBrokerImage implements the CreateBrokerImage RPC method.
func (h *Service) CreateBrokerImage(ctx context.Context, req *brokerpb.CreateBrokerImageRequest) (*brokerpb.CreateBrokerImageResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.brokers.create")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return nil, err
	}

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetBrokerId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetBrokerId()), zap.Error(err))
		return &brokerpb.CreateBrokerImageResponse{}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Construct the BrokerImageInput object from the request
	brokerImageInput := models.BrokerImage{
		ID:       uuid.New(),
		BrokerID: brokerID,
		Name:     req.GetName(),
		Data:     req.GetData(),
	}

	// Validate the broker image
	if ok, err := brokerImageInput.IsValid(); !ok {
		zap.L().Warn("BrokerImage is not valid", zap.Error(err))
		return &brokerpb.CreateBrokerImageResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Verify broker has no image
	ok, err := repositories.R().B().HasImage(brokerID)
	if err != nil {
		zap.L().Error("Cannot verify broker image existence", zap.Error(err))
		return &brokerpb.CreateBrokerImageResponse{}, status.Error(codes.Internal, err.Error())
	}
	if ok {
		zap.L().Warn("Broker already has an image", zap.String("broker_id", brokerID.String()))
		return &brokerpb.CreateBrokerImageResponse{}, status.Error(codes.InvalidArgument, "Broker already has an image")
	}

	// Create the broker image
	err = repositories.R().I().Create(brokerImageInput)
	if err != nil {
		zap.L().Error("PostBrokerImage.Create", zap.Error(err))
		return &brokerpb.CreateBrokerImageResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get the broker image back from the database
	brokerImage, found, err := repositories.R().I().Get(brokerImageInput.ID)
	if err != nil {
		zap.L().Error("Cannot get broker image", zap.String("uuid", brokerImageInput.ID.String()), zap.Error(err))
		return &brokerpb.CreateBrokerImageResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("Broker image not found after create", zap.String("uuid", brokerImageInput.ID.String()))
		return &brokerpb.CreateBrokerImageResponse{}, status.Error(codes.Internal, "Broker image not found after create")
	}

	// Set the broker image
	err = repositories.R().B().SetImage(brokerID, brokerImage.ID)
	if err != nil {
		zap.L().Error("SetImageBroker.Create", zap.Error(err))
		return &brokerpb.CreateBrokerImageResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &brokerpb.CreateBrokerImageResponse{
		Image: mappers.BrokerImageToProto(brokerImage),
	}, nil
}

// GetBrokerImage implements the GetBrokerImage RPC method.
func (h *Service) GetBrokerImage(ctx context.Context, req *brokerpb.GetBrokerImageRequest) (*brokerpb.GetBrokerImageResponse, error) {
	// Parse the image ID from the request
	imageID, err := uuid.Parse(req.GetImageId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid image ID", zap.String("image_id", req.GetImageId()), zap.Error(err))
		return &brokerpb.GetBrokerImageResponse{}, status.Error(codes.InvalidArgument, "Invalid image ID")
	}

	// Get the broker image
	brokerImage, found, err := repositories.R().I().Get(imageID)
	if err != nil {
		zap.L().Error("Cannot get brokerImage", zap.String("image_id", imageID.String()), zap.Error(err))
		return &brokerpb.GetBrokerImageResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Warn("BrokerImage does not exist", zap.String("image_id", imageID.String()))
		return &brokerpb.GetBrokerImageResponse{}, status.Error(codes.NotFound, "BrokerImage does not exist")
	}

	return &brokerpb.GetBrokerImageResponse{
		Name: brokerImage.Name,
		Data: brokerImage.Data,
	}, nil
}

// UpdateBrokerImage implements the UpdateBrokerImage RPC method.
func (h *Service) UpdateBrokerImage(ctx context.Context, req *brokerpb.UpdateBrokerImageRequest) (*brokerpb.UpdateBrokerImageResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.brokers.update")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return nil, err
	}

	// Parse the image ID from the request
	imageID, err := uuid.Parse(req.GetImageId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid image ID", zap.String("image_id", req.GetImageId()), zap.Error(err))
		return &brokerpb.UpdateBrokerImageResponse{}, status.Error(codes.InvalidArgument, "Invalid image ID")
	}

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetBrokerId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetBrokerId()), zap.Error(err))
		return &brokerpb.UpdateBrokerImageResponse{}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Construct the BrokerImageInput object from the request
	brokerImageInput := models.BrokerImage{
		ID:       imageID,
		BrokerID: brokerID,
		Name:     req.GetName(),
		Data:     req.GetData(),
	}

	// Validate the broker image
	if ok, err := brokerImageInput.IsValid(); !ok {
		zap.L().Warn("BrokerImage is not valid", zap.Error(err))
		return &brokerpb.UpdateBrokerImageResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Verify imageBroker existence
	exists, err := repositories.R().I().Exists(brokerID, imageID)
	if err != nil {
		zap.L().Error("Check imageBroker exists", zap.Error(err))
		return &brokerpb.UpdateBrokerImageResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		zap.L().Warn("ImageBroker not found", zap.String("broker_id", brokerID.String()), zap.String("image_id", imageID.String()))
		return &brokerpb.UpdateBrokerImageResponse{}, status.Error(codes.NotFound, "ImageBroker not found")
	}

	// Update the broker image
	err = repositories.R().I().Update(brokerImageInput)
	if err != nil {
		zap.L().Error("UpdateBrokerImage.Update", zap.Error(err))
		return &brokerpb.UpdateBrokerImageResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get the broker image back from the database
	brokerImage, found, err := repositories.R().I().Get(brokerImageInput.ID)
	if err != nil {
		zap.L().Error("Cannot get broker image", zap.String("uuid", brokerImageInput.ID.String()), zap.Error(err))
		return &brokerpb.UpdateBrokerImageResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("Broker image not found after create", zap.String("uuid", brokerImageInput.ID.String()))
		return &brokerpb.UpdateBrokerImageResponse{}, status.Error(codes.Internal, "Broker image not found after create")
	}

	return &brokerpb.UpdateBrokerImageResponse{
		Image: mappers.BrokerImageToProto(brokerImage),
	}, nil
}

// DeleteBrokerImage implements the DeleteBrokerImage RPC method.
func (h *Service) DeleteBrokerImage(ctx context.Context, req *brokerpb.DeleteBrokerImageRequest) (*brokerpb.DeleteBrokerImageResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.brokers.delete")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return nil, err
	}

	// Parse the image ID from the request
	imageID, err := uuid.Parse(req.GetImageId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid image ID", zap.String("image_id", req.GetImageId()), zap.Error(err))
		return &brokerpb.DeleteBrokerImageResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, "Invalid image ID")
	}

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetBrokerId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetBrokerId()), zap.Error(err))
		return &brokerpb.DeleteBrokerImageResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Verify imageBroker existence
	exists, err := repositories.R().I().Exists(brokerID, imageID)
	if err != nil {
		zap.L().Error("Check imageBroker exists", zap.Error(err))
		return &brokerpb.DeleteBrokerImageResponse{
			Success: false,
		}, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		zap.L().Warn("ImageBroker not found", zap.String("broker_id", brokerID.String()), zap.String("image_id", imageID.String()))
		return &brokerpb.DeleteBrokerImageResponse{
			Success: false,
		}, status.Error(codes.NotFound, "ImageBroker not found")
	}

	// Delete the broker image
	err = repositories.R().I().Delete(imageID)
	if err != nil {
		zap.L().Error("DeleteBrokerImage.Delete", zap.Error(err))
		return &brokerpb.DeleteBrokerImageResponse{
			Success: false,
		}, status.Error(codes.Internal, err.Error())
	}

	return &brokerpb.DeleteBrokerImageResponse{
		Success: true,
	}, nil
}
