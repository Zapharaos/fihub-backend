package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/cmd/user/app/service"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/grpcutil"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"google.golang.org/grpc"
	"time"
)

func main() {

	// Setup Environment
	err := app.InitConfiguration("user")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	// Setup gRPC microservice
	serviceName := "USER"
	lis, err := grpcutil.SetupServer(serviceName)
	if err != nil {
		return
	}

	// Setup gRPC clients
	securityConn := grpcutil.ConnectToClient("SECURITY")
	publicSecurityClient := securitypb.NewPublicSecurityServiceClient(securityConn)
	security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))

	// Register gRPC service
	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, &service.Service{})

	// Setup Database
	if app.InitPostgres() {
		setupPostgresRepositories()
	}

	// Start databases health monitoring
	database.StartHealthMonitoring("Postgres", 30*time.Second, database.DB().Postgres(), func() {
		if app.InitPostgres() {
			setupPostgresRepositories()
		}
	})

	// Register gRPC health service
	grpcutil.RegisterHealthServer(s, 30*time.Second, serviceName, serverHealthStatusIsHealthy)

	// Start gRPC server
	grpcutil.StartServer(s, lis, serviceName)
}

// setupPostgresRepositories initializes the Postgres repositories for the microservice.
func setupPostgresRepositories() {
	repositories.ReplaceGlobals(repositories.NewPostgresRepository(database.DB().Postgres().DB))
}

// serverHealthStatusIsHealthy indicates whether the server is healthy.
func serverHealthStatusIsHealthy() bool {
	return database.DB().Postgres().IsHealthy()
}
