package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/protogen/health"
	"go.uber.org/zap"
)

// Service is the implementation of the HealthService interface.
type Service struct {
	health.UnimplementedHealthServiceServer
}

// CheckHealth implements the CheckHealth RPC method.
func (h *Service) CheckHealth(ctx context.Context, req *health.HealthRequest) (*health.HealthResponse, error) {

	zap.L().Info("Checking service", zap.String("service_name", req.ServiceName))

	// TODO : check global health status or specific service health status

	// Example logic for service check
	if req.ServiceName == "" {
		zap.L().Error("Service name is required")
		return &health.HealthResponse{
			IsHealthy: false,
			Message:   "Service name is required",
		}, nil
	}

	zap.L().Info("Service is healthy", zap.String("service_name", req.ServiceName))

	return &health.HealthResponse{
		IsHealthy: true,
		Message:   "Service is healthy",
	}, nil
}
