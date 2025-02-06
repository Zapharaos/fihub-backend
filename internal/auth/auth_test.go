package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	// Create a new instance of Auth
	auth := New(CheckHeader | CheckQuery)

	// Check the instance
	assert.NotNil(t, auth)
	assert.NotEmpty(t, auth.signingKey)
	assert.Equal(t, int8(CheckHeader|CheckQuery), auth.checks)
}

// TestGetToken tests the GetToken function
func TestGetToken(t *testing.T) {
	// Define test data
	auth := New(CheckHeader)
	userWithPassword := users.UserWithPassword{
		User:     users.User{Email: "test@example.com"},
		Password: "password",
	}
	userBody, _ := json.Marshal(userWithPassword)

	// Create a new httptest server
	ts := httptest.NewServer(http.HandlerFunc(auth.GetToken))
	defer ts.Close()

	// Define test cases
	tests := []struct {
		name        string                // Test case name
		users       mocks.UsersRepository // Users repository mocks
		body        []byte
		status      int  // Expected status code
		expectEmpty bool // Expected empty token
	}{
		{
			name:        "invalid data",
			status:      http.StatusBadRequest,
			body:        nil,
			expectEmpty: true,
		},
		{
			name:        "authentication error",
			users:       mocks.UsersRepository{Err: errors.New("error")},
			body:        userBody,
			status:      http.StatusInternalServerError,
			expectEmpty: true,
		},
		{
			name:        "authentication failed",
			users:       mocks.UsersRepository{Found: false},
			body:        userBody,
			status:      http.StatusBadRequest,
			expectEmpty: true,
		},
		{
			name:        "authentication success",
			users:       mocks.UsersRepository{Found: true, User: userWithPassword.User},
			body:        userBody,
			status:      http.StatusOK,
			expectEmpty: false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock repository
			users.ReplaceGlobals(mocks.NewUsersRepository(tt.users))

			// Create a new request to the test server
			r := httptest.NewRequest("POST", ts.URL+"/api/v1/auth/token", bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()

			// Send the request
			auth.GetToken(w, r)
			response := w.Result()
			defer response.Body.Close()

			// Get the response body
			var token JwtToken
			data, err := io.ReadAll(response.Body)
			err = json.Unmarshal(data, &token)

			// Check the response
			assert.Equal(t, tt.status, response.StatusCode)
			if tt.expectEmpty {
				assert.Empty(t, token.Token)
			} else {
				assert.NotEmpty(t, token.Token)
				assert.NoError(t, err)
			}
		})
	}
}

// TestGenerateToken tests the GenerateToken function
func TestGenerateToken(t *testing.T) {
	// Define test data
	auth := New(CheckHeader)
	user := users.User{
		ID: uuid.New(),
	}

	// Generate a token
	token, err := auth.GenerateToken(user)

	// Check the response
	assert.NoError(t, err)
	assert.NotEmpty(t, token.Token)
}

// TestValidateToken tests the ValidateToken function
func TestValidateToken(t *testing.T) {
	// Define test data
	auth := New(CheckHeader)
	user := users.User{
		ID: uuid.New(),
	}

	// Generate a valid token
	validToken, err := auth.GenerateToken(user)
	assert.NoError(t, err)

	// Define test cases
	tests := []struct {
		name      string
		token     string
		expectErr bool
	}{
		{
			name:      "valid token",
			token:     validToken.Token,
			expectErr: false,
		},
		{
			name:      "invalid token",
			token:     "invalid.token.string",
			expectErr: true,
		},
		{
			name:      "empty token",
			token:     "",
			expectErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := auth.ValidateToken(tt.token)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, user.ID.String(), claims["id"])
			}
		})
	}
}

// TestMiddleware tests the Middleware function
func TestMiddleware(t *testing.T) {
	// Define test data
	auth := New(CheckHeader)
	user := users.User{
		ID: uuid.New(),
	}

	// Generate a valid token
	validToken, err := auth.GenerateToken(user)
	assert.NoError(t, err)

	// Create a new httptest server
	ts := httptest.NewServer(http.HandlerFunc(auth.GetToken))
	defer ts.Close()

	// Define test cases
	tests := []struct {
		name        string                // Test case name
		token       string                // Token to set in the request
		target      string                // Target URL
		users       mocks.UsersRepository // Users repository mocks
		roles       mocks.RolesRepository // Roles repository mocks
		expectCode  int                   // Expected status code
		expectCtx   bool                  // Expected context
		expectErr   bool
		expectFound bool
	}{
		{
			name:       "skip middleware",
			target:     "/users",
			expectCode: http.StatusOK,
			expectCtx:  false,
		},
		{
			name:       "empty token",
			token:      "",
			expectCode: http.StatusUnauthorized,
			expectCtx:  false,
		},
		{
			name:       "invalid token",
			token:      "invalid.token.string",
			expectCode: http.StatusUnauthorized,
			expectCtx:  false,
		},
		{
			name:  "load user error",
			token: validToken.Token,
			users: mocks.UsersRepository{
				Err:   fmt.Errorf("error"),
				Found: false,
			},
			expectCode: http.StatusUnauthorized,
			expectCtx:  false,
		},
		{
			name:        "valid token",
			token:       validToken.Token,
			expectErr:   false,
			expectFound: true,
			users: mocks.UsersRepository{
				User:  user,
				Found: true,
			},
			roles: mocks.RolesRepository{
				Roles: []roles.Role{{Id: uuid.New(), Name: "admin"}},
			},
			expectCode: http.StatusOK,
			expectCtx:  true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Mock the repositories
			users.ReplaceGlobals(mocks.NewUsersRepository(tt.users))
			roles.ReplaceGlobals(mocks.NewRolesRepository(tt.roles))
			permissions.ReplaceGlobals(mocks.NewPermissionsRepository(mocks.PermissionsRepository{}))

			// Create a new request
			r := httptest.NewRequest("POST", ts.URL+"/api/v1"+tt.target, nil)
			w := httptest.NewRecorder()
			if tt.token != "" {
				r.Header.Set("Authorization", tt.token)
			}

			// Create a simple handler function
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check the context
				_, ok := r.Context().Value(app.ContextKeyUser).(users.UserWithRoles)
				// If the context is expected, check if it is set
				assert.Equal(t, tt.expectCtx, ok)
			})

			// Setups the middleware and
			handler := auth.Middleware(next)
			handler.ServeHTTP(w, r)

			// Call the handler
			handler.ServeHTTP(w, r)
			response := w.Result()
			defer response.Body.Close()

			// Check the response
			assert.Equal(t, tt.expectCode, response.StatusCode)
		})
	}
}
