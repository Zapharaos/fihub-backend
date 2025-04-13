package handlers

//go:generate mockgen -source=utils.go -destination=../../test/mocks/handlers_utils.go --package=mocks Utils

import (
	"errors"
	"fmt"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"io"
	"net/http"
	"strconv"
	"sync"
)

// Utils defines the interface for handler utility functions
type Utils interface {
	CheckPermission(w http.ResponseWriter, r *http.Request, permission string) bool
	GetUserFromContext(r *http.Request) (users.UserWithRoles, bool)
	ParseParamString(w http.ResponseWriter, r *http.Request, key string) (string, bool)
	ParseParamUUID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool)
	ParseParamLanguage(w http.ResponseWriter, r *http.Request) language.Tag
	ParseParamBool(w http.ResponseWriter, r *http.Request, key string) (bool, bool)
	ParseUUIDPair(w http.ResponseWriter, r *http.Request, key string) (baseID, keyID uuid.UUID, ok bool)
	ReadImage(w http.ResponseWriter, r *http.Request) ([]byte, string, bool)
}

var (
	_globalUtilsMu sync.RWMutex
	_globalUtils   Utils
)

// U is used to access the global utils singleton
func U() Utils {
	_globalUtilsMu.RLock()
	defer _globalUtilsMu.RUnlock()

	utils := _globalUtils
	return utils
}

// ReplaceGlobals affect a new utils to the global utils singleton
func ReplaceGlobals(utils Utils) func() {
	_globalUtilsMu.Lock()
	defer _globalUtilsMu.Unlock()

	prev := _globalUtils
	_globalUtils = utils
	return func() { ReplaceGlobals(prev) }
}

type utils struct{}

func NewUtils() Utils {
	u := &utils{}
	var utils Utils = u
	return utils
}

// CheckPermission checks if context user has permission
// if user has no permission, it returns 403 status code
// returns true to indicate that the user has the permission
// returns false to indicate that the user has not the permission and the request should be stopped
func (u *utils) CheckPermission(w http.ResponseWriter, r *http.Request, permission string) bool {
	user, ok := U().GetUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	// check if the user has the permission
	if !user.HasPermission(permission) {
		w.WriteHeader(http.StatusForbidden)
		return false
	}

	return true
}

// GetUserFromContext extract the logged user from the request context
func (u *utils) GetUserFromContext(r *http.Request) (users.UserWithRoles, bool) {
	_user := r.Context().Value(app.ContextKeyUser)
	if _user == nil {
		zap.L().Warn("No context user provided")
		return users.UserWithRoles{}, false
	}
	user, ok := _user.(users.UserWithRoles)
	if !ok {
		zap.L().Warn("Invalid user type in context")
		return users.UserWithRoles{}, false
	}
	return user, true
}

// ParseParamString parses a string from the request parameters (using key parameter)
func (u *utils) ParseParamString(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	value := chi.URLParam(r, key)
	if value == "" {
		zap.L().Debug("Parse string", zap.String("key", key))
		render.BadRequest(w, r, fmt.Errorf("invalid %s", key))
		return "", false
	}

	return value, true
}

// ParseParamUUID parses an uuid from the request parameters (using key parameter)
func (u *utils) ParseParamUUID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	value := chi.URLParam(r, key)

	result, err := uuid.Parse(value)
	if err != nil {
		zap.L().Debug("Parse uuid", zap.String("key", key), zap.Error(err))
		render.BadRequest(w, r, fmt.Errorf("invalid %s", key))
		return uuid.UUID{}, false
	}

	return result, true
}

// ParseParamLanguage parses a language from the request parameters
func (u *utils) ParseParamLanguage(w http.ResponseWriter, r *http.Request) language.Tag {
	langParam := r.URL.Query().Get("lang")
	lang, err := language.Parse(langParam)
	if langParam == "" || err != nil {
		// If no language is provided, use the default language
		defaultLang := viper.GetString("DEFAULT_LANG")
		return language.MustParse(defaultLang)
	}
	return lang
}

// ParseParamBool parses a boolean from the request parameters (using key parameter)
func (u *utils) ParseParamBool(w http.ResponseWriter, r *http.Request, key string) (bool, bool) {
	value := chi.URLParam(r, key)
	if value == "" {
		zap.L().Debug("Parse bool", zap.String("key", key))
		render.BadRequest(w, r, fmt.Errorf("invalid %s", key))
		return false, false
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		zap.L().Debug("Parse bool", zap.String("key", key), zap.Error(err))
		render.BadRequest(w, r, fmt.Errorf("invalid %s", key))
		return false, false
	}

	return result, true
}

// ParseUUIDPair is a helper function to parse a key and base UUIDs from the request
// using the key "id" for the base UUID
func (u *utils) ParseUUIDPair(w http.ResponseWriter, r *http.Request, key string) (baseID, keyID uuid.UUID, ok bool) {
	keyID, ok = U().ParseParamUUID(w, r, key)
	if !ok {
		return
	}
	baseID, ok = U().ParseParamUUID(w, r, "id")
	return
}

// ReadImage reads an image from a multipart form
func (u *utils) ReadImage(w http.ResponseWriter, r *http.Request) ([]byte, string, bool) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		render.BadRequest(w, r, err)
		return nil, "", false
	}

	// Get the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		zap.L().Warn("Form file", zap.Error(err))
		render.BadRequest(w, r, err)
		return nil, "", false
	}
	defer file.Close()

	// Read the file
	data, err := io.ReadAll(file)
	if err != nil {
		zap.L().Warn("Read file", zap.Error(err))
		render.BadRequest(w, r, err)
		return nil, "", false
	}

	// Check the MIME type
	mimeType := http.DetectContentType(data)
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		zap.L().Warn("Invalid MIME type", zap.String("mimeType", mimeType))
		render.BadRequest(w, r, errors.New("invalid-type"))
		return nil, "", false
	}

	return data, header.Filename, true
}
