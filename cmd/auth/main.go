package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/auth/app/service"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/grpcutil"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/Zapharaos/fihub-backend/pkg/translation"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/language"
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

	defer app.RecoverPanic()   // Catch and log panics
	defer app.CleanResources() // Clean up regardless of shutdown cause

	// Setup Email
	email.ReplaceGlobals(email.NewSendgridService())

	// Setup Translations
	defaultLang := language.MustParse(viper.GetString("DEFAULT_LANGUAGE"))
	translation.ReplaceGlobals(translation.NewI18nService(defaultLang))

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

	// Start databases health monitoring
	database.StartHealthMonitoring("Redis", 30*time.Second, database.DB().Redis(), func() {
		app.InitRedis()
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

// serverHealthStatusIsHealthy indicates whether the server is healthy.
func serverHealthStatusIsHealthy() bool {
	return database.DB().Postgres().IsHealthy() &&
		database.DB().Redis().IsHealthy()
}
