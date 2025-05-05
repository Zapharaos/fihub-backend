package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/auth/app/service"
	userrepositories "github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/grpcutil"
	"github.com/Zapharaos/fihub-backend/internal/password"
	"google.golang.org/grpc"
	"time"
)

func main() {

	// Setup Environment
	err := app.InitConfiguration("auth")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	// Setup gRPC microservice
	serviceName := "AUTH"
	lis, err := grpcutil.SetupServer(serviceName)
	if err != nil {
		return
	}

	// Setup gRPC clients
	userConn := grpcutil.ConnectToClient("USER")
	userClient := userpb.NewUserServiceClient(userConn)

	// Register gRPC service
	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, service.NewAuthService(userClient))

	// TODO : remove
	// Setup Database
	if app.ConnectPostgres() {
		setupPostgresRepositories()
	}

	// TODO : remove
	// Maintain postgres health status
	app.StartPostgresHealthCheck(30*time.Second, setupPostgresRepositories)

	// Register gRPC health service
	grpcutil.RegisterHealthServer(s, 30*time.Second, serviceName, serverHealthStatusIsHealthy)

	// Start gRPC server
	grpcutil.StartServer(s, lis, serviceName)
}

// setupPostgresRepositories initializes the Postgres repositories for the microservice.
func setupPostgresRepositories() {
	// TODO : remove
	userrepositories.ReplaceGlobals(userrepositories.NewPostgresRepository(database.DB().Postgres().DB))
	password.ReplaceGlobals(password.NewPostgresRepository(database.DB().Postgres().DB))
}

// serverHealthStatusIsHealthy indicates whether the server is healthy.
func serverHealthStatusIsHealthy() bool {
	// TODO : remove
	return database.DB().Postgres().IsHealthy()
}
