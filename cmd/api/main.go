package main

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/router"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/server"
	userrepositories "github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/brokerpb"
	"github.com/Zapharaos/fihub-backend/gen/go/healthpb"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/gen/go/transactionpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/grpcutil"
	"github.com/Zapharaos/fihub-backend/internal/password"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/Zapharaos/fihub-backend/pkg/translation"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//	@version		1.0
//	@title			Fihub API Swagger
//	@description	Fihub API Swagger
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Zapharaos
//	@contact.url	https://matthieu-freitag.com/
//	@contact.email	contact@matthieu-freitag.com

// @host      localhost:8080

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
func main() {

	setup()

	defer app.RecoverPanic()   // Catch and log panics
	defer app.CleanResources() // Clean up regardless of shutdown cause

	zap.L().Info("Starting Fihub Backend", zap.String("version", app.Version), zap.String("build_date", app.BuildDate))

	// Server configuration
	serverPort := viper.GetString("HTTP_SERVER_PORT")
	serverEnableTLS := viper.GetBool("HTTP_SERVER_ENABLE_TLS")
	serverTLSCert := viper.GetString("HTTP_SERVER_TLS_FILE_CRT")
	serverTLSKey := viper.GetString("HTTP_SERVER_TLS_FILE_KEY")

	// Auth configuration
	authConfig := server.Config{
		CORS:        viper.GetBool("HTTP_SERVER_API_ENABLE_CORS"),
		Security:    viper.GetBool("HTTP_SERVER_API_ENABLE_SECURITY"),
		GatewayMode: viper.GetBool("HTTP_SERVER_API_ENABLE_GATEWAY_MODE"),
		AllowOrigin: viper.GetString("CORS_ALLOWED_ORIGIN"),
	}

	// Create api
	srv := &http.Server{
		Addr:         ":" + serverPort,
		Handler:      router.New(authConfig),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// Create a channel to receive termination signals (SIGINT, SIGTERM)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the api in a separate goroutine
	go func() {
		if serverEnableTLS {
			// Start the api with TLS
			err := srv.ListenAndServeTLS(serverTLSCert, serverTLSKey)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				zap.L().Fatal("api listen with TLS", zap.Error(err))
			}
		} else {
			// Start the api without TLS
			err := srv.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				zap.L().Fatal("api listen", zap.Error(err))
			}
		}
	}()
	zap.L().Info("Server started", zap.String("addr", srv.Addr))

	// Wait for a termination signal
	<-done

	// Create a context with a 5-second timeout for graceful shutdown
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	// Gracefully shut down the api
	if err := srv.Shutdown(ctxShutDown); err != nil {
		zap.L().Fatal("Server shutdown failed", zap.Error(err))
	}

	zap.L().Info("Server shutdown")
}

// initGrpcClients initializes the gRPC clients for the application.
func initGrpcClients() {
	// Connect to microservices
	healthConn := grpcutil.ConnectToClient("HEALTH")
	userConn := grpcutil.ConnectToClient("USER")
	authConn := grpcutil.ConnectToClient("AUTH")
	securityConn := grpcutil.ConnectToClient("SECURITY")
	brokerConn := grpcutil.ConnectToClient("BROKER")
	transactionConn := grpcutil.ConnectToClient("TRANSACTION")

	// Create gRPC clients
	healthClient := healthpb.NewHealthServiceClient(healthConn)
	userClient := userpb.NewUserServiceClient(userConn)
	authClient := authpb.NewAuthServiceClient(authConn)
	securityClient := securitypb.NewSecurityServiceClient(securityConn)
	publicSecurityClient := securitypb.NewPublicSecurityServiceClient(securityConn)
	brokerClient := brokerpb.NewBrokerServiceClient(brokerConn)
	transactionClient := transactionpb.NewTransactionServiceClient(transactionConn)

	// Setup facades
	security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))

	// Initialize the gRPC clients
	clients.ReplaceGlobals(clients.NewClients(
		clients.WithHealthClient(healthClient),
		clients.WithUserClient(userClient),
		clients.WithAuthClient(authClient),
		clients.WithSecurityClient(securityClient),
		clients.WithBrokerClient(brokerClient),
		clients.WithTransactionClient(transactionClient),
	))
}

func setup() {
	// Setup Environment
	err := app.InitConfiguration("api")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	// Setup api clients
	initGrpcClients()

	// TODO : remove once auth fully migrated

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

	// Setup Email
	email.ReplaceGlobals(email.NewSendgridService())

	// Setup Translations
	defaultLang := language.MustParse(viper.GetString("DEFAULT_LANGUAGE"))
	translation.ReplaceGlobals(translation.NewI18nService(defaultLang))
}

// setupPostgresRepositories initializes the Postgres repositories for the microservice.
func setupPostgresRepositories() {
	// TODO : remove once auth fully migrated
	userrepositories.ReplaceGlobals(userrepositories.NewPostgresRepository(database.DB().Postgres().DB))
	password.ReplaceGlobals(password.NewPostgresRepository(database.DB().Postgres().DB))
}
