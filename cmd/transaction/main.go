package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/transaction/app/transaction"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/database"
	gentransaction "github.com/Zapharaos/fihub-backend/protogen/transaction"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func main() {

	// Setup Environment
	err := app.InitConfiguration("transaction")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	// Setup Database
	app.InitDatabase()

	// Transactions Repository
	transaction.ReplaceGlobals(transaction.NewPostgresRepository(database.DB().Postgres()))

	// Start gRPC microservice
	port := viper.GetString("TRANSACTION_MICROSERVICE_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		zap.L().Error("Failed to listen Transaction microservice: %v", zap.Error(err))
	}

	s := grpc.NewServer()

	// Register your gRPC service here
	gentransaction.RegisterTransactionServiceServer(s, &transaction.Service{})

	zap.L().Info("gRPC Transaction microservice is running on port : " + port)
	if err := s.Serve(lis); err != nil {
		zap.L().Error("Failed to serve Transaction microservice: %v", zap.Error(err))
	}
}
