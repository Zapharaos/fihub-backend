package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreatePermission implements the CreatePermission RPC method.
func (s *Service) CreatePermission(ctx context.Context, req *protogen.CreatePermissionRequest) (*protogen.CreatePermissionResponse, error) {

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
		return &protogen.CreatePermissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Create the permission
	permissionID, err := repositories.R().P().Create(permission)
	if err != nil {
		zap.L().Warn("Create permission", zap.Error(err))
		return &protogen.CreatePermissionResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get the permission from the database
	permission, found, err := repositories.R().P().Get(permissionID)
	if err != nil {
		zap.L().Error("Cannot get permission", zap.String("uuid", permissionID.String()), zap.Error(err))
		return &protogen.CreatePermissionResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("Permission not found after creation", zap.String("uuid", permissionID.String()))
		return &protogen.CreatePermissionResponse{}, status.Error(codes.Internal, "Permission not found after creation")
	}

	return &protogen.CreatePermissionResponse{
		Permission: permission.ToProtogenPermission(),
	}, nil
}

// GetPermission implements the GetPermission RPC method.
func (s *Service) GetPermission(ctx context.Context, req *protogen.GetPermissionRequest) (*protogen.GetPermissionResponse, error) {

	// Parse the permission ID from the request
	permissionID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid permission ID", zap.String("permission_id", req.GetId()), zap.Error(err))
		return &protogen.GetPermissionResponse{}, status.Error(codes.InvalidArgument, "Invalid permission ID")
	}

	// Get the permission from the database
	permission, found, err := repositories.R().P().Get(permissionID)
	if err != nil {
		zap.L().Error("Cannot load permission", zap.String("uuid", permissionID.String()), zap.Error(err))
		return &protogen.GetPermissionResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Debug("Permission not found", zap.String("uuid", permissionID.String()))
		return &protogen.GetPermissionResponse{}, status.Error(codes.NotFound, "Permission not found")
	}

	return &protogen.GetPermissionResponse{
		Permission: permission.ToProtogenPermission(),
	}, nil
}

// UpdatePermission implements the UpdatePermission RPC method.
func (s *Service) UpdatePermission(ctx context.Context, req *protogen.UpdatePermissionRequest) (*protogen.UpdatePermissionResponse, error) {

	// Parse the permission ID from the request
	permissionID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid permission ID", zap.String("permission_id", req.GetId()), zap.Error(err))
		return &protogen.UpdatePermissionResponse{}, status.Error(codes.InvalidArgument, "Invalid permission ID")
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
		return &protogen.UpdatePermissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Update the permission
	err = repositories.R().P().Update(permission)
	if err != nil {
		zap.L().Error("Update permission", zap.Error(err))
		return &protogen.UpdatePermissionResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Get the permission from the database
	permission, found, err := repositories.R().P().Get(permissionID)
	if err != nil {
		zap.L().Error("Cannot get permission", zap.String("uuid", permissionID.String()), zap.Error(err))
		return &protogen.UpdatePermissionResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !found {
		zap.L().Error("Permission not found after update", zap.String("uuid", permissionID.String()))
		return &protogen.UpdatePermissionResponse{}, status.Error(codes.Internal, "Permission not found after update")
	}

	return &protogen.UpdatePermissionResponse{
		Permission: permission.ToProtogenPermission(),
	}, nil
}

// DeletePermission implements the DeletePermission RPC method.
func (s *Service) DeletePermission(ctx context.Context, req *protogen.DeletePermissionRequest) (*protogen.DeletePermissionResponse, error) {

	// Parse the permission ID from the request
	permissionID, err := uuid.Parse(req.GetId())
	if err != nil {
		// Log the error and return an invalid response
		zap.L().Error("Invalid permission ID", zap.String("permission_id", req.GetId()), zap.Error(err))
		return &protogen.DeletePermissionResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, "Invalid permission ID")
	}

	// Delete the permission
	err = repositories.R().P().Delete(permissionID)
	if err != nil {
		zap.L().Error("Delete permission", zap.Error(err))
		return &protogen.DeletePermissionResponse{
			Success: false,
		}, status.Error(codes.Internal, err.Error())
	}

	return &protogen.DeletePermissionResponse{
		Success: true,
	}, nil
}

// ListPermissions implements the ListPermissions RPC method.
func (s *Service) ListPermissions(ctx context.Context, req *protogen.ListPermissionsRequest) (*protogen.ListPermissionsResponse, error) {
	// Get all permissions from the database
	result, err := repositories.R().P().List()
	if err != nil {
		zap.L().Error("Cannot list permissions", zap.Error(err))
		return &protogen.ListPermissionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	// Convert the permissions to the protogen format
	protogenPermissions := make([]*protogen.Permission, len(result))
	for i, permission := range result {
		protogenPermissions[i] = permission.ToProtogenPermission()
	}

	return &protogen.ListPermissionsResponse{
		Permissions: protogenPermissions,
	}, nil
}
