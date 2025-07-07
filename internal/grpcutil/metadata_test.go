package grpcutil

import (
	"context"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"testing"
)

func TestAddUserIDToContextMetadata(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		expectKey bool
	}{
		{"empty userID", "", false},
		{"valid userID", "user123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = AddUserIDToContextMetadata(ctx, tt.userID)
			md, ok := metadata.FromOutgoingContext(ctx)
			if tt.expectKey {
				if !ok || len(md[string(app.ContextKeyUserID)]) == 0 || md[string(app.ContextKeyUserID)][0] != tt.userID {
					t.Errorf("expected userID in metadata")
				}
			} else {
				if ok && len(md[string(app.ContextKeyUserID)]) > 0 {
					t.Errorf("did not expect userID in metadata")
				}
			}
		})
	}
}

func TestGetUserIDFromContextMetadata(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		expectErr bool
		expectVal string
	}{
		{
			name:      "no metadata",
			ctx:       context.Background(),
			expectErr: true,
		},
		{
			name:      "no userID in metadata",
			ctx:       metadata.NewIncomingContext(context.Background(), metadata.Pairs()),
			expectErr: true,
		},
		{
			name:      "userID present",
			ctx:       metadata.NewIncomingContext(context.Background(), metadata.Pairs(string(app.ContextKeyUserID), "user123")),
			expectErr: false,
			expectVal: "user123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := GetUserIDFromContextMetadata(tt.ctx)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				st, _ := status.FromError(err)
				if st == nil || st.Code() != codes.Unauthenticated {
					t.Errorf("expected Unauthenticated error code")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if val != tt.expectVal {
					t.Errorf("expected %v, got %v", tt.expectVal, val)
				}
			}
		})
	}
}

func TestPropagateContextMetadata(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		expectKey bool
	}{
		{
			name:      "no incoming metadata",
			ctx:       context.Background(),
			expectKey: false,
		},
		{
			name:      "with incoming metadata",
			ctx:       metadata.NewIncomingContext(context.Background(), metadata.Pairs("foo", "bar")),
			expectKey: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outCtx := PropagateContextMetadata(tt.ctx)
			md, ok := metadata.FromOutgoingContext(outCtx)
			if tt.expectKey {
				if !ok || md["foo"][0] != "bar" {
					t.Errorf("expected metadata to be propagated")
				}
			} else {
				if ok && len(md["foo"]) > 0 {
					t.Errorf("did not expect metadata to be propagated")
				}
			}
		})
	}
}
