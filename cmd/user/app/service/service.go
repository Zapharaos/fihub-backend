package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/mappers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service is the implementation of the UserService interface.
type Service struct {
	userpb.UnimplementedUserServiceServer
}

// CreateUser implements the CreateUser RPC method.
func (s *Service) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	// Construct the user input object
	userInputCreate := models.UserInputCreate{
		UserInputPassword: models.UserInputPassword{
			UserWithPassword: models.UserWithPassword{
				User: models.User{
					ID:    uuid.New(),
					Email: req.GetEmail(),
				},
				Password: req.GetPassword(),
			},
			Confirmation: req.GetConfirmation(),
		},
		Checkbox: req.GetCheckbox(),
	}

	// Validate user
	if ok, err := userInputCreate.IsValid(); !ok {
		zap.L().Warn("User is not valid", zap.Error(err))
		return &userpb.CreateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert to UserWithPassword
	userWithPassword := userInputCreate.UserWithPassword

	// Verify user existence
	exists, err := repositories.R().Exists(userWithPassword.Email)
	if err != nil {
		zap.L().Error("Check user exists", zap.Error(err))
		return &userpb.CreateUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	if exists {
		zap.L().Warn("User already exists", zap.String("email", userWithPassword.Email))
		return &userpb.CreateUserResponse{}, status.Error(codes.AlreadyExists, "email-used")
	}

	// Create user
	userID, err := repositories.R().Create(userWithPassword)
	if err != nil {
		zap.L().Error("PostUser.Create", zap.Error(err))
		return &userpb.CreateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get user back from database
	user, found, err := repositories.R().Get(userID)
	if err != nil {
		zap.L().Error("Cannot get user", zap.String("uuid", userID.String()), zap.Error(err))
		return &userpb.CreateUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("User not found after creation", zap.String("uuid", userID.String()))
		return &userpb.CreateUserResponse{}, status.Error(codes.Internal, "User not found after creation")
	}

	return &userpb.CreateUserResponse{
		User: mappers.UserToProto(user),
	}, nil
}

// GetUser implements the GetUser RPC method.
func (s *Service) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &userpb.GetUserResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Check user permissions
	err = security.Facade().CheckPermission(ctx, "admin.users.read", userID)
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &userpb.GetUserResponse{}, err
	}

	// the user default accessible data
	user, found, err := repositories.R().Get(userID)
	if err != nil {
		zap.L().Error("GetUser.Get", zap.Error(err))
		return &userpb.GetUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Warn("User not found", zap.String("uuid", userID.String()))
		return &userpb.GetUserResponse{}, status.Error(codes.NotFound, "User not found")
	}

	return &userpb.GetUserResponse{
		User: mappers.UserToProto(user),
	}, nil
}

// GetByEmail implements the GetByEmail RPC method.
func (s *Service) GetByEmail(ctx context.Context, req *userpb.GetByEmailRequest) (*userpb.GetByEmailResponse, error) {
	user, found, err := repositories.R().GetByEmail(req.GetEmail())
	if err != nil {
		zap.L().Error("GetByEmail.GetByEmail", zap.Error(err))
		return &userpb.GetByEmailResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Warn("User not found", zap.String("email", req.GetEmail()))
		return &userpb.GetByEmailResponse{}, status.Error(codes.NotFound, "User not found")
	}

	return &userpb.GetByEmailResponse{
		User: mappers.UserToProto(user),
	}, nil
}

// UpdateUser implements the UpdateUser RPC method.
func (s *Service) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &userpb.UpdateUserResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	err = security.Facade().CheckPermission(ctx, "admin.users.update", userID)
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &userpb.UpdateUserResponse{}, err
	}

	// Construct the user object
	user := models.User{
		ID:    userID,
		Email: req.GetEmail(),
	}

	// Validate user
	if ok, err := user.IsValid(); !ok {
		zap.L().Warn("User is not valid", zap.Error(err))
		return &userpb.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Update user
	err = repositories.R().Update(user)
	if err != nil {
		zap.L().Error("PutUser.Update", zap.Error(err))
		return &userpb.UpdateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get user back from database
	user, found, err := repositories.R().Get(userID)
	if err != nil {
		zap.L().Error("Cannot get user", zap.String("uuid", userID.String()), zap.Error(err))
		return &userpb.UpdateUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("User not found after update", zap.String("uuid", userID.String()))
		return &userpb.UpdateUserResponse{}, status.Error(codes.Internal, "User not found after update")
	}

	return &userpb.UpdateUserResponse{
		User: mappers.UserToProto(user),
	}, nil
}

// UpdateUserPassword implements the UpdateUserPassword RPC method.
func (s *Service) UpdateUserPassword(ctx context.Context, req *userpb.UpdateUserPasswordRequest) (*userpb.UpdateUserPasswordResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &userpb.UpdateUserPasswordResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	err = security.Facade().CheckPermission(ctx, "admin.users.update", userID)
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &userpb.UpdateUserPasswordResponse{}, err
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
		return &userpb.UpdateUserPasswordResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Update password
	err = repositories.R().UpdateWithPassword(userInputPassword.UserWithPassword)
	if err != nil {
		zap.L().Error("PutUser.UpdatePassword", zap.Error(err))
		return &userpb.UpdateUserPasswordResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &userpb.UpdateUserPasswordResponse{
		Success: true,
	}, nil
}

// DeleteUser implements the DeleteUser RPC method.
func (s *Service) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &userpb.DeleteUserResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	err = security.Facade().CheckPermission(ctx, "admin.users.delete", userID)
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &userpb.DeleteUserResponse{}, err
	}

	err = repositories.R().Delete(userID)
	if err != nil {
		zap.L().Error("DeleteUser.Delete", zap.Error(err))
		return &userpb.DeleteUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &userpb.DeleteUserResponse{
		Success: true,
	}, nil
}

// ListUsers implements the ListUsers RPC method.
func (s *Service) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.users.list")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &userpb.ListUsersResponse{}, err
	}

	// List users
	users, err := repositories.R().List()
	if err != nil {
		zap.L().Error("GetUsers.List", zap.Error(err))
		return &userpb.ListUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &userpb.ListUsersResponse{
		Users: mappers.UsersToProto(users),
	}, nil
}

// AuthenticateUser implements the AuthenticateUser RPC method.
func (s *Service) AuthenticateUser(ctx context.Context, req *userpb.AuthenticateUserRequest) (*userpb.AuthenticateUserResponse, error) {
	// Try to authenticate the user
	user, found, err := repositories.R().Authenticate(req.Email, req.Password)
	if err != nil || !found {
		zap.L().Error("AuthenticateUser", zap.Error(err))
		return &userpb.AuthenticateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &userpb.AuthenticateUserResponse{
		User: mappers.UserToProto(user),
	}, nil
}
