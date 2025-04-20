package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/repositories"
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/service"
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
	err := app.InitConfiguration("broker")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	// Setup Database
	app.InitDatabase()

	// Broker repositories
	brokerRepository := repositories.NewPostgresRepository(database.DB().Postgres())
	userBrokerRepository := repositories.NewUserPostgresRepository(database.DB().Postgres())
	imageBrokerRepository := repositories.NewImagePostgresRepository(database.DB().Postgres())
	repositories.ReplaceGlobals(repositories.NewRepository(brokerRepository, userBrokerRepository, imageBrokerRepository))

	// Start gRPC microservice
	port := viper.GetString("BROKER_MICROSERVICE_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		zap.L().Error("Failed to listen Broker microservice: %v", zap.Error(err))
	}

	s := grpc.NewServer()

	// Register gRPC service
	protogen.RegisterBrokerServiceServer(s, &service.Service{})

	zap.L().Info("gRPC Broker microservice is running on port : " + port)
	if err := s.Serve(lis); err != nil {
		zap.L().Error("Failed to serve Broker microservice: %v", zap.Error(err))
	}
}
