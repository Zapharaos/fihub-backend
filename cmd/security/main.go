package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/service"
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
	err := app.InitConfiguration("security")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	// Setup Database
	app.InitDatabase()

	// User repositories
	roleRepository := repositories.NewRolePostgresRepository(database.DB().Postgres())
	permissionRepository := repositories.NewPermissionPostgresRepository(database.DB().Postgres())
	repositories.ReplaceGlobals(repositories.NewRepository(roleRepository, permissionRepository))

	// Start gRPC microservice
	port := viper.GetString("SECURITY_MICROSERVICE_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		zap.L().Error("Failed to listen Security microservice: %v", zap.Error(err))
	}

	s := grpc.NewServer()

	// Register gRPC service
	protogen.RegisterSecurityServiceServer(s, &service.Service{})

	zap.L().Info("gRPC Security microservice is running on port : " + port)
	if err := s.Serve(lis); err != nil {
		zap.L().Error("Failed to serve Security microservice: %v", zap.Error(err))
	}
}
