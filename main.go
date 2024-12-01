package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func main() {

	// Load the .env file in the current directory
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a new router
	r := mux.NewRouter()

	// Blackjack game endpoints
	//r.HandleFunc("/api/games/blackjack", games.BlackjackCreate).Methods("POST")

	// Get the allowed origins
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	// Middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigins},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Create handler
	handler := c.Handler(r)

	// Get the port
	httpPort := os.Getenv("GO_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":"+httpPort, handler))
}
