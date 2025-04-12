package health

import (
	"context"
	"github.com/Zapharaos/fihub-backend/protogen/health"
)

// Service is the implementation of the HealthService interface.
type Service struct {
	health.UnimplementedHealthServiceServer
}

// CheckHealth implements the CheckHealth RPC method.
func (h *Service) CheckHealth(ctx context.Context, req *health.HealthRequest) (*health.HealthResponse, error) {
	// Example logic for health check
	if req.ServiceName == "" {
		return &health.HealthResponse{
			IsHealthy: false,
			Message:   "Service name is required",
		}, nil
	}

	return &health.HealthResponse{
		IsHealthy: true,
		Message:   "Service is healthy",
	}, nil
}
