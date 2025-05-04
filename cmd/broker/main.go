package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/repositories"
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/service"
	"github.com/Zapharaos/fihub-backend/gen/go/brokerpb"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/grpcutil"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"google.golang.org/grpc"
	"time"
)

func main() {

	// Setup Environment
	err := app.InitConfiguration("broker")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	// Setup gRPC microservice
	serviceName := "BROKER"
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
	brokerpb.RegisterBrokerServiceServer(s, &service.Service{})

	// Setup Database
	if app.ConnectPostgres() {
		setupPostgresRepositories()
	}

	// Maintain postgres health status
	app.StartPostgresHealthCheck(30*time.Second, setupPostgresRepositories)

	// Register gRPC health service
	grpcutil.RegisterHealthServer(s, 30*time.Second, serviceName, serverHealthStatusIsHealthy)

	// Start gRPC server
	grpcutil.StartServer(s, lis, serviceName)
}

// setupPostgresRepositories initializes the Postgres repositories for the microservice.
func setupPostgresRepositories() {
	brokerRepository := repositories.NewPostgresRepository(database.DB().Postgres().DB)
	userBrokerRepository := repositories.NewUserPostgresRepository(database.DB().Postgres().DB)
	imageBrokerRepository := repositories.NewImagePostgresRepository(database.DB().Postgres().DB)
	repositories.ReplaceGlobals(repositories.NewRepository(brokerRepository, userBrokerRepository, imageBrokerRepository))
}

// serverHealthStatusIsHealthy indicates whether the server is healthy.
func serverHealthStatusIsHealthy() bool {
	return database.DB().Postgres().IsHealthy()
}
