package main

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/auth"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)
import "github.com/Zapharaos/fihub-backend/internal/router"

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

	// Setup application
	app.Init()

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

	// Create server
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

	// Start the server in a separate goroutine
	go func() {
		if serverEnableTLS {
			// Start the server with TLS
			err := srv.ListenAndServeTLS(serverTLSCert, serverTLSKey)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				zap.L().Fatal("server listen with TLS", zap.Error(err))
			}
		} else {
			// Start the server without TLS
			err := srv.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				zap.L().Fatal("server listen", zap.Error(err))
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

	// Gracefully shut down the server
	if err := srv.Shutdown(ctxShutDown); err != nil {
		zap.L().Fatal("Server shutdown failed", zap.Error(err))
	}

	zap.L().Info("Server shutdown")
}
