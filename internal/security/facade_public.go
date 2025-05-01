package security

import (
	"context"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// PermissionChecker defines the methods required by the facade
type PermissionChecker interface {
	CheckPermission(ctx context.Context, req *securitypb.CheckPermissionRequest) (*securitypb.CheckPermissionResponse, error)
}

// GrpcClientAdapter adapts the securitypb.PublicSecurityServiceClient to PermissionChecker
type GrpcClientAdapter struct {
	client securitypb.PublicSecurityServiceClient
}

// CheckPermission implements PermissionChecker for the gRPC client
func (a *GrpcClientAdapter) CheckPermission(ctx context.Context, req *securitypb.CheckPermissionRequest) (*securitypb.CheckPermissionResponse, error) {
	return a.client.CheckPermission(ctx, req)
}

// NewGrpcClientAdapter creates a new adapter for the gRPC client
func NewGrpcClientAdapter(client securitypb.PublicSecurityServiceClient) *GrpcClientAdapter {
	return &GrpcClientAdapter{
		client: client,
	}
}

type PublicSecurityFacade struct {
	service PermissionChecker
}

func NewPublicSecurityFacade(service PermissionChecker) *PublicSecurityFacade {
	return &PublicSecurityFacade{
		service: service,
	}
}

// NewPublicSecurityFacadeWithGrpcClient creates a new facade with a gRPC client
func NewPublicSecurityFacadeWithGrpcClient(client securitypb.PublicSecurityServiceClient) *PublicSecurityFacade {
	return &PublicSecurityFacade{
		service: NewGrpcClientAdapter(client),
	}
}

// CheckPermission wraps the CheckPermission call
func (s *PublicSecurityFacade) CheckPermission(ctx context.Context, permission string, userIDs ...uuid.UUID) error {
	// If any, propagate metadata from the incoming context to the outgoing context
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// Setup the request
	req := &securitypb.CheckPermissionRequest{
		Permission: permission,
	}

	if len(userIDs) > 0 {
		req.UserId = userIDs[0].String()
	}

	response, err := s.service.CheckPermission(ctx, req)
	if err != nil {
		zap.L().Error("PublicSecurityFacade.CheckPermission", zap.Error(err))
		return err
	}

	if !response.GetHasPermission() {
		zap.L().Error("PublicSecurityFacade.PermissionDenied", zap.String("permission", permission))
		return status.Error(codes.PermissionDenied, "Permission denied")
	}

	return nil
}

// Global instance of the PublicSecurityFacade
var _globalPublicSecurityFacade *PublicSecurityFacade

// Facade returns the global PublicSecurityFacade instance
func Facade() *PublicSecurityFacade {
	return _globalPublicSecurityFacade
}

// ReplaceGlobals sets the global PublicSecurityFacade instance
func ReplaceGlobals(facade *PublicSecurityFacade) {
	_globalPublicSecurityFacade = facade
}
