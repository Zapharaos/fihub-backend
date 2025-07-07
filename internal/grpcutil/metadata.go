package grpcutil

import (
	"context"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AddUserIDToContextMetadata(ctx context.Context, userID string) context.Context {
	if userID == "" {
		return ctx
	}

	// Create a new context with the user ID in metadata
	md := metadata.Pairs(string(app.ContextKeyUserID), userID)
	return metadata.NewOutgoingContext(ctx, md)
}

func GetUserIDFromContextMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Missing metadata")
	}

	userIDs := md[string(app.ContextKeyUserID)]
	if len(userIDs) == 0 {
		return "", status.Error(codes.Unauthenticated, "Missing user ID in metadata")
	}

	return userIDs[0], nil
}

func PropagateContextMetadata(ctx context.Context) context.Context {
	// If any, propagate metadata from the incoming context to the outgoing context
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		return metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}
