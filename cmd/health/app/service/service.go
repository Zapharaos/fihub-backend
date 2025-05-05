package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/health/app/clients"
	"github.com/Zapharaos/fihub-backend/gen/go/healthpb"
	"github.com/Zapharaos/fihub-backend/internal/grpcutil"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service is the implementation of the HealthService interface.
type Service struct {
	healthpb.UnimplementedHealthServiceServer
}

// CheckHealth implements the CheckHealth RPC method.
func (h *Service) CheckHealth(ctx context.Context, req *healthpb.HealthRequest) (*healthpb.HealthStatus, error) {

	zap.L().Info("Checking service", zap.String("service_name", req.ServiceName))

	// Searches for the service in the clients map
	conn, err := clients.GetTypedClient[*grpc.ClientConn](clients.C(), req.GetServiceName())
	if err != nil {
		zap.L().Error("Service not found", zap.String("service_name", req.ServiceName))
		return &healthpb.HealthStatus{}, status.Error(codes.NotFound, "Service not found")
	}

	// Check the health of the service using the gRPC client
	response, err := grpcutil.CheckClientHealth(conn, req.GetServiceName())
	if err != nil {
		zap.L().Error("Failed to check service health", zap.String("service_name", req.ServiceName), zap.Error(err))
		return &healthpb.HealthStatus{}, status.Error(codes.Internal, "Failed to check service health")
	}

	// Return the health status of the service
	return &healthpb.HealthStatus{
		Status: response.GetStatus().String(),
	}, nil
}
