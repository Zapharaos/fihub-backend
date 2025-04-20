package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/protogen"
	"go.uber.org/zap"
)

// Service is the implementation of the HealthService interface.
type Service struct {
	protogen.UnimplementedHealthServiceServer
}

// CheckHealth implements the CheckHealth RPC method.
func (h *Service) CheckHealth(ctx context.Context, req *protogen.HealthRequest) (*protogen.HealthResponse, error) {

	zap.L().Info("Checking service", zap.String("service_name", req.ServiceName))

	// TODO : check global health status or specific service health status

	// Example logic for service check
	if req.ServiceName == "" {
		zap.L().Error("Service name is required")
		return &protogen.HealthResponse{
			IsHealthy: false,
			Message:   "Service name is required",
		}, nil
	}

	zap.L().Info("Service is healthy", zap.String("service_name", req.ServiceName))

	return &protogen.HealthResponse{
		IsHealthy: true,
		Message:   "Service is healthy",
	}, nil
}
