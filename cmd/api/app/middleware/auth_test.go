package middleware

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/server"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	inputUserID := uuid.New().String()

	// Define test cases
	tests := []struct {
		name       string
		mockSetup  func(ctrl *gomock.Controller)
		config     server.Config
		expectCode int
		expectCtx  bool
	}{
		{
			name: "no security mode",
			mockSetup: func(ctrl *gomock.Controller) {
				authClient := mocks.NewMockAuthServiceClient(ctrl)
				authClient.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).Times(0)
				authClient.EXPECT().ExtractUserIDFromToken(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(authClient),
				))
			},
			config: server.Config{
				Security: false,
			},
			expectCode: http.StatusOK,
			expectCtx:  false,
		},
		{
			name: "fails in gateway mode",
			mockSetup: func(ctrl *gomock.Controller) {
				authClient := mocks.NewMockAuthServiceClient(ctrl)
				authClient.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).Times(0)
				authClient.EXPECT().ExtractUserIDFromToken(gomock.Any(), gomock.Any()).Return(&authpb.ExtractUserIDFromTokenResponse{}, errors.New("some error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(authClient),
				))
			},
			config: server.Config{
				Security:    true,
				GatewayMode: true,
			},
			expectCode: http.StatusBadRequest,
			expectCtx:  false,
		},
		{
			name: "fails in default mode",
			mockSetup: func(ctrl *gomock.Controller) {
				authClient := mocks.NewMockAuthServiceClient(ctrl)
				authClient.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).Return(&authpb.ValidateTokenResponse{}, errors.New("some error"))
				authClient.EXPECT().ExtractUserIDFromToken(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(authClient),
				))
			},
			config: server.Config{
				Security: true,
			},
			expectCode: http.StatusUnauthorized,
			expectCtx:  false,
		},
		{
			name: "success in gateway mode",
			mockSetup: func(ctrl *gomock.Controller) {
				authClient := mocks.NewMockAuthServiceClient(ctrl)
				authClient.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).Times(0)
				authClient.EXPECT().ExtractUserIDFromToken(gomock.Any(), gomock.Any()).Return(&authpb.ExtractUserIDFromTokenResponse{
					UserId: inputUserID,
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(authClient),
				))
			},
			config: server.Config{
				Security:    true,
				GatewayMode: true,
			},
			expectCode: http.StatusOK,
			expectCtx:  true,
		},
		{
			name: "success in default mode",
			mockSetup: func(ctrl *gomock.Controller) {
				authClient := mocks.NewMockAuthServiceClient(ctrl)
				authClient.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).Return(&authpb.ValidateTokenResponse{
					UserId: inputUserID,
				}, nil)
				authClient.EXPECT().ExtractUserIDFromToken(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(authClient),
				))
			},
			config: server.Config{
				Security: true,
			},
			expectCode: http.StatusOK,
			expectCtx:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock dependencies
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Create middleware
			middleware := AuthMiddleware(tt.config)

			// Create a test request
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", "token")
			rr := httptest.NewRecorder()

			// Create a next handler to verify context
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, ok := r.Context().Value(app.ContextKeyUserID).(string)
				if tt.expectCtx {
					assert.True(t, ok, "User ID should be set in context")
					assert.Equal(t, inputUserID, userID, "User ID should match")
				} else {
					assert.False(t, ok, "User ID should not be set in context")
				}
				w.WriteHeader(http.StatusOK)
			})

			// Execute middleware
			handler := middleware(next)
			handler.ServeHTTP(rr, req)

			// Verify response
			assert.Equal(t, tt.expectCode, rr.Code, "Response status code should match")
		})
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name        string
		headerToken string
		queryToken  string
		want        string
	}{
		{
			name:        "token in header",
			headerToken: "header-token",
			queryToken:  "query-token",
			want:        "header-token",
		},
		{
			name:        "token in query only",
			headerToken: "",
			queryToken:  "query-token",
			want:        "query-token",
		},
		{
			name:        "no token",
			headerToken: "",
			queryToken:  "",
			want:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.headerToken != "" {
				req.Header.Set("Authorization", tt.headerToken)
			}
			if tt.queryToken != "" {
				q := req.URL.Query()
				q.Set("token", tt.queryToken)
				req.URL.RawQuery = q.Encode()
			}
			got := extractToken(req)
			assert.Equal(t, tt.want, got)
		})
	}
}
