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

	// Init application
	app.Init()
	defer app.Stop()

	// TODO : Configure

	// Router
	r := router.New()

	// Server
	srv := &http.Server{
		Addr:         ":" + env.GetString("GO_PORT", "8080"),
		Handler:      r,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Println("Listening on " + srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
