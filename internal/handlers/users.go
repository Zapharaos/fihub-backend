package handlers

import (
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var userWithPassword users.UserWithPassword
	err := json.NewDecoder(r.Body).Decode(&userWithPassword)
	if err != nil {
		zap.L().Warn("User json decode", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate user
	if ok, err := userWithPassword.IsValid(); !ok {
		zap.L().Warn("User is not valid", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		zap.L().Warn("User already exists", zap.String("login", userWithPassword.Email))
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// Create user
	userID, err := users.R().Create(&userWithPassword)
	if err != nil {
		zap.L().Error("PostUser.Create", zap.Error(err))
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
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

	// Marshal user to JSON
	json.NewEncoder(w).Encode(user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	// Retrieve userID from URL
	id := chi.URLParam(r, "id")
	userID, err := uuid.Parse(id)
	if err != nil {
		zap.L().Warn("Parse user id", zap.Error(err))
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Retrieve user from database
	user, err := users.R().Get(userID)
	if err != nil {
		zap.L().Error("Cannot load user", zap.String("uuid", userID.String()), zap.Error(err))
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	} else if user == nil {
		zap.L().Debug("User not found", zap.String("uuid", userID.String()))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Marshal user to JSON
	json.NewEncoder(w).Encode(user)
}
