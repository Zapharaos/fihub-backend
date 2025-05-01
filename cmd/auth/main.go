package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/auth/app/service"
	"github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func main() {

	// Setup Environment
	err := app.InitConfiguration("auth")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	// Setup Database
	app.InitDatabase()

	// TODO : remove this
	// User repositories
	repositories.ReplaceGlobals(repositories.NewPostgresRepository(database.DB().Postgres()))

	// Start gRPC microservice
	port := viper.GetString("AUTH_MICROSERVICE_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		zap.L().Error("Failed to listen Auth microservice: %v", zap.Error(err))
	}

	s := grpc.NewServer()

	// Register gRPC service
	protogen.RegisterAuthServiceServer(s, service.NewAuthService())

	zap.L().Info("gRPC Auth microservice is running on port : " + port)
	if err := s.Serve(lis); err != nil {
		zap.L().Error("Failed to serve Auth microservice: %v", zap.Error(err))
	}
}
