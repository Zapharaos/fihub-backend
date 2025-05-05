package main

import (
	"github.com/Zapharaos/fihub-backend/cmd/health/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/health/app/service"
	"github.com/Zapharaos/fihub-backend/gen/go/healthpb"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/grpcutil"
	"google.golang.org/grpc"
)

func main() {

	// Setup Environment
	err := app.InitConfiguration("health")
	if err != nil {
		return
	}

	// Setup Logger
	app.InitLogger()

	// Setup gRPC microservice
	serviceName := "HEALTH"
	lis, err := grpcutil.SetupServer(serviceName)
	if err != nil {
		return
	}

	// Register gRPC connections
	registerGrpcConnections()

	// Register gRPC service
	s := grpc.NewServer()
	healthpb.RegisterHealthServiceServer(s, &service.Service{})

	// Start gRPC server
	grpcutil.StartServer(s, lis, serviceName)
}

// registerGrpcConnections registers the gRPC connections for the health statuses.
func registerGrpcConnections() {
	// Register the gRPC connections
	clients.ReplaceGlobals(clients.NewClients())
	clients.C().Register("USER", grpcutil.ConnectToClient("USER"))
	clients.C().Register("AUTH", grpcutil.ConnectToClient("AUTH"))
	clients.C().Register("SECURITY", grpcutil.ConnectToClient("SECURITY"))
	clients.C().Register("BROKER", grpcutil.ConnectToClient("BROKER"))
	clients.C().Register("TRANSACTION", grpcutil.ConnectToClient("TRANSACTION"))
}
