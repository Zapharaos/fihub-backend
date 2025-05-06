package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/service"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
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
	err := app.InitConfiguration("security")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	defer app.RecoverPanic()   // Catch and log panics
	defer app.CleanResources() // Clean up regardless of shutdown cause

	// Setup gRPC microservice
	serviceName := "SECURITY"
	lis, err := grpcutil.SetupServer(serviceName)
	if err != nil {
		return
	}

	// Register gRPC services
	s := grpc.NewServer()
	publicService := &service.PublicService{}
	securitypb.RegisterPublicSecurityServiceServer(s, publicService)
	securitypb.RegisterSecurityServiceServer(s, &service.Service{})
	security.ReplaceGlobals(security.NewPublicSecurityFacade(publicService))

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
	roleRepository := repositories.NewRolePostgresRepository(database.DB().Postgres().DB)
	permissionRepository := repositories.NewPermissionPostgresRepository(database.DB().Postgres().DB)
	repositories.ReplaceGlobals(repositories.NewRepository(roleRepository, permissionRepository))
}

// serverHealthStatusIsHealthy indicates whether the server is healthy.
func serverHealthStatusIsHealthy() bool {
	return database.DB().Postgres().IsHealthy()
}
