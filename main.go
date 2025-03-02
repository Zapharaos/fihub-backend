package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"github.com/adshao/go-binance/v2"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

	// Create a new binance client
	apiKey := env.GetString("BINANCE_API_KEY", "")
	secretKey := env.GetString("BINANCE_API_SECRET", "")
	client := binance.NewClient(apiKey, secretKey)

	// Retrieve SPOT account information
	spot, err := client.NewGetUserAsset().Do(context.Background())
	if err != nil {
		log.Fatalf("Error retrieving SPOT account information: %v", err)
	}

	// Retrieve Simple Earn account information
	earn, err := client.NewSimpleEarnService().FlexibleService().GetPosition().Do(context.Background())
	if err != nil {
		log.Fatalf("Error retrieving Simple Earn account information: %v", err)
	}

	// Process wallet data
	assets := make(map[string]float64)
	for _, row := range earn.Rows {
		amount, err := strconv.ParseFloat(row.TotalAmount, 64)
		if err != nil {
			log.Fatalf("Error parsing amount: %v", err)
		}
		assets[row.Asset] = amount
	}
	for _, row := range spot {
		free, err := strconv.ParseFloat(row.Free, 64)
		if err != nil {
			log.Fatalf("Error parsing amount: %v", err)
		}
		locked, err := strconv.ParseFloat(row.Locked, 64)
		if err != nil {
			log.Fatalf("Error parsing amount: %v", err)
		}
		if amount, exists := assets[row.Asset]; exists {
			// Increment the amount if the asset already exists
			assets[row.Asset] = amount + free + locked
		} else {
			// Add the asset to the map
			assets[row.Asset] = free + locked
		}
	}

	for asset, amount := range assets {
		fmt.Printf("Asset: %s, Amount: %f\n", asset, amount)
	}

	exchange, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatalf("Error retrieving exchange info: %v", err)
	}

	fmt.Printf("Symbols: %+v\n", len(exchange.Symbols))
	/*for _, symbol := range exchange.Symbols {
		fmt.Printf("Symbol: %+v\n", symbol.Symbol)
	}*/

	trades, err := client.NewListTradesService().Symbol("ETHUSDT").Do(context.Background())
	if err != nil {
		log.Println("Error retrieving fiat payments history:", err)
		return
	}

	fmt.Printf("ETH Trades: %+v\n", len(trades))
	/*for _, trade := range trades {
		fmt.Printf("Trade: %+v\n", trade)
	}*/

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
