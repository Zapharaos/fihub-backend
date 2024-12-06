package handlers

import (
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"go.uber.org/zap"
	"net/http"
)

// GetUserFromContext extract the logged user from the request context
func GetUserFromContext(r *http.Request) (users.User, bool) {
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
