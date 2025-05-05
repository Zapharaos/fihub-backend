package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/internal/mappers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateRole implements the CreateRole RPC method.
func (s *Service) CreateRole(ctx context.Context, req *securitypb.CreateRoleRequest) (*securitypb.CreateRoleResponse, error) {
	// Check user permissions for creating a role
	err := security.Facade().CheckPermission(ctx, "admin.roles.create")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.CreateRoleResponse{}, err
	}

	// Construct the Role object from the request
	role := models.Role{
		Id:   uuid.New(),
		Name: req.GetName(),
	}

	if ok, err := role.IsValid(); !ok {
		zap.L().Warn("Role is not valid", zap.Error(err))
		return &securitypb.CreateRoleResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check user permissions for updating a role permissions
	err = security.Facade().CheckPermission(ctx, "admin.roles.permissions.update")

	// If the user has permission to update role permissions, validate and set them
	var permissions models.RolePermissionsInput
	if err == nil {

		// Construct the Permissions object from the request
		permissions = models.RolePermissionsInputFromUUIDs(req.GetPermissions())

		// Validate the permissions
		if ok, err := permissions.IsValid(); !ok {
			zap.L().Warn("Permissions are not valid", zap.Error(err))
			return &securitypb.CreateRoleResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	// Create the role in the database
	roleID, err := repositories.R().R().Create(role, permissions)
	if err != nil {
		zap.L().Error("Cannot create role", zap.Error(err))
		return &securitypb.CreateRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get the role + permissions from the database
	result, found, err := repositories.R().R().GetWithPermissions(roleID)
	if err != nil {
		zap.L().Error("Cannot get role", zap.String("uuid", roleID.String()), zap.Error(err))
		return &securitypb.CreateRoleResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("Role not found after creation", zap.String("uuid", roleID.String()))
		return &securitypb.CreateRoleResponse{}, status.Error(codes.Internal, "Role not found after creation")
	}

	// Convert the role to the gen format
	return &securitypb.CreateRoleResponse{
		Role: mappers.RoleWithPermissionsToProto(result),
	}, nil
}

// GetRole implements the GetRole RPC method.
func (s *Service) GetRole(ctx context.Context, req *securitypb.GetRoleRequest) (*securitypb.GetRoleResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.roles.read")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.GetRoleResponse{}, err
	}

	// Parse the user ID from the request
	roleID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &securitypb.GetRoleResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Get the role from the database
	role, found, err := repositories.R().R().GetWithPermissions(roleID)
	if err != nil {
		zap.L().Error("Cannot load role", zap.String("uuid", roleID.String()), zap.Error(err))
		return &securitypb.GetRoleResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Debug("Role not found", zap.String("uuid", roleID.String()))
		return &securitypb.GetRoleResponse{}, status.Error(codes.NotFound, "Role not found")
	}

	// Convert the role to the gen format
	return &securitypb.GetRoleResponse{
		Role: mappers.RoleWithPermissionsToProto(role),
	}, nil
}

// UpdateRole implements the UpdateRole RPC method.
func (s *Service) UpdateRole(ctx context.Context, req *securitypb.UpdateRoleRequest) (*securitypb.UpdateRoleResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.roles.update")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.UpdateRoleResponse{}, err
	}

	// Parse the user ID from the request
	roleID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &securitypb.UpdateRoleResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Construct the Role object from the request
	role := models.Role{
		Id:   roleID,
		Name: req.GetName(),
	}

	if ok, err := role.IsValid(); !ok {
		zap.L().Warn("Role is not valid", zap.Error(err))
		return &securitypb.UpdateRoleResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check user permissions for updating a role permissions
	err = security.Facade().CheckPermission(ctx, "admin.roles.permissions.update")

	// If the user has permission to update role permissions, validate and set them
	var permissions models.RolePermissionsInput
	if err == nil {

		// Construct the Permissions object from the request
		permissions = models.RolePermissionsInputFromUUIDs(req.GetPermissions())

		// Validate the permissions
		if ok, err := permissions.IsValid(); !ok {
			zap.L().Warn("Permissions are not valid", zap.Error(err))
			return &securitypb.UpdateRoleResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	// Update the role in the database
	err = repositories.R().R().Update(role, permissions)
	if err != nil {
		zap.L().Error("Cannot update role", zap.Error(err))
		return &securitypb.UpdateRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get the role + permissions from the database
	result, found, err := repositories.R().R().GetWithPermissions(roleID)
	if err != nil {
		zap.L().Error("Cannot get role", zap.String("uuid", roleID.String()), zap.Error(err))
		return &securitypb.UpdateRoleResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("Role not found after update", zap.String("uuid", roleID.String()))
		return &securitypb.UpdateRoleResponse{}, status.Error(codes.Internal, "Role not found after update")
	}

	// Convert the role to the gen format
	return &securitypb.UpdateRoleResponse{
		Role: mappers.RoleWithPermissionsToProto(result),
	}, nil
}

// DeleteRole implements the DeleteRole RPC method.
func (s *Service) DeleteRole(ctx context.Context, req *securitypb.DeleteRoleRequest) (*securitypb.DeleteRoleResponse, error) {

	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.roles.delete")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.DeleteRoleResponse{}, err
	}

	// Parse the user ID from the request
	roleID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &securitypb.DeleteRoleResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Delete the role from the database
	err = repositories.R().R().Delete(roleID)
	if err != nil {
		zap.L().Error("Cannot delete role", zap.String("uuid", roleID.String()), zap.Error(err))
		return &securitypb.DeleteRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &securitypb.DeleteRoleResponse{
		Success: true,
	}, nil
}

// ListRoles implements the ListRoles RPC method.
func (s *Service) ListRoles(ctx context.Context, req *securitypb.ListRolesRequest) (*securitypb.ListRolesResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.roles.list")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.ListRolesResponse{}, err
	}

	// Get all roles from the database
	result, err := repositories.R().R().ListWithPermissions()
	if err != nil {
		zap.L().Error("Cannot list roles", zap.Error(err))
		return &securitypb.ListRolesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &securitypb.ListRolesResponse{
		Roles: mappers.RolesWithPermissionsToProto(result),
	}, nil
}

// ListRolePermissions implements the ListRolePermissions RPC method.
func (s *Service) ListRolePermissions(ctx context.Context, req *securitypb.ListRolePermissionsRequest) (*securitypb.ListRolePermissionsResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.roles.permissions.list")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.ListRolePermissionsResponse{}, err
	}

	// Parse the user ID from the request
	roleID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &securitypb.ListRolePermissionsResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// List all role permissions from the database
	permissions, err := repositories.R().R().ListPermissionsByRoleId(roleID)
	if err != nil {
		zap.L().Error("Cannot list role permissions", zap.String("uuid", roleID.String()), zap.Error(err))
		return &securitypb.ListRolePermissionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &securitypb.ListRolePermissionsResponse{
		Permissions: mappers.PermissionsToProto(permissions),
	}, nil
}

// SetRolePermissions implements the SetRolePermissions RPC method.
func (s *Service) SetRolePermissions(ctx context.Context, req *securitypb.SetRolePermissionsRequest) (*securitypb.SetRolePermissionsResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.roles.permissions.update")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.SetRolePermissionsResponse{}, err
	}

	// Parse the user ID from the request
	roleID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid user ID", zap.String("user_id", req.GetId()), zap.Error(err))
		return &securitypb.SetRolePermissionsResponse{}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Construct the Permissions object from the request
	permissionsInput := models.RolePermissionsInputFromUUIDs(req.GetPermissions())
	if ok, err := permissionsInput.IsValid(); !ok {
		zap.L().Warn("Permissions are not valid", zap.Error(err))
		return &securitypb.SetRolePermissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Set the role permissions in the database
	err = repositories.R().R().SetPermissionsByRoleId(roleID, permissionsInput)
	if err != nil {
		zap.L().Error("Failed to set permissions", zap.Error(err))
		return &securitypb.SetRolePermissionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	// List all role permissions from the database
	permissions, err := repositories.R().R().ListPermissionsByRoleId(roleID)
	if err != nil {
		zap.L().Error("Cannot list role permissions", zap.String("uuid", roleID.String()), zap.Error(err))
		return &securitypb.SetRolePermissionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &securitypb.SetRolePermissionsResponse{
		Permissions: mappers.PermissionsToProto(permissions),
	}, nil
}
