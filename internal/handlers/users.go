package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/handlers/render"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

// CreateUser godoc
//
//	@Id				CreateUser
//
//	@Summary		Create a new user
//	@Description	Create a new user.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body	users.UserWithPassword	true	"user (json)"
//	@Security		Bearer
//	@Success		200	{object}	users.User				"user"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var userWithPassword users.UserWithPassword
	err := json.NewDecoder(r.Body).Decode(&userWithPassword)
	if err != nil {
		zap.L().Warn("User json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate user
	if ok, err := userWithPassword.IsValid(); !ok {
		zap.L().Warn("User is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Verify user existence
	exists, err := users.R().Exists(userWithPassword.Email)
	if err != nil {
		zap.L().Error("Check user exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		zap.L().Warn("User already exists", zap.String("email", userWithPassword.Email))
		render.BadRequest(w, r, fmt.Errorf("email-used"))
		return
	}

	// Create user
	userID, err := users.R().Create(&userWithPassword)
	if err != nil {
		zap.L().Error("PostUser.Create", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get user back from database
	user, err := users.R().Get(userID)
	if err != nil {
		zap.L().Error("Cannot get user", zap.String("uuid", userID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		zap.L().Error("User not found after creation", zap.String("uuid", userID.String()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, user)
}

// GetUser godoc
//
//	@Id				GetUser
//
//	@Summary		Get a user
//	@Description	Get a user by id.
//	@Tags			Users
//	@Produce		json
//	@Param			id	path	string	true	"user id"
//	@Security		Bearer
//	@Success		200	{object}	users.User				"user"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		404	{string}	string					"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {

	// Retrieve userID from URL
	id := chi.URLParam(r, "id")
	userID, err := uuid.Parse(id)
	if err != nil {
		zap.L().Warn("Parse user id", zap.Error(err))
		render.BadRequest(w, r, fmt.Errorf("invalid user id"))
		return
	}

	// Retrieve user from database
	user, err := users.R().Get(userID)
	if err != nil {
		zap.L().Error("Cannot load user", zap.String("uuid", userID.String()), zap.Error(err))
		render.Error(w, r, nil, "")
		return
	} else if user == nil {
		zap.L().Debug("User not found", zap.String("uuid", userID.String()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Marshal user to JSON
	json.NewEncoder(w).Encode(user)
}

// GetUserSelf godoc
//
//	@Id				GetUserSelf
//
//	@Summary		Get the currently authenticated user
//	@Description	Retrieves the currently authenticated user.
//	@Tags			Users
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{object}	users.User				"user"
//	@Failure		400	{string}	string					"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/me [get]
func GetUserSelf(w http.ResponseWriter, r *http.Request) {
	userCtx, found := GetUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, userCtx)
}
