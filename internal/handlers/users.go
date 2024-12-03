package handlers

import (
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user users.UserWithPassword
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create user
	_, err = users.R().Create(&user)
	if err != nil {
		log.Println("User.Create:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := users.R().Get(userID)
	if err != nil {
		log.Println("User.Get:", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	} else if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Marshal user to JSON
	json.NewEncoder(w).Encode(user)
}
