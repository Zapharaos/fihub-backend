package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/internal/grpcutil"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PublicService is the implementation of the PublicSecurityService interface.
type PublicService struct {
	securitypb.UnimplementedPublicSecurityServiceServer
}

// CheckPermission implements the CheckPermission RPC method.
func (s *PublicService) CheckPermission(ctx context.Context, req *securitypb.CheckPermissionRequest) (*securitypb.CheckPermissionResponse, error) {

	// Retrieve the user ID from the context
	userID, err := grpcutil.GetUserIDFromContextMetadata(ctx)
	if err != nil {
		return &securitypb.CheckPermissionResponse{
			HasPermission: false,
		}, err
	}

	// If the user ID is provided in the request, we should check if it matches the user ID in the metadata
	if req.GetUserId() != "" && userID == req.GetUserId() {
		// User is performing request for himself : authorized
		return &securitypb.CheckPermissionResponse{
			HasPermission: true,
		}, nil
	}

	// Retrieve the user roles with permissions from the repository
	userRolesWithPermissions, err := repositories.R().R().ListWithPermissionsByUserId(uuid.MustParse(userID))
	if err != nil {
		zap.L().Error("Cannot list a user roles with permissions", zap.String("uuid", userID), zap.Error(err))
		return &securitypb.CheckPermissionResponse{
			HasPermission: false,
		}, status.Error(codes.Internal, err.Error())
	}

	// Check if the user has the permission
	if !userRolesWithPermissions.HasPermission(req.GetPermission()) {
		zap.L().Warn("Permission not found in context")
		return &securitypb.CheckPermissionResponse{
			HasPermission: false,
		}, status.Error(codes.PermissionDenied, "Missing permission for action")
	}

	return &securitypb.CheckPermissionResponse{
		HasPermission: true,
	}, nil
}
