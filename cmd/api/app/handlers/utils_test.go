package handlers_test

import (
	"bytes"
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"
	"mime/multipart"
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
		contextUser  models.UserWithRoles
		contextOk    bool
		permission   string
		expectOK     bool
		expectStatus int
	}{
		{
			name:         "fail to retrieve user from context",
			contextUser:  models.UserWithRoles{},
			contextOk:    false,
			permission:   "some-permission",
			expectOK:     false,
			expectStatus: http.StatusUnauthorized,
		},
		{
			name:         "user does not have roles",
			contextUser:  models.UserWithRoles{Roles: models.RolesWithPermissions{}},
			contextOk:    true,
			permission:   "some-permission",
			expectOK:     false,
			expectStatus: http.StatusForbidden,
		},
		{
			name: "user does not have permission",
			contextUser: models.UserWithRoles{Roles: models.RolesWithPermissions{
				{Permissions: models.Permissions{}},
			}},
			contextOk:    true,
			permission:   "some-permission",
			expectOK:     false,
			expectStatus: http.StatusForbidden,
		},
		{
			name: "user has permission",
			contextUser: models.UserWithRoles{Roles: models.RolesWithPermissions{
				{Permissions: models.Permissions{models.Permission{Value: "valid-permission"}}},
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
	user := models.UserWithRoles{
		User: models.User{
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
		expectUser models.UserWithRoles
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

// TestParseParamString tests the ParseParamString function
func TestParseParamString(t *testing.T) {
	// Replace the global utils with a new instance
	handlers.ReplaceGlobals(handlers.NewUtils())

	// Define the test cases
	tests := []struct {
		name       string
		paramValue string
		paramKey   string
		expectOK   bool
		expectCode int
		expectStr  string
	}{
		{
			name:       "missing string",
			paramValue: "",
			paramKey:   "id",
			expectOK:   false,
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "valid string",
			paramValue: "valid-string",
			paramKey:   "id",
			expectOK:   true,
			expectCode: http.StatusOK,
			expectStr:  "valid-string",
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
			resultStr, ok := handlers.U().ParseParamString(w, r, tt.paramKey)

			// Check the results
			assert.Equal(t, tt.expectOK, ok)
			assert.Equal(t, tt.expectCode, w.Code)
			if tt.expectOK {
				assert.Equal(t, tt.expectStr, resultStr)
			}
		})
	}
}

// TestParseParamLanguage tests the ParseParamLanguage function
func TestParseParamLanguage(t *testing.T) {
	// Define data
	defaultLanguage := language.English

	// Replace the global utils with a new instance
	handlers.ReplaceGlobals(handlers.NewUtils())

	// Mock DEFAULT_LANGUAGE
	viper.Set("DEFAULT_LANGUAGE", "en")

	// Define the test cases
	tests := []struct {
		name       string
		langParam  string
		expectOK   bool
		expectCode int
		expectLang language.Tag
	}{
		{
			name:       "missing language parameter",
			langParam:  "",
			expectLang: defaultLanguage,
		},
		{
			name:       "invalid language parameter",
			langParam:  "invalid-lang",
			expectLang: defaultLanguage,
		},
		{
			name:       "valid language parameter",
			langParam:  "fr",
			expectLang: language.MustParse("fr"),
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			// Add the language parameter to the request
			q := r.URL.Query()
			q.Add("lang", tt.langParam)
			r.URL.RawQuery = q.Encode()

			// Call the function
			resultLang := handlers.U().ParseParamLanguage(w, r)

			// Check the results
			assert.Equal(t, tt.expectLang, resultLang)
		})
	}
}

// TestParseParamBool tests the ParseParamBool function
func TestParseParamBool(t *testing.T) {
	// Replace the global utils with a new instance
	handlers.ReplaceGlobals(handlers.NewUtils())

	// Define the test cases
	tests := []struct {
		name       string
		paramValue string
		paramKey   string
		expectOK   bool
		expectCode int
		expectBool bool
	}{
		{
			name:       "missing bool",
			paramValue: "",
			paramKey:   "flag",
			expectOK:   false,
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "invalid bool",
			paramValue: "invalid",
			paramKey:   "flag",
			expectOK:   false,
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "valid true bool",
			paramValue: "true",
			paramKey:   "flag",
			expectOK:   true,
			expectCode: http.StatusOK,
			expectBool: true,
		},
		{
			name:       "valid false bool",
			paramValue: "false",
			paramKey:   "flag",
			expectOK:   true,
			expectCode: http.StatusOK,
			expectBool: false,
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
			resultBool, ok := handlers.U().ParseParamBool(w, r, tt.paramKey)

			// Check the results
			assert.Equal(t, tt.expectOK, ok)
			assert.Equal(t, tt.expectCode, w.Code)
			if tt.expectOK {
				assert.Equal(t, tt.expectBool, resultBool)
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

// TestReadImage tests the ReadImage function
func TestReadImage(t *testing.T) {
	// Replace the global utils with a new instance
	handlers.ReplaceGlobals(handlers.NewUtils())

	// Define the test cases
	tests := []struct {
		name        string
		fileContent []byte
		fileName    string
		fileType    string
		expectOK    bool
		expectCode  int
		expectData  []byte
		expectName  string
	}{
		{
			name:       "missing file",
			expectOK:   false,
			expectCode: http.StatusBadRequest,
		},
		{
			name:        "invalid MIME type",
			fileContent: []byte("invalid content"),
			fileName:    "invalid.txt",
			fileType:    "text/plain",
			expectOK:    false,
			expectCode:  http.StatusBadRequest,
		},
		{
			name:        "valid JPEG image",
			fileContent: []byte{0xFF, 0xD8, 0xFF}, // JPEG header
			fileName:    "image.jpg",
			fileType:    "image/jpeg",
			expectOK:    true,
			expectCode:  http.StatusOK,
			expectData:  []byte{0xFF, 0xD8, 0xFF},
			expectName:  "image.jpg",
		},
		{
			name:        "valid PNG image",
			fileContent: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, // PNG header
			fileName:    "image.png",
			fileType:    "image/png",
			expectOK:    true,
			expectCode:  http.StatusOK,
			expectData:  []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			expectName:  "image.png",
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", nil)

			if tt.fileContent != nil {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", tt.fileName)
				_, err := part.Write(tt.fileContent)
				if err != nil {
					assert.Fail(t, "failed to decode response")
				}
				writer.Close()
				r = httptest.NewRequest("POST", "/", body)
				r.Header.Set("Content-Type", writer.FormDataContentType())
			}

			// Call the function
			data, name, ok := handlers.U().ReadImage(w, r)

			// Check the results
			assert.Equal(t, tt.expectOK, ok)
			assert.Equal(t, tt.expectCode, w.Code)
			if tt.expectOK {
				assert.Equal(t, tt.expectData, data)
				assert.Equal(t, tt.expectName, name)
			}
		})
	}
}
