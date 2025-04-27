package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"go.uber.org/zap"
)

// PublicService is the implementation of the PublicSecurityService interface.
type PublicService struct {
	protogen.UnimplementedPublicSecurityServiceServer
}

// CheckPermission implements the CheckPermission RPC method.
func (s *PublicService) CheckPermission(ctx context.Context, req *protogen.CheckPermissionRequest) (*protogen.CheckPermissionResponse, error) {
	// TODO : separate ctx models.UserWithRoles into public user vs roles
	_user := ctx.Value(app.ContextKeyUser)
	if _user == nil {
		zap.L().Warn("No context user provided")
		return &protogen.CheckPermissionResponse{
			HasPermission: false,
		}, nil
	}
	user, ok := _user.(models.UserWithRoles)
	if !ok {
		zap.L().Warn("Invalid user type in context")
		return &protogen.CheckPermissionResponse{
			HasPermission: false,
		}, nil
	}

	// TODO : check user id context with request user id
	// TODO : does not match : retrieve user roles from db

	// check if the user has the permission
	if !user.HasPermission(req.GetPermission()) {
		zap.L().Warn("Permission not found in context")
		return &protogen.CheckPermissionResponse{
			HasPermission: false,
		}, nil
	}

	return &protogen.CheckPermissionResponse{
		HasPermission: true,
	}, nil
}
