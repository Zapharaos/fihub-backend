package handlers

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/gen/go/healthpb"
	"log"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call the health check method on the gRPC client
	_, err := clients.C().Health().CheckHealth(ctx, &healthpb.HealthRequest{
		ServiceName: "fihub-backend",
	})
	if err != nil {
		log.Printf("Health check failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the response
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("Unable to write response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
