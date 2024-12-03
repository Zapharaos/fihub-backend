package main

import (
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"log"
	"net/http"
	"time"
)
import "github.com/Zapharaos/fihub-backend/internal/router"

func main() {

	app.Init()

	// Create router
	r := router.New()

	// Create server
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
