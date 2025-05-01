package middleware

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/server"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/protogen"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"net/http"
)

// extractToken extracts the token from the request.
func extractToken(r *http.Request) string {
	if token := r.Header.Get("Authorization"); token != "" {
		return token
	}
	return r.URL.Query().Get("token")
}

// AuthMiddleware is a middleware for authenticating requests.
func AuthMiddleware(config server.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// No auth in insecure mode
			if !config.Security {
				next.ServeHTTP(w, r)
				return
			}

			token := extractToken(r)
			userID := ""

			// Gateway mode: skip validation
			// WARNING: this is a security risk, don't use unless you know what you're doing.
			if config.GatewayMode {
				// Extract user ID from token
				response, err := clients.C().Auth().ExtractUserID(r.Context(), &protogen.ExtractUserIDRequest{
					Token: token,
				})
				if err != nil {
					zap.L().Error("ExtractUserID", zap.Error(err))
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				userID = response.GetUserId()
			} else {
				// Validate token
				response, err := clients.C().Auth().ValidateToken(r.Context(), &protogen.ValidateTokenRequest{
					Token: token,
				})
				if err != nil {
					zap.L().Error("ValidateToken", zap.Error(err))
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				userID = response.GetUserId()
			}

			// Setup metadata for gRPC clients as context
			md := metadata.Pairs("x-user-id", userID)
			ctx := metadata.NewOutgoingContext(r.Context(), md)
			r = r.WithContext(ctx)

			// Set user ID in context
			ctx = context.WithValue(r.Context(), app.ContextKeyUserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
