package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/internal/mappers"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddUsersToRole implements the AddUsersToRole RPC method.
func (s *Service) AddUsersToRole(ctx context.Context, req *securitypb.AddUsersToRoleRequest) (*securitypb.AddUsersToRoleResponse, error) {
	// Check user permissions for creating a role
	err := security.Facade().CheckPermission(ctx, "admin.roles.users.update")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.AddUsersToRoleResponse{}, err
	}

	// Parse the user ID from the request
	roleID, err := uuid.Parse(req.GetRoleId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetRoleId()), zap.Error(err))
		return &securitypb.AddUsersToRoleResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse user UUIDs
	uuidUsers, err := utils.StringsToUUIDs(req.GetUserIds())
	if err != nil {
		zap.L().Error("Invalid user UUIDs", zap.Error(err))
		return &securitypb.AddUsersToRoleResponse{}, status.Error(codes.InvalidArgument, "Invalid user UUIDs")
	}

	err = repositories.R().R().AddToUsers(uuidUsers, roleID)
	if err != nil {
		zap.L().Error("AddUsersToRole", zap.Error(err))
		return &securitypb.AddUsersToRoleResponse{}, status.Error(codes.Internal, "Failed to add users to role")
	}

	// Retrieve users by role ID from database
	users, err := repositories.R().R().ListUsersByRoleId(roleID)
	if err != nil {
		zap.L().Error("ListUsersForRole", zap.Error(err))
		return &securitypb.AddUsersToRoleResponse{}, status.Error(codes.Internal, "Failed to list users for role")
	}

	return &securitypb.AddUsersToRoleResponse{
		UserIds: users,
	}, nil
}

// RemoveUsersFromRole implements the RemoveUsersFromRole RPC method.
func (s *Service) RemoveUsersFromRole(ctx context.Context, req *securitypb.RemoveUsersFromRoleRequest) (*securitypb.RemoveUsersFromRoleResponse, error) {
	// Check user permissions for creating a role
	err := security.Facade().CheckPermission(ctx, "admin.roles.users.delete")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.RemoveUsersFromRoleResponse{}, err
	}

	// Parse the user ID from the request
	roleID, err := uuid.Parse(req.GetRoleId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetRoleId()), zap.Error(err))
		return &securitypb.RemoveUsersFromRoleResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse user UUIDs
	uuidUsers, err := utils.StringsToUUIDs(req.GetUserIds())
	if err != nil {
		zap.L().Error("Invalid user UUIDs", zap.Error(err))
		return &securitypb.RemoveUsersFromRoleResponse{}, status.Error(codes.InvalidArgument, "Invalid user UUIDs")
	}

	err = repositories.R().R().RemoveFromUsers(uuidUsers, roleID)
	if err != nil {
		zap.L().Error("RemoveUsersFromRole", zap.Error(err))
		return &securitypb.RemoveUsersFromRoleResponse{}, status.Error(codes.Internal, "Failed to remove users from role")
	}

	// Retrieve users by role ID from database
	users, err := repositories.R().R().ListUsersByRoleId(roleID)
	if err != nil {
		zap.L().Error("ListUsersForRole", zap.Error(err))
		return &securitypb.RemoveUsersFromRoleResponse{}, status.Error(codes.Internal, "Failed to list users for role")
	}

	return &securitypb.RemoveUsersFromRoleResponse{
		UserIds: users,
	}, nil
}

// ListUsersForRole implements the ListUsersForRole RPC method.
func (s *Service) ListUsersForRole(ctx context.Context, req *securitypb.ListUsersForRoleRequest) (*securitypb.ListUsersForRoleResponse, error) {
	// Check user permissions for creating a role
	err := security.Facade().CheckPermission(ctx, "admin.roles.users.list")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.ListUsersForRoleResponse{}, err
	}

	// Parse the user ID from the request
	roleID, err := uuid.Parse(req.GetRoleId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetRoleId()), zap.Error(err))
		return &securitypb.ListUsersForRoleResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Retrieve users by role ID from database
	users, err := repositories.R().R().ListUsersByRoleId(roleID)
	if err != nil {
		zap.L().Error("ListUsersForRole", zap.Error(err))
		return &securitypb.ListUsersForRoleResponse{}, status.Error(codes.Internal, "Failed to list users for role")
	}

	return &securitypb.ListUsersForRoleResponse{
		UserIds: users,
	}, nil
}

// SetRolesForUser implements the SetRolesForUser RPC method.
func (s *Service) SetRolesForUser(ctx context.Context, req *securitypb.SetRolesForUserRequest) (*securitypb.SetRolesForUserResponse, error) {
	// Check user permissions for creating a role
	err := security.Facade().CheckPermission(ctx, "admin.users.roles.update")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.SetRolesForUserResponse{}, err
	}

	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &securitypb.SetRolesForUserResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Parse role UUIDs
	uuidRoles, err := utils.StringsToUUIDs(req.GetRoleIds())
	if err != nil {
		zap.L().Error("Invalid role UUIDs", zap.Error(err))
		return &securitypb.SetRolesForUserResponse{}, status.Error(codes.InvalidArgument, "Invalid role UUIDs")
	}

	// Set roles on user
	err = repositories.R().R().SetForUser(userID, uuidRoles)
	if err != nil {
		zap.L().Error("SetRolesForUser", zap.Error(err))
		return &securitypb.SetRolesForUserResponse{}, status.Error(codes.Internal, "Failed to set user roles")
	}

	// Get all roles with permissions for user from the database
	roles, err := repositories.R().R().ListWithPermissionsByUserId(userID)
	if err != nil {
		zap.L().Error("Cannot list roles with permissions", zap.Error(err))
		return &securitypb.SetRolesForUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &securitypb.SetRolesForUserResponse{
		Roles: mappers.RolesWithPermissionsToProto(roles),
	}, nil
}

// ListRolesForUser implements the ListRolesForUser RPC method.
func (s *Service) ListRolesForUser(ctx context.Context, req *securitypb.ListRolesForUserRequest) (*securitypb.ListRolesForUserResponse, error) {
	// Check user permissions for creating a role
	err := security.Facade().CheckPermission(ctx, "admin.users.roles.list")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.ListRolesForUserResponse{}, err
	}

	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &securitypb.ListRolesForUserResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Get all roles for user from the database
	roles, err := repositories.R().R().ListByUserId(userID)
	if err != nil {
		zap.L().Error("Cannot list roles", zap.Error(err))
		return &securitypb.ListRolesForUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &securitypb.ListRolesForUserResponse{
		Roles: mappers.RolesToProto(roles),
	}, nil
}

// ListRolesWithPermissionsForUser implements the ListRolesWithPermissionsForUser RPC method.
func (s *Service) ListRolesWithPermissionsForUser(ctx context.Context, req *securitypb.ListRolesWithPermissionsForUserRequest) (*securitypb.ListRolesWithPermissionsForUserResponse, error) {
	// Check user permissions for creating a role
	err := security.Facade().CheckPermission(ctx, "admin.users.roles.list")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.ListRolesWithPermissionsForUserResponse{}, err
	}

	// Parse the user ID from the request
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return &securitypb.ListRolesWithPermissionsForUserResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Get all roles with permissions for user from the database
	roles, err := repositories.R().R().ListWithPermissionsByUserId(userID)
	if err != nil {
		zap.L().Error("Cannot list roles with permissions", zap.Error(err))
		return &securitypb.ListRolesWithPermissionsForUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &securitypb.ListRolesWithPermissionsForUserResponse{
		Roles: mappers.RolesWithPermissionsToProto(roles),
	}, nil
}

// ListUsersFull implements the ListUsersFull RPC method.
func (s *Service) ListUsersFull(ctx context.Context, req *securitypb.ListUsersFullRequest) (*securitypb.ListUsersFullResponse, error) {
	// Check user permissions for creating a role
	err := security.Facade().CheckPermission(ctx, "admin.users.list")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.ListUsersFullResponse{}, err
	}

	// Retrieve users from database
	userIds, err := repositories.R().R().ListUsers()
	if err != nil {
		zap.L().Error("ListUsers", zap.Error(err))
		return &securitypb.ListUsersFullResponse{}, status.Error(codes.Internal, "Failed to list users")
	}

	users := make([]*securitypb.UserWithRoles, 0, len(userIds))
	for _, userID := range userIds {
		roles, err := repositories.R().R().ListByUserId(uuid.MustParse(userID))
		if err != nil {
			zap.L().Error("Cannot list roles", zap.Error(err))
			return &securitypb.ListUsersFullResponse{}, status.Error(codes.Internal, err.Error())
		}

		users = append(users, &securitypb.UserWithRoles{
			UserId: userID,
			Roles:  mappers.RolesToProto(roles),
		})
	}

	return &securitypb.ListUsersFullResponse{
		Users: users,
	}, nil
}
