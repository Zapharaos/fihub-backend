package main

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)
import "github.com/Zapharaos/fihub-backend/internal/router"

func main() {

	// Setup application
	app.Init()

	// Create router
	r := router.New()

	// Create server
	srv := &http.Server{
		Addr:         ":" + env.GetString("GO_PORT", "8080"),
		Handler:      r,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// Create a channel to receive termination signals (SIGINT, SIGTERM)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a separate goroutine
	go func() {
		// Listen for and handle incoming requests
		err := srv.ListenAndServe()

		// Check if the error is due to a graceful shutdown
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("server listen", zap.Error(err))
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
