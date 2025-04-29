package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/transaction/app/repositories"
	"github.com/Zapharaos/fihub-backend/cmd/transaction/app/service"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/grpcconn"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"github.com/Zapharaos/fihub-backend/protogen"
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
	repositories.ReplaceGlobals(repositories.NewPostgresRepository(database.DB().Postgres()))

	// gRPC clients
	securityConn := grpcconn.ConnectGRPCService("SECURITY")
	publicSecurityClient := protogen.NewPublicSecurityServiceClient(securityConn)
	security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
	
	// Start gRPC microservice
	port := viper.GetString("TRANSACTION_MICROSERVICE_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		zap.L().Error("Failed to listen Transaction microservice: %v", zap.Error(err))
	}

	s := grpc.NewServer()

	// Register gRPC service
	protogen.RegisterTransactionServiceServer(s, &service.Service{})

	zap.L().Info("gRPC Transaction microservice is running on port : " + port)
	if err := s.Serve(lis); err != nil {
		zap.L().Error("Failed to serve Transaction microservice: %v", zap.Error(err))
	}
}
