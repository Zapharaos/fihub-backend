package main

import (
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"log"
	"net/http"
	"time"
)
import "github.com/Zapharaos/fihub-backend/internal/router"

func main() {

	// Load the .env file
	err := env.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Connect to database
	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Stop(db)

	// Create router
	r := router.New()

	// Prepare server
	srv := &http.Server{
		Addr:         ":" + env.GetString("GO_PORT", "8080"),
		Handler:      r,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Println("Server about to listen on " + srv.Addr)

	// Start server
	log.Fatal(srv.ListenAndServe())
}
