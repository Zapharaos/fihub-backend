package main

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/auth"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/router"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/Zapharaos/fihub-backend/pkg/translation"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	// Connect to the gRPC service microservice
	healthHost := viper.GetString("HEALTH_MICROSERVICE_HOST")
	healthPort := viper.GetString("HEALTH_MICROSERVICE_PORT")
	healthConn, err := grpc.NewClient(healthHost+":"+healthPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Fatal("Failed to connect to health gRPC service", zap.Error(err))
	} else {
		zap.L().Info("Connected to health gRPC service", zap.String("address", healthConn.Target()))
	}
	healthClient := protogen.NewHealthServiceClient(healthConn)

	// Connect to the gRPC broker microservice
	brokerHost := viper.GetString("BROKER_MICROSERVICE_HOST")
	brokerPort := viper.GetString("BROKER_MICROSERVICE_PORT")
	brokerConn, err := grpc.NewClient(brokerHost+":"+brokerPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Fatal("Failed to connect to Broker gRPC service", zap.Error(err))
	} else {
		zap.L().Info("Connected to Broker gRPC service", zap.String("address", brokerConn.Target()))
	}
	brokerClient := protogen.NewBrokerServiceClient(brokerConn)

	// Connect to the gRPC transaction microservice
	transactionHost := viper.GetString("TRANSACTION_MICROSERVICE_HOST")
	transactionPort := viper.GetString("TRANSACTION_MICROSERVICE_PORT")
	transactionConn, err := grpc.NewClient(transactionHost+":"+transactionPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Fatal("Failed to connect to Transaction gRPC service", zap.Error(err))
	} else {
		zap.L().Info("Connected to transaction gRPC service", zap.String("address", transactionConn.Target()))
	}
	transactionClient := protogen.NewTransactionServiceClient(transactionConn)

	// Initialize the gRPC clients
	clients.ReplaceGlobals(clients.NewClients(healthClient, brokerClient, transactionClient))
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
