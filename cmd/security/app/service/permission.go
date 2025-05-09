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

// CreatePermission implements the CreatePermission RPC method.
func (s *Service) CreatePermission(ctx context.Context, req *securitypb.CreatePermissionRequest) (*securitypb.CreatePermissionResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.permissions.create")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.CreatePermissionResponse{}, err
	}

	// Construct the Permission object from the request
	permission := models.Permission{
		Id:          uuid.New(),
		Value:       req.GetValue(),
		Scope:       req.GetScope(),
		Description: req.GetDescription(),
	}

	// Validate the permission
	if ok, err := permission.IsValid(); !ok {
		zap.L().Warn("Permission is not valid", zap.Error(err))
		return &securitypb.CreatePermissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Create the permission
	permissionID, err := repositories.R().P().Create(permission)
	if err != nil {
		zap.L().Warn("Create permission", zap.Error(err))
		return &securitypb.CreatePermissionResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get the permission from the database
	permission, found, err := repositories.R().P().Get(permissionID)
	if err != nil {
		zap.L().Error("Cannot get permission", zap.String("uuid", permissionID.String()), zap.Error(err))
		return &securitypb.CreatePermissionResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("Permission not found after creation", zap.String("uuid", permissionID.String()))
		return &securitypb.CreatePermissionResponse{}, status.Error(codes.Internal, "Permission not found after creation")
	}

	return &securitypb.CreatePermissionResponse{
		Permission: mappers.PermissionToProto(permission),
	}, nil
}

// GetPermission implements the GetPermission RPC method.
func (s *Service) GetPermission(ctx context.Context, req *securitypb.GetPermissionRequest) (*securitypb.GetPermissionResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.permissions.read")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.GetPermissionResponse{}, err
	}

	// Parse the permission ID from the request
	permissionID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid permission ID", zap.String("permission_id", req.GetId()), zap.Error(err))
		return &securitypb.GetPermissionResponse{}, status.Error(codes.InvalidArgument, "Invalid permission ID")
	}

	// Get the permission from the database
	permission, found, err := repositories.R().P().Get(permissionID)
	if err != nil {
		zap.L().Error("Cannot load permission", zap.String("uuid", permissionID.String()), zap.Error(err))
		return &securitypb.GetPermissionResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Debug("Permission not found", zap.String("uuid", permissionID.String()))
		return &securitypb.GetPermissionResponse{}, status.Error(codes.NotFound, "Permission not found")
	}

	return &securitypb.GetPermissionResponse{
		Permission: mappers.PermissionToProto(permission),
	}, nil
}

// UpdatePermission implements the UpdatePermission RPC method.
func (s *Service) UpdatePermission(ctx context.Context, req *securitypb.UpdatePermissionRequest) (*securitypb.UpdatePermissionResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.permissions.update")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.UpdatePermissionResponse{}, err
	}

	// Parse the permission ID from the request
	permissionID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid permission ID", zap.String("permission_id", req.GetId()), zap.Error(err))
		return &securitypb.UpdatePermissionResponse{}, status.Error(codes.InvalidArgument, "Invalid permission ID")
	}

	// Construct the Permission object from the request
	permission := models.Permission{
		Id:          permissionID,
		Value:       req.GetValue(),
		Scope:       req.GetScope(),
		Description: req.GetDescription(),
	}

	// Validate the permission
	if ok, err := permission.IsValid(); !ok {
		zap.L().Warn("Permission is not valid", zap.Error(err))
		return &securitypb.UpdatePermissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Update the permission
	err = repositories.R().P().Update(permission)
	if err != nil {
		zap.L().Error("Update permission", zap.Error(err))
		return &securitypb.UpdatePermissionResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get the permission from the database
	permission, found, err := repositories.R().P().Get(permissionID)
	if err != nil {
		zap.L().Error("Cannot get permission", zap.String("uuid", permissionID.String()), zap.Error(err))
		return &securitypb.UpdatePermissionResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("Permission not found after update", zap.String("uuid", permissionID.String()))
		return &securitypb.UpdatePermissionResponse{}, status.Error(codes.Internal, "Permission not found after update")
	}

	return &securitypb.UpdatePermissionResponse{
		Permission: mappers.PermissionToProto(permission),
	}, nil
}

// DeletePermission implements the DeletePermission RPC method.
func (s *Service) DeletePermission(ctx context.Context, req *securitypb.DeletePermissionRequest) (*securitypb.DeletePermissionResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.permissions.delete")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.DeletePermissionResponse{}, err
	}

	// Parse the permission ID from the request
	permissionID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid permission ID", zap.String("permission_id", req.GetId()), zap.Error(err))
		return &securitypb.DeletePermissionResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, "Invalid permission ID")
	}

	// Delete the permission
	err = repositories.R().P().Delete(permissionID)
	if err != nil {
		zap.L().Error("Delete permission", zap.Error(err))
		return &securitypb.DeletePermissionResponse{
			Success: false,
		}, status.Error(codes.Internal, err.Error())
	}

	return &securitypb.DeletePermissionResponse{
		Success: true,
	}, nil
}

// ListPermissions implements the ListPermissions RPC method.
func (s *Service) ListPermissions(ctx context.Context, req *securitypb.ListPermissionsRequest) (*securitypb.ListPermissionsResponse, error) {
	// Check user permissions
	err := security.Facade().CheckPermission(ctx, "admin.permissions.list")
	if err != nil {
		zap.L().Error("CheckPermission", zap.Error(err))
		return &securitypb.ListPermissionsResponse{}, err
	}

	// Get all permissions from the database
	result, err := repositories.R().P().List()
	if err != nil {
		zap.L().Error("Cannot list permissions", zap.Error(err))
		return &securitypb.ListPermissionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Convert the permissions to the gen format
	return &securitypb.ListPermissionsResponse{
		Permissions: mappers.PermissionsToProto(result),
	}, nil
}
