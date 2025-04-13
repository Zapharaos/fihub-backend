package handlers

import (
	"context"
	genhealth "github.com/Zapharaos/fihub-backend/protogen/health"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Connect to the gRPC health microservice
	conn, err := grpc.NewClient("health:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Error("Failed to connect to gRPC server", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := genhealth.NewHealthServiceClient(conn)

	// Call the HealthCheck method
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.CheckHealth(ctx, &genhealth.HealthRequest{
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
