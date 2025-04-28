package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PublicService is the implementation of the PublicSecurityService interface.
type PublicService struct {
	protogen.UnimplementedPublicSecurityServiceServer
}

// CheckPermission implements the CheckPermission RPC method.
func (s *PublicService) CheckPermission(ctx context.Context, req *protogen.CheckPermissionRequest) (*protogen.CheckPermissionResponse, error) {
	// Retrieve the user from the context
	_ctxUser := ctx.Value(app.ContextKeyUser)
	if _ctxUser != nil {
		zap.L().Warn("No context user provided")
		return &protogen.CheckPermissionResponse{
			HasPermission: false,
		}, nil
	}
	user, ok := _ctxUser.(models.User)
	if !ok {
		zap.L().Warn("Invalid user type in context")
		return &protogen.CheckPermissionResponse{
			HasPermission: false,
		}, nil
	}

	// If the user ID is provided in the request, we should check if it matches the user ID in the context
	if req.GetUserId() != "" && user.ID.String() == req.GetUserId() {
		// User is performing request for himself : authorized
		return &protogen.CheckPermissionResponse{
			HasPermission: true,
		}, nil
	}

	// Prepare roles data
	userRolesWithPermissions := models.RolesWithPermissions{}

	// Retrieve the userRoles from the context
	_ctxUserRolesWithPermissions := ctx.Value(app.ContextKeyUserRolesWithPermissions)
	if _ctxUserRolesWithPermissions != nil {
		// _ctxUserRoles to models.Roles
		_userRolesWithPermissions, ok := _ctxUserRolesWithPermissions.(models.RolesWithPermissions)
		if ok {
			userRolesWithPermissions = _userRolesWithPermissions
		}
	}

	// Could not retrieve userRoles from context, retrieving from database
	if len(userRolesWithPermissions) == 0 {
		result, err := repositories.R().R().ListWithPermissionsByUserId(user.ID)
		if err != nil {
			zap.L().Error("Cannot list a user roles with permissions", zap.String("uuid", user.ID.String()), zap.Error(err))
			return &protogen.CheckPermissionResponse{
				HasPermission: false,
			}, status.Error(codes.Internal, err.Error())
		}
		userRolesWithPermissions = result
	}

	// Check if the user has the permission
	if !userRolesWithPermissions.HasPermission(req.GetPermission()) {
		zap.L().Warn("Permission not found in context")
		return &protogen.CheckPermissionResponse{
			HasPermission: false,
		}, nil
	}

	return &protogen.CheckPermissionResponse{
		HasPermission: true,
	}, nil
}
