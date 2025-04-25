package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/service"
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
	err := app.InitConfiguration("user")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	// Setup Database
	app.InitDatabase()

	// User repositories
	userRepository := repositories.NewUserPostgresRepository(database.DB().Postgres())
	roleRepository := repositories.NewRolePostgresRepository(database.DB().Postgres())
	permissionRepository := repositories.NewPermissionPostgresRepository(database.DB().Postgres())
	repositories.ReplaceGlobals(repositories.NewRepository(userRepository, roleRepository, permissionRepository))

	// Start gRPC microservice
	port := viper.GetString("USER_MICROSERVICE_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		zap.L().Error("Failed to listen User microservice: %v", zap.Error(err))
	}

	s := grpc.NewServer()

	// Register gRPC service
	protogen.RegisterBrokerServiceServer(s, &service.Service{})

	zap.L().Info("gRPC User microservice is running on port : " + port)
	if err := s.Serve(lis); err != nil {
		zap.L().Error("Failed to serve User microservice: %v", zap.Error(err))
	}
}
