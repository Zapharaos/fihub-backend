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

	// Setup Database
	app.InitRedis()

	// TODO : remove once auth fully migrated to redis
	if app.InitPostgres() {
		setupPostgresRepositories()
	}

	// Start databases health monitoring
	// TODO : remove once auth fully migrated to redis
	healthMonitor := database.NewHealthMonitor(30 * time.Second)
	healthMonitor.AddTarget("Postgres", database.DB().Postgres(), func() {
		if app.InitPostgres() {
			setupPostgresRepositories()
		}
	})
	healthMonitor.AddTarget("Redis", database.DB().Redis(), func() {
		app.InitRedis()
	})
	healthMonitor.Start()
	// TODO : uncomment once auth fully migrated to redis
	/*database.StartHealthMonitoring("Redis", 30*time.Second, database.DB().Redis(), func() {
		app.InitRedis()
	})*/

	// Register gRPC health service
	grpcutil.RegisterHealthServer(s, 30*time.Second, serviceName, serverHealthStatusIsHealthy)

	// Start gRPC server
	grpcutil.StartServer(s, lis, serviceName)
}

// setupPostgresRepositories initializes the Postgres repositories for the microservice.
func setupPostgresRepositories() {
	// TODO : remove once auth fully migrated to redis
	userrepositories.ReplaceGlobals(userrepositories.NewPostgresRepository(database.DB().Postgres().DB))
	password.ReplaceGlobals(password.NewPostgresRepository(database.DB().Postgres().DB))
}

// serverHealthStatusIsHealthy indicates whether the server is healthy.
func serverHealthStatusIsHealthy() bool {
	return database.DB().Postgres().IsHealthy() &&
		database.DB().Redis().IsHealthy()
}
