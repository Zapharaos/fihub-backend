package handlers

import (
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/handlers/render"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

// getUserFromContext extract the logged user from the request context
func getUserFromContext(r *http.Request) (users.User, bool) {
	_user := r.Context().Value(app.ContextKeyUser)
	if _user == nil {
		zap.L().Warn("No context user provided")
		return users.User{}, false
	}
	user, ok := _user.(users.User)
	if !ok {
		zap.L().Warn("Invalid user type in context")
		return users.User{}, false
	}
	return user, true
}

// ParseParamUUID parses an uuid from the request parameters (using key parameter)
func ParseParamUUID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	value := chi.URLParam(r, key)

	result, err := uuid.Parse(value)
	if err != nil {
		zap.L().Debug("Parse uuid", zap.String("key", key), zap.Error(err))
		render.BadRequest(w, r, fmt.Errorf("invalid %s", key))
		return uuid.UUID{}, false
	}

	return result, true
}

// parseUUIDPair is a helper function to parse a key and base UUIDs from the request
// using the key "id" for the base UUID
func parseUUIDPair(w http.ResponseWriter, r *http.Request, key string) (baseID, keyID uuid.UUID, ok bool) {
	keyID, ok = ParseParamUUID(w, r, key)
	if !ok {
		return
	}
	baseID, ok = ParseParamUUID(w, r, "id")
	return
}
