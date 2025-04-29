package main

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/auth"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/router"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/grpcconn"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/Zapharaos/fihub-backend/pkg/translation"
	"github.com/Zapharaos/fihub-backend/protogen"
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

	zap.L().Info("Starting Fihub Backend", zap.String("version", app.Version), zap.String("build_date", app.BuildDate))

	// Server configuration
	serverPort := viper.GetString("HTTP_SERVER_PORT")
	serverEnableTLS := viper.GetBool("HTTP_SERVER_ENABLE_TLS")
	serverTLSCert := viper.GetString("HTTP_SERVER_TLS_FILE_CRT")
	serverTLSKey := viper.GetString("HTTP_SERVER_TLS_FILE_KEY")

	// Auth configuration
	authConfig := auth.Config{
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
	healthConn := grpcconn.ConnectGRPCService("HEALTH")
	userConn := grpcconn.ConnectGRPCService("USER")
	securityConn := grpcconn.ConnectGRPCService("SECURITY")
	brokerConn := grpcconn.ConnectGRPCService("BROKER")
	transactionConn := grpcconn.ConnectGRPCService("TRANSACTION")

	// Create gRPC clients
	healthClient := protogen.NewHealthServiceClient(healthConn)
	userClient := protogen.NewUserServiceClient(userConn)
	securityClient := protogen.NewSecurityServiceClient(securityConn)
	publicSecurityClient := protogen.NewPublicSecurityServiceClient(securityConn)
	brokerClient := protogen.NewBrokerServiceClient(brokerConn)
	transactionClient := protogen.NewTransactionServiceClient(transactionConn)

	// Setup facades
	security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))

	// Initialize the gRPC clients
	clients.ReplaceGlobals(clients.NewClients(
		clients.WithHealthClient(healthClient),
		clients.WithSecurityClient(securityClient),
		clients.WithUserClient(userClient),
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

	// Setup Database
	app.InitDatabase()

	// Initialize the postgres repositories
	app.InitPostgres(database.DB().Postgres())

	// Setup api clients
	initGrpcClients()

	// Setup Email
	email.ReplaceGlobals(email.NewSendgridService())

	// Setup Translations
	defaultLang := language.MustParse(viper.GetString("DEFAULT_LANGUAGE"))
	translation.ReplaceGlobals(translation.NewI18nService(defaultLang))
}
