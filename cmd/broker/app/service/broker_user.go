package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/brokerpb"
	"github.com/Zapharaos/fihub-backend/internal/mappers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateBrokerUser implements the CreateBrokerUser RPC method.
func (h *Service) CreateBrokerUser(ctx context.Context, req *brokerpb.CreateBrokerUserRequest) (*brokerpb.CreateBrokerUserResponse, error) {

	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &brokerpb.CreateBrokerUserResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetBrokerId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetUserId()), zap.Error(err))
		return &brokerpb.CreateBrokerUserResponse{}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Construct the BrokerUserInput object from the request
	brokerUserInput := models.BrokerUserInput{
		UserID:   userID.String(),
		BrokerID: brokerID.String(),
	}

	// Validate broker user
	/*if ok, err := brokerUserInput.IsValid(); !ok {
		zap.L().Warn("BrokerUser validation failed", zap.Error(err))
		return &brokerpb.CreateBrokerUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}*/

	// Convert to BrokerUser
	userBroker := brokerUserInput.ToUser()

	// Verify broker existence
	broker, exists, err := repositories.R().B().Get(userBroker.Broker.ID)
	if err != nil {
		zap.L().Error("Cannot get broker", zap.Error(err))
		return &brokerpb.CreateBrokerUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		zap.L().Warn("Broker not found",
			zap.String("UserID", userBroker.UserID.String()),
			zap.String("BrokerID", userBroker.Broker.ID.String()))
		return &brokerpb.CreateBrokerUserResponse{}, status.Error(codes.NotFound, "Broker not found")
	}

	// Check if broker is disabled
	if broker.Disabled {
		zap.L().Warn("Broker is disabled", zap.String("BrokerID", userBroker.Broker.ID.String()))
		return &brokerpb.CreateBrokerUserResponse{}, status.Error(codes.FailedPrecondition, "Broker disabled")
	}

	// Verify BrokerUser existence
	exists, err = repositories.R().U().Exists(userBroker)
	if err != nil {
		zap.L().Error("Verify BrokerUser existence", zap.Error(err))
		return &brokerpb.CreateBrokerUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	if exists {
		zap.L().Warn("BrokerUser already exists", zap.String("BrokerID", userBroker.Broker.ID.String()))
		return &brokerpb.CreateBrokerUserResponse{}, status.Error(codes.AlreadyExists, "Broker not found")
	}

	// Create BrokerUser
	err = repositories.R().U().Create(userBroker)
	if err != nil {
		zap.L().Error("Create BrokerUser", zap.Error(err))
		return &brokerpb.CreateBrokerUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get BrokerUsers back from database
	userBrokers, err := repositories.R().U().GetAll(userID)
	if err != nil {
		zap.L().Error("Cannot get UserBrokers", zap.String("uuid", userID.String()), zap.Error(err))
		return &brokerpb.CreateBrokerUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	if len(userBrokers) == 0 {
		zap.L().Error("UserBrokers not found", zap.String("uuid", userID.String()))
		return &brokerpb.CreateBrokerUserResponse{}, status.Error(codes.Internal, "UserBrokers not found")
	}

	// Convert userBrokers to gRPC format
	return &brokerpb.CreateBrokerUserResponse{
		UserBrokers: mappers.BrokerUsersToProto(userBrokers),
	}, nil
}

// GetBrokerUser implements the GetBrokerUser RPC method.
func (h *Service) GetBrokerUser(ctx context.Context, req *brokerpb.GetBrokerUserRequest) (*brokerpb.GetBrokerUserResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &brokerpb.GetBrokerUserResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetBrokerId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetUserId()), zap.Error(err))
		return &brokerpb.GetBrokerUserResponse{}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Retrieve the BrokerUser from the database
	brokerUser, exists, err := repositories.R().U().Get(models.BrokerUser{
		UserID: userID,
		Broker: models.Broker{ID: brokerID},
	})
	if err != nil {
		zap.L().Error("Cannot get user broker", zap.Error(err))
		return &brokerpb.GetBrokerUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		zap.L().Warn("BrokerUser not found",
			zap.String("UserID", userID.String()),
			zap.String("BrokerID", brokerID.String()))
		return &brokerpb.GetBrokerUserResponse{}, status.Error(codes.NotFound, "Broker not found")
	}

	return &brokerpb.GetBrokerUserResponse{
		BrokerUser: mappers.BrokerUserToProto(brokerUser),
	}, nil
}

// DeleteBrokerUser implements the DeleteBrokerUser RPC method.
func (h *Service) DeleteBrokerUser(ctx context.Context, req *brokerpb.DeleteBrokerUserRequest) (*brokerpb.DeleteBrokerUserResponse, error) {

	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &brokerpb.DeleteBrokerUserResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse the broker ID from the request
	brokerID, err := uuid.Parse(req.GetBrokerId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid broker ID", zap.String("broker_id", req.GetUserId()), zap.Error(err))
		return &brokerpb.DeleteBrokerUserResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, "Invalid broker ID")
	}

	// Construct the BrokerUser object from the request
	brokerUser := models.BrokerUser{
		UserID: userID,
		Broker: models.Broker{ID: brokerID},
	}

	// Verify BrokerUser existence
	exists, err := repositories.R().U().Exists(brokerUser)
	if err != nil {
		zap.L().Error("Check userBroker exists", zap.Error(err))
		return &brokerpb.DeleteBrokerUserResponse{
			Success: false,
		}, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		zap.L().Warn("BrokerUser not found",
			zap.String("UserID", userID.String()),
			zap.String("BrokerID", brokerID.String()))
		return &brokerpb.DeleteBrokerUserResponse{
			Success: false,
		}, status.Error(codes.NotFound, "BrokerUser not found")
	}

	// Remove BrokerUser
	err = repositories.R().U().Delete(brokerUser)
	if err != nil {
		zap.L().Error("Cannot remove user broker", zap.String("uuid", userID.String()), zap.Error(err))
		return &brokerpb.DeleteBrokerUserResponse{
			Success: false,
		}, status.Error(codes.Internal, err.Error())
	}

	return &brokerpb.DeleteBrokerUserResponse{
		Success: true,
	}, nil
}

// ListUserBrokers implements the ListUserBrokers RPC method.
func (h *Service) ListUserBrokers(ctx context.Context, req *brokerpb.ListUserBrokersRequest) (*brokerpb.ListUserBrokersResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &brokerpb.ListUserBrokersResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Get userBrokers back from database
	userBrokers, err := repositories.R().U().GetAll(userID)
	if err != nil {
		zap.L().Error("Cannot get user brokers", zap.String("uuid", userID.String()), zap.Error(err))
		return &brokerpb.ListUserBrokersResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Convert userBrokers to gRPC format
	return &brokerpb.ListUserBrokersResponse{
		UserBrokers: mappers.BrokerUsersToProto(userBrokers),
	}, nil
}
