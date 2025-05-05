package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/transaction/app/repositories"
	"github.com/Zapharaos/fihub-backend/cmd/transaction/app/service"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/gen/go/transactionpb"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/grpcutil"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func main() {

	// Setup Environment
	err := app.InitConfiguration("transaction")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	defer app.RecoverPanic()   // Catch and log panics
	defer app.CleanResources() // Clean up regardless of shutdown cause

	// Setup gRPC microservice
	serviceName := "TRANSACTION"
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
	transactionpb.RegisterTransactionServiceServer(s, &service.Service{})

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
	<-grpcutil.WaitForShutdown()

	// Shutdown
	zap.L().Info("Shutdown gRPC server", zap.String("service", serviceName))
	s.GracefulStop() // Stop server cleanly
}

// setupPostgresRepositories initializes the Postgres repositories for the microservice.
func setupPostgresRepositories() {
	repositories.ReplaceGlobals(repositories.NewPostgresRepository(database.DB().Postgres().DB))
}

// serverHealthStatusIsHealthy indicates whether the server is healthy.
func serverHealthStatusIsHealthy() bool {
	return database.DB().Postgres().IsHealthy()
}
