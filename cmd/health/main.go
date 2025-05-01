package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/health/app/service"
	"github.com/Zapharaos/fihub-backend/gen/go/healthpb"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func main() {

	// Setup Environment
	err := app.InitConfiguration("health")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	// Start gRPC microservice
	port := viper.GetString("HEALTH_MICROSERVICE_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		zap.L().Error("Failed to listen Health microservice: %v", zap.Error(err))
	}

	s := grpc.NewServer()

	// Register gRPC service
	healthpb.RegisterHealthServiceServer(s, &service.Service{})

	zap.L().Info("gRPC Health microservice is running on port : " + port)
	if err := s.Serve(lis); err != nil {
		zap.L().Error("Failed to serve health microservice: %v", zap.Error(err))
	}
}
