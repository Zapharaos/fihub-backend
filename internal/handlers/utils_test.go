package handlers_test

import (
	"context"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/handlers"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCheckPermission tests the CheckPermission function
func TestCheckPermission(t *testing.T) {
	// Create a new controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Define the test cases
	tests := []struct {
		name         string
		contextUser  users.UserWithRoles
		contextOk    bool
		permission   string
		expectOK     bool
		expectStatus int
	}{
		{
			name:         "fail to retrieve user from context",
			contextUser:  users.UserWithRoles{},
			contextOk:    false,
			permission:   "some-permission",
			expectOK:     false,
			expectStatus: http.StatusUnauthorized,
		},
		{
			name:         "user does not have roles",
			contextUser:  users.UserWithRoles{Roles: roles.RolesWithPermissions{}},
			contextOk:    true,
			permission:   "some-permission",
			expectOK:     false,
			expectStatus: http.StatusForbidden,
		},
		{
			name: "user does not have permission",
			contextUser: users.UserWithRoles{Roles: roles.RolesWithPermissions{
				{Permissions: permissions.Permissions{}},
			}},
			contextOk:    true,
			permission:   "some-permission",
			expectOK:     false,
			expectStatus: http.StatusForbidden,
		},
		{
			name: "user has permission",
			contextUser: users.UserWithRoles{Roles: roles.RolesWithPermissions{
				{Permissions: permissions.Permissions{permissions.Permission{Value: "valid-permission"}}},
			}},
			contextOk:    true,
			permission:   "valid-permission",
			expectOK:     true,
			expectStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			// Set up the expectations
			m := mocks.NewMockUtils(ctrl)
			m.EXPECT().GetUserFromContext(r).Return(tt.contextUser, tt.contextOk)

			// Replace the global utils with the mock
			resolve := handlers.ReplaceGlobals(m)
			defer resolve()

			// Call the function
			ok := handlers.NewUtils().CheckPermission(w, r, tt.permission)

			// Check the results
			assert.Equal(t, tt.expectOK, ok)
			assert.Equal(t, tt.expectStatus, w.Code)
		})
	}
}

// TestGetUserFromContext tests the GetUserFromContext function
func TestGetUserFromContext(t *testing.T) {
	// Define valid user data
	user := users.UserWithRoles{
		User: users.User{
			ID: uuid.New(),
		},
	}

	// Replace the global utils with a new instance
	handlers.ReplaceGlobals(handlers.NewUtils())

	// Define the test cases
	tests := []struct {
		name       string
		context    context.Context
		expectOK   bool
		expectUser users.UserWithRoles
	}{
		{
			name:     "no context",
			context:  context.Background(),
			expectOK: false,
		},
		{
			name:     "wrong struct in context",
			context:  context.WithValue(context.Background(), app.ContextKeyUser, "wrong struct"),
			expectOK: false,
		},
		{
			name:       "valid user in context",
			context:    context.WithValue(context.Background(), app.ContextKeyUser, user),
			expectOK:   true,
			expectUser: user,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request with the context
			r := httptest.NewRequest("GET", "/", nil).WithContext(tt.context)

			// Call the function
			resultUser, ok := handlers.U().GetUserFromContext(r)

			// Check the results
			assert.Equal(t, tt.expectOK, ok)
			if tt.expectOK {
				assert.Equal(t, tt.expectUser, resultUser)
			}
		})
	}
}

// TestParseParamUUID tests the ParseParamUUID function
func TestParseParamUUID(t *testing.T) {
	// Define valid data
	validUUID, _ := uuid.NewUUID()
	validString := validUUID.String()

	// Replace the global utils with a new instance
	handlers.ReplaceGlobals(handlers.NewUtils())

	// Define the test cases
	tests := []struct {
		name       string
		paramValue string
		paramKey   string
		expectOK   bool
		expectCode int
		expectUUID uuid.UUID
	}{
		{
			name:       "missing UUID",
			paramValue: "",
			paramKey:   "id",
			expectOK:   false,
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "invalid UUID",
			paramValue: "invalid-uuid",
			paramKey:   "id",
			expectOK:   false,
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "valid UUID",
			paramValue: validString,
			paramKey:   "id",
			expectOK:   true,
			expectCode: http.StatusOK,
			expectUUID: validUUID,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			// Create a new route context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add(tt.paramKey, tt.paramValue)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			// Call the function
			resultUUID, ok := handlers.U().ParseParamUUID(w, r, tt.paramKey)

			// Check the results
			assert.Equal(t, tt.expectOK, ok)
			assert.Equal(t, tt.expectCode, w.Code)
			if tt.expectOK {
				assert.Equal(t, tt.expectUUID, resultUUID)
			}
		})
	}
}

// TestParseUUIDPair tests the ParseUUIDPair function
func TestParseUUIDPair(t *testing.T) {
	// Define valid data
	key := "key"
	baseUUID, _ := uuid.NewUUID()
	keyUUID, _ := uuid.NewUUID()

	// Create a new controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Define the test cases
	tests := []struct {
		name       string
		keyOK      bool
		keyUUID    uuid.UUID
		baseOK     bool
		baseUUID   uuid.UUID
		expectOK   bool
		expectBase uuid.UUID
		expectKey  uuid.UUID
	}{
		{
			name:       "invalid key UUID",
			keyOK:      false,
			keyUUID:    uuid.Nil,
			expectOK:   false,
			expectBase: uuid.Nil,
			expectKey:  uuid.Nil,
		},
		{
			name:       "invalid base UUID",
			keyOK:      true,
			keyUUID:    keyUUID,
			baseOK:     false,
			baseUUID:   uuid.Nil,
			expectOK:   false,
			expectBase: uuid.Nil,
			expectKey:  keyUUID,
		},
		{
			name:       "valid UUIDs",
			keyOK:      true,
			keyUUID:    keyUUID,
			baseOK:     true,
			baseUUID:   baseUUID,
			expectOK:   true,
			expectBase: baseUUID,
			expectKey:  keyUUID,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			// Set up the expectations
			m := mocks.NewMockUtils(ctrl)
			m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), key).Return(tt.keyUUID, tt.keyOK)
			if tt.keyOK {
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(tt.baseUUID, tt.baseOK)
			}

			// Replace the global utils with the mock
			resolve := handlers.ReplaceGlobals(m)
			defer resolve()

			// Call the function
			baseID, keyID, ok := handlers.NewUtils().ParseUUIDPair(w, r, key)

			// Check the results
			assert.Equal(t, tt.expectOK, ok)
			assert.Equal(t, tt.expectBase, baseID)
			assert.Equal(t, tt.expectKey, keyID)
		})
	}
}
