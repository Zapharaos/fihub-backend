package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddUsersToRole implements the AddUsersToRole RPC method.
func (s *Service) AddUsersToRole(ctx context.Context, req *protogen.AddUsersToRoleRequest) (*protogen.AddUsersToRoleResponse, error) {
	// Check user permissions for creating a role
	err := security.Facade().CheckPermission(ctx, "admin.roles.users.update")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return nil, err
	}

	// Parse the user ID from the request
	roleID, err := uuid.Parse(req.GetRoleId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetRoleId()), zap.Error(err))
		return &protogen.AddUsersToRoleResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse user UUIDs
	uuidUsers, err := utils.StringsToUUIDs(req.GetUserIds())
	if err != nil {
		zap.L().Error("Invalid user UUIDs", zap.Error(err))
		return &protogen.AddUsersToRoleResponse{}, status.Error(codes.InvalidArgument, "Invalid user UUIDs")
	}

	err = repositories.R().R().AddToUsers(uuidUsers, roleID)
	if err != nil {
		zap.L().Error("AddUsersToRole", zap.Error(err))
		return &protogen.AddUsersToRoleResponse{}, status.Error(codes.Internal, "Failed to add users to role")
	}

	return &protogen.AddUsersToRoleResponse{
		Success: true,
	}, nil
}

// RemoveUsersFromRole implements the RemoveUsersFromRole RPC method.
func (s *Service) RemoveUsersFromRole(ctx context.Context, req *protogen.RemoveUsersFromRoleRequest) (*protogen.RemoveUsersFromRoleResponse, error) {
	// Check user permissions for creating a role
	err := security.Facade().CheckPermission(ctx, "admin.roles.users.delete")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return nil, err
	}

	// Parse the user ID from the request
	roleID, err := uuid.Parse(req.GetRoleId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetRoleId()), zap.Error(err))
		return &protogen.RemoveUsersFromRoleResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse user UUIDs
	uuidUsers, err := utils.StringsToUUIDs(req.GetUserIds())
	if err != nil {
		zap.L().Error("Invalid user UUIDs", zap.Error(err))
		return &protogen.RemoveUsersFromRoleResponse{}, status.Error(codes.InvalidArgument, "Invalid user UUIDs")
	}

	err = repositories.R().R().RemoveFromUsers(uuidUsers, roleID)
	if err != nil {
		zap.L().Error("RemoveUsersFromRole", zap.Error(err))
		return &protogen.RemoveUsersFromRoleResponse{}, status.Error(codes.Internal, "Failed to remove users from role")
	}

	return &protogen.RemoveUsersFromRoleResponse{
		Success: true,
	}, nil
}

// SetRolesForUser implements the SetRolesForUser RPC method.
func (s *Service) SetRolesForUser(ctx context.Context, req *protogen.SetRolesForUserRequest) (*protogen.SetRolesForUserResponse, error) {
	// Check user permissions for creating a role
	err := security.Facade().CheckPermission(ctx, "admin.users.roles.update")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return nil, err
	}

	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &protogen.SetRolesForUserResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse role UUIDs
	uuidRoles, err := utils.StringsToUUIDs(req.GetRoleIds())
	if err != nil {
		zap.L().Error("Invalid role UUIDs", zap.Error(err))
		return &protogen.SetRolesForUserResponse{}, status.Error(codes.InvalidArgument, "Invalid role UUIDs")
	}

	// Set roles on user
	err = repositories.R().R().SetForUser(userID, uuidRoles)
	if err != nil {
		zap.L().Error("SetRolesForUser", zap.Error(err))
		return &protogen.SetRolesForUserResponse{}, status.Error(codes.Internal, "Failed to set user roles")
	}

	return &protogen.SetRolesForUserResponse{
		Success: true,
	}, nil
}
