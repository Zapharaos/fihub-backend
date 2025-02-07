package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/auth"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	// Create a new instance of Auth
	a := auth.New(auth.CheckHeader | auth.CheckQuery)

	// Check the instance
	assert.NotNil(t, a)
	assert.NotEmpty(t, a.SigningKey)
	assert.Equal(t, int8(auth.CheckHeader|auth.CheckQuery), a.Checks)
}

// TestGetToken tests the GetToken function
func TestGetToken(t *testing.T) {
	// Define test data
	a := auth.New(auth.CheckHeader)
	userWithPassword := users.UserWithPassword{
		User:     users.User{Email: "test@example.com"},
		Password: "password",
	}
	userBody, _ := json.Marshal(userWithPassword)

	// Create a new httptest server
	ts := httptest.NewServer(http.HandlerFunc(a.GetToken))
	defer ts.Close()

	// Define test cases
	tests := []struct {
		name        string // Test case name
		mockSetup   func(*gomock.Controller)
		body        []byte
		status      int  // Expected status code
		expectEmpty bool // Expected empty token
	}{
		{
			name:        "invalid data",
			mockSetup:   func(ctrl *gomock.Controller) {},
			status:      http.StatusBadRequest,
			body:        nil,
			expectEmpty: true,
		},
		{
			name: "authentication error",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewUsersRepository(ctrl)
				m.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(users.User{}, false, errors.New("error"))
				users.ReplaceGlobals(m)
			},
			body:        userBody,
			status:      http.StatusInternalServerError,
			expectEmpty: true,
		},
		{
			name: "authentication failed",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewUsersRepository(ctrl)
				m.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(users.User{}, false, nil)
				users.ReplaceGlobals(m)
			},
			body:        userBody,
			status:      http.StatusBadRequest,
			expectEmpty: true,
		},
		{
			name: "authentication success",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewUsersRepository(ctrl)
				m.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(userWithPassword.User, true, nil)
				users.ReplaceGlobals(m)
			},
			body:        userBody,
			status:      http.StatusOK,
			expectEmpty: false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request to the test server
			r := httptest.NewRequest("POST", ts.URL+"/api/v1/auth/token", bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()

			// Apply mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tt.mockSetup(ctrl)

			// Send the request
			a.GetToken(w, r)
			response := w.Result()
			defer response.Body.Close()

			// Get the response body
			var token auth.JwtToken
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
	a := auth.New(auth.CheckHeader)
	user := users.User{
		ID: uuid.New(),
	}

	// Generate a token
	token, err := a.GenerateToken(user)

	// Check the response
	assert.NoError(t, err)
	assert.NotEmpty(t, token.Token)
}

// TestValidateToken tests the ValidateToken function
func TestValidateToken(t *testing.T) {
	// Define test data
	a := auth.New(auth.CheckHeader)
	user := users.User{
		ID: uuid.New(),
	}

	// Generate a valid token
	validToken, err := a.GenerateToken(user)
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
			claims, err := a.ValidateToken(tt.token)
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
	a := auth.New(auth.CheckHeader)
	user := users.User{
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
		name        string                   // Test case name
		token       string                   // Token to set in the request
		target      string                   // Target URL
		mockSetup   func(*gomock.Controller) // Mock setup function
		roles       mocks.RolesRepository    // Roles repository mocks
		expectCode  int                      // Expected status code
		expectCtx   bool                     // Expected context
		expectErr   bool
		expectFound bool
	}{
		{
			name:   "skip middleware",
			target: "/users",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectCode: http.StatusOK,
			expectCtx:  false,
		},
		{
			name:  "empty token",
			token: "",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectCode: http.StatusUnauthorized,
			expectCtx:  false,
		},
		{
			name:  "invalid token",
			token: "invalid.token.string",
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectCode: http.StatusUnauthorized,
			expectCtx:  false,
		},
		{
			name:  "load user error",
			token: validToken.Token,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(users.User{}, false, errors.New("error"))
				users.ReplaceGlobals(u)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Times(0)
				roles.ReplaceGlobals(r)
			},
			expectCode: http.StatusUnauthorized,
			expectCtx:  false,
		},
		{
			name:  "valid token",
			token: validToken.Token,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(user, true, nil)
				users.ReplaceGlobals(u)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return([]roles.Role{{Id: uuid.New(), Name: "admin"}}, nil)
				roles.ReplaceGlobals(r)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Return([]permissions.Permission{}, nil)
				permissions.ReplaceGlobals(p)
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
			handler := a.Middleware(next)
			handler.ServeHTTP(w, r)
			response := w.Result()
			defer response.Body.Close()

			// Check the response
			assert.Equal(t, tt.expectCode, response.StatusCode)
		})
	}
}
