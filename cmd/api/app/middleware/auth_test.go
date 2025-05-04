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
				authClient.EXPECT().ExtractUserID(gomock.Any(), gomock.Any()).Times(0)
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
				authClient.EXPECT().ExtractUserID(gomock.Any(), gomock.Any()).Return(&authpb.ExtractUserIDResponse{}, errors.New("some error"))
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
				authClient.EXPECT().ExtractUserID(gomock.Any(), gomock.Any()).Times(0)
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
				authClient.EXPECT().ExtractUserID(gomock.Any(), gomock.Any()).Return(&authpb.ExtractUserIDResponse{
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
				authClient.EXPECT().ExtractUserID(gomock.Any(), gomock.Any()).Times(0)
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

/*// TestMiddleware tests the Middleware function
func TestMiddleware(t *testing.T) {
	// Define test data
	a := server.New(server.CheckHeader, server.Config{})
	user := models.User{
		ID: uuid.New(),
	}

	// Generate a valid token
	validToken, err := a.GenerateToken(user)
	assert.NoError(t, err)

	// Create a new httptest server
	ts := httptest.NewServer(http.HandlerFunc(a.GetToken))
	defer ts.Close()

	// Define test cases
	tests := []struct {
		name        string                       // Test case name
		token       string                       // Token to set in the request
		target      string                       // Target URL
		mockSetup   func(*gomock.Controller)     // Mock setup function
		roles       mocks.SecurityRoleRepository // Roles repository mocks
		expectCode  int                          // Expected status code
		expectCtx   bool                         // Expected context
		expectErr   bool
		expectFound bool
	}{
		{
			name:   "skip middleware",
			target: "/users",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(u)
			},
			expectCode: http.StatusOK,
			expectCtx:  false,
		},
		{
			name:  "empty token",
			token: "",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(u)
			},
			expectCode: http.StatusUnauthorized,
			expectCtx:  false,
		},
		{
			name:  "invalid token",
			token: "invalid.token.string",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(u)
			},
			expectCode: http.StatusUnauthorized,
			expectCtx:  false,
		},
		{
			name:  "load user error",
			token: validToken.Token,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(models.User{}, false, errors.New("error"))
				repositories.ReplaceGlobals(u)
				r := mocks.NewSecurityRoleRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Times(0)
				securityrepositories.ReplaceGlobals(securityrepositories.NewRepository(r, nil))
			},
			expectCode: http.StatusUnauthorized,
			expectCtx:  false,
		},
		{
			name:  "valid token",
			token: validToken.Token,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUserRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(user, true, nil)
				repositories.ReplaceGlobals(u)
				r := mocks.NewSecurityRoleRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return([]models.Role{{Id: uuid.New(), Name: "admin"}}, nil)
				p := mocks.NewSecurityPermissionRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Return(models.Permissions{}, nil)
				securityrepositories.ReplaceGlobals(securityrepositories.NewRepository(r, p))
			},
			expectErr:   false,
			expectFound: true,
			expectCode:  http.StatusOK,
			expectCtx:   true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Apply mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tt.mockSetup(ctrl)

			// Create a new request
			apiBasePath := viper.GetString("API_BASE_PATH")
			r := httptest.NewRequest("POST", ts.URL+apiBasePath+tt.target, nil)
			w := httptest.NewRecorder()
			if tt.token != "" {
				r.Header.Set("Authorization", tt.token)
			}

			// Create a simple handler function
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check the context
				_, ok := r.Context().Value(app.ContextKeyUser).(models.UserWithRoles)
				// If the context is expected, check if it is set
				assert.Equal(t, tt.expectCtx, ok)
			})

			// Setups the middleware and
			handler := a.Middleware(next)
			handler.ServeHTTP(w, r)
			response := w.Result()
			defer response.Body.Close()

			// Check the response
			assert.Equal(t, tt.expectCode, response.StatusCode)
		})
	}
}*/
