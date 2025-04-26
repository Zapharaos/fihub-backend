package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service is the implementation of the UserService interface.
type Service struct {
	protogen.UnimplementedUserServiceServer
}

// CreateUser implements the CreateUser RPC method.
func (s *Service) CreateUser(ctx context.Context, req *protogen.CreateUserRequest) (*protogen.CreateUserResponse, error) {
	// Construct the user input object
	userInputCreate := models.UserInputPassword{
		UserWithPassword: models.UserWithPassword{
			User: models.User{
				ID:    uuid.New(),
				Email: req.GetEmail(),
			},
			Password: req.GetPassword(),
		},
		Confirmation: req.GetConfirmation(),
	}

	// Validate user
	if ok, err := userInputCreate.IsValid(); !ok {
		zap.L().Warn("User is not valid", zap.Error(err))
		return &protogen.CreateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert to UserWithPassword
	userWithPassword := userInputCreate.UserWithPassword

	// Verify user existence
	exists, err := repositories.R().Exists(userWithPassword.Email)
	if err != nil {
		zap.L().Error("Check user exists", zap.Error(err))
		return &protogen.CreateUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	if exists {
		zap.L().Warn("User already exists", zap.String("email", userWithPassword.Email))
		return &protogen.CreateUserResponse{}, status.Error(codes.AlreadyExists, "email-used")
	}

	// Create user
	userID, err := repositories.R().Create(userWithPassword)
	if err != nil {
		zap.L().Error("PostUser.Create", zap.Error(err))
		return &protogen.CreateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get user back from database
	user, found, err := repositories.R().Get(userID)
	if err != nil {
		zap.L().Error("Cannot get user", zap.String("uuid", userID.String()), zap.Error(err))
		return &protogen.CreateUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("User not found after creation", zap.String("uuid", userID.String()))
		return &protogen.CreateUserResponse{}, status.Error(codes.Internal, "User not found after creation")
	}

	return &protogen.CreateUserResponse{
		User: user.ToProtogenUser(),
	}, nil
}

// GetUser implements the GetUser RPC method.
func (s *Service) GetUser(ctx context.Context, req *protogen.GetUserRequest) (*protogen.GetUserResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &protogen.GetUserResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// the user default accessible data
	user, found, err := repositories.R().Get(userID)
	if err != nil {
		zap.L().Error("GetUser.Get", zap.Error(err))
		return &protogen.GetUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Warn("User not found", zap.String("uuid", userID.String()))
		return &protogen.GetUserResponse{}, status.Error(codes.NotFound, "User not found")
	}

	return &protogen.GetUserResponse{
		User: user.ToProtogenUser(),
	}, nil
}

// UpdateUser implements the UpdateUser RPC method.
func (s *Service) UpdateUser(ctx context.Context, req *protogen.UpdateUserRequest) (*protogen.UpdateUserResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &protogen.UpdateUserResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Construct the user object
	user := models.User{
		ID:    userID,
		Email: req.GetEmail(),
	}

	// Validate user
	if ok, err := user.IsValid(); !ok {
		zap.L().Warn("User is not valid", zap.Error(err))
		return &protogen.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Update user
	err = repositories.R().Update(user)
	if err != nil {
		zap.L().Error("PutUser.Update", zap.Error(err))
		return &protogen.UpdateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get user back from database
	user, found, err := repositories.R().Get(userID)
	if err != nil {
		zap.L().Error("Cannot get user", zap.String("uuid", userID.String()), zap.Error(err))
		return &protogen.UpdateUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("User not found after update", zap.String("uuid", userID.String()))
		return &protogen.UpdateUserResponse{}, status.Error(codes.Internal, "User not found after update")
	}

	return &protogen.UpdateUserResponse{
		User: user.ToProtogenUser(),
	}, nil
}

// UpdateUserPassword implements the UpdateUserPassword RPC method.
func (s *Service) UpdateUserPassword(ctx context.Context, req *protogen.UpdateUserPasswordRequest) (*protogen.UpdateUserPasswordResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &protogen.UpdateUserPasswordResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Construct the user input object
	userInputPassword := models.UserInputPassword{
		UserWithPassword: models.UserWithPassword{
			User: models.User{
				ID: userID,
			},
			Password: req.GetPassword(),
		},
		Confirmation: req.GetConfirmation(),
	}

	// Validate password
	if ok, err := userInputPassword.IsValidPassword(); !ok {
		zap.L().Warn("User is not valid", zap.Error(err))
		return &protogen.UpdateUserPasswordResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Update password
	err = repositories.R().UpdateWithPassword(userInputPassword.UserWithPassword)
	if err != nil {
		zap.L().Error("PutUser.UpdatePassword", zap.Error(err))
		return &protogen.UpdateUserPasswordResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &protogen.UpdateUserPasswordResponse{
		Success: true,
	}, nil
}

// DeleteUser implements the DeleteUser RPC method.
func (s *Service) DeleteUser(ctx context.Context, req *protogen.DeleteUserRequest) (*protogen.DeleteUserResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &protogen.DeleteUserResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	err = repositories.R().Delete(userID)
	if err != nil {
		zap.L().Error("DeleteUser.Delete", zap.Error(err))
		return &protogen.DeleteUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &protogen.DeleteUserResponse{
		Success: true,
	}, nil
}
