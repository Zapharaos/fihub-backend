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

// CreateBroker implements the CreateBroker RPC method.
func (h *Service) CreateBroker(ctx context.Context, req *brokerpb.CreateBrokerRequest) (*brokerpb.CreateBrokerResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.brokers.create")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &brokerpb.CreateBrokerResponse{}, err
	}

	// Construct the Broker object from the request
	broker := models.Broker{
		ID:       uuid.New(),
		Name:     req.GetName(),
		Disabled: req.GetDisabled(),
	}

	// Validate the broker
	if valid, err := broker.IsValid(); !valid {
		zap.L().Warn("Broker is not valid", zap.Error(err))
		return &brokerpb.CreateBrokerResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Verify that the broker does not already exist
	exists, err := repositories.R().B().ExistsByName(broker.Name)
	if err != nil {
		zap.L().Error("Check broker exists", zap.Error(err))
		return &brokerpb.CreateBrokerResponse{}, status.Error(codes.Internal, err.Error())
	}
	if exists {
		zap.L().Warn("Broker already exists", zap.String("Name", broker.Name))
		return &brokerpb.CreateBrokerResponse{}, status.Error(codes.AlreadyExists, "name-used")
	}

	// Create the broker
	brokerID, err := repositories.R().B().Create(broker)
	if err != nil {
		zap.L().Warn("Create broker", zap.Error(err))
		return &brokerpb.CreateBrokerResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get the broker from the database
	broker, found, err := repositories.R().B().Get(brokerID)
	if err != nil {
		zap.L().Error("Cannot get broker", zap.String("uuid", brokerID.String()), zap.Error(err))
		return &brokerpb.CreateBrokerResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("Broker not found after creation", zap.String("uuid", brokerID.String()))
		return &brokerpb.CreateBrokerResponse{}, status.Error(codes.Internal, "Broker not found after creation")
	}

	return &brokerpb.CreateBrokerResponse{
		Broker: mappers.BrokerToProto(broker),
	}, nil
}

// GetBroker implements the GetBroker RPC method.
func (h *Service) GetBroker(ctx context.Context, req *brokerpb.GetBrokerRequest) (*brokerpb.GetBrokerResponse, error) {

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetId()), zap.Error(err))
		return &brokerpb.GetBrokerResponse{}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Get the broker from the database
	broker, found, err := repositories.R().B().Get(brokerID)
	if err != nil {
		zap.L().Error("Cannot get broker", zap.String("uuid", brokerID.String()), zap.Error(err))
		return &brokerpb.GetBrokerResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("Broker not found", zap.String("uuid", brokerID.String()))
		return &brokerpb.GetBrokerResponse{}, status.Error(codes.NotFound, "Broker not found")
	}

	return &brokerpb.GetBrokerResponse{
		Broker: mappers.BrokerToProto(broker),
	}, nil
}

// UpdateBroker implements the UpdateBroker RPC method.
func (h *Service) UpdateBroker(ctx context.Context, req *brokerpb.UpdateBrokerRequest) (*brokerpb.UpdateBrokerResponse, error) {

	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.brokers.update")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &brokerpb.UpdateBrokerResponse{}, err
	}

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetId()), zap.Error(err))
		return &brokerpb.UpdateBrokerResponse{}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Construct the Broker object from the request
	broker := models.Broker{
		ID:       brokerID,
		Name:     req.GetName(),
		Disabled: req.GetDisabled(),
	}

	// Validate the broker
	if valid, err := broker.IsValid(); !valid {
		zap.L().Warn("Broker is not valid", zap.Error(err))
		return &brokerpb.UpdateBrokerResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Retrieve the broker from the database and verify its existence
	oldBroker, found, err := repositories.R().B().Get(brokerID)
	if err != nil {
		zap.L().Error("Cannot get broker", zap.String("uuid", brokerID.String()), zap.Error(err))
		return &brokerpb.UpdateBrokerResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("Broker not found", zap.String("uuid", brokerID.String()))
		return &brokerpb.UpdateBrokerResponse{}, status.Error(codes.NotFound, "Broker not found")
	}

	// Check if the broker name has changed
	if oldBroker.Name != broker.Name {
		// Verify that the broker name is not already used
		exists, err := repositories.R().B().ExistsByName(broker.Name)
		if err != nil {
			zap.L().Error("Check broker name exists", zap.Error(err))
			return &brokerpb.UpdateBrokerResponse{}, status.Error(codes.Internal, err.Error())
		}
		if exists {
			zap.L().Warn("Broker name already used", zap.String("Name", broker.Name))
			return &brokerpb.UpdateBrokerResponse{}, status.Error(codes.AlreadyExists, "name-used")
		}
	}

	// Update the broker
	err = repositories.R().B().Update(broker)
	if err != nil {
		zap.L().Warn("Update broker", zap.Error(err))
		return &brokerpb.UpdateBrokerResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get the broker from the database
	broker, found, err = repositories.R().B().Get(brokerID)
	if err != nil {
		zap.L().Error("Cannot get broker", zap.String("uuid", brokerID.String()), zap.Error(err))
		return &brokerpb.UpdateBrokerResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("Broker not found after update", zap.String("uuid", brokerID.String()))
		return &brokerpb.UpdateBrokerResponse{}, status.Error(codes.Internal, "Broker not found after update")
	}

	return &brokerpb.UpdateBrokerResponse{
		Broker: mappers.BrokerToProto(broker),
	}, nil
}

// DeleteBroker implements the DeleteBroker RPC method.
func (h *Service) DeleteBroker(ctx context.Context, req *brokerpb.DeleteBrokerRequest) (*brokerpb.DeleteBrokerResponse, error) {

	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.brokers.delete")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &brokerpb.DeleteBrokerResponse{}, err
	}

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetId()), zap.Error(err))
		return &brokerpb.DeleteBrokerResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Verify that the broker exists
	exists, err := repositories.R().B().Exists(brokerID)
	if err != nil {
		zap.L().Error("Check broker exists", zap.Error(err))
		return &brokerpb.DeleteBrokerResponse{
			Success: false,
		}, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		zap.L().Warn("Broker not found", zap.String("BrokerID", brokerID.String()))
		return &brokerpb.DeleteBrokerResponse{
			Success: false,
		}, status.Error(codes.NotFound, "Broker not found")
	}

	// Delete the broker
	err = repositories.R().B().Delete(brokerID)
	if err != nil {
		zap.L().Warn("Delete broker", zap.Error(err))
		return &brokerpb.DeleteBrokerResponse{
			Success: false,
		}, status.Error(codes.Internal, err.Error())
	}

	return &brokerpb.DeleteBrokerResponse{
		Success: true,
	}, nil
}

// ListBrokers implements the ListBrokers RPC method.
func (h *Service) ListBrokers(ctx context.Context, req *brokerpb.ListBrokersRequest) (*brokerpb.ListBrokersResponse, error) {

	var (
		result []models.Broker
		err    error
	)

	// Check if the query parameter is set to true
	if req.GetEnabledOnly() {
		result, err = repositories.R().B().GetAllEnabled()
	} else {
		result, err = repositories.R().B().GetAll()
	}

	if err != nil {
		return &brokerpb.ListBrokersResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Convert userBrokers to gRPC format
	return &brokerpb.ListBrokersResponse{
		Brokers: mappers.BrokersToProto(result),
	}, nil
}
