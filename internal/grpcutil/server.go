package grpcutil

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// SetupServer creates a new gRPC server listener for the specified service
func SetupServer(serviceName string) (net.Listener, error) {
	port := viper.GetString(fmt.Sprintf("%s_MICROSERVICE_PORT", serviceName))
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		zap.L().Error("Failed to listen: %v", zap.String("service", serviceName), zap.Error(err))
		return nil, err
	}
	return lis, nil
}

// StartServer starts the gRPC server and listens for incoming connections
func StartServer(s *grpc.Server, lis net.Listener, serviceName string) {
	go func() {
		// Starting the microservice server
		if err := s.Serve(lis); err != nil {
			zap.L().Error("Failed to serve microservice: %v", zap.Error(err))
		}
	}()
	zap.L().Info("Served gRPC microservice", zap.String("service", serviceName), zap.String("address", lis.Addr().String()))
}

// WaitForShutdown waits for an OS signal
func WaitForShutdown() <-chan os.Signal {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return done
}

// RegisterHealthServer registers a health check service with the gRPC server
func RegisterHealthServer(s *grpc.Server, interval time.Duration, serviceName string, isHealthy func() bool) {
	// Register health service
	healthServer := health.NewServer()
	healthServer.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	// Maintain server health status
	maintainServerHealthStatus(healthServer, serviceName, isHealthy)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			maintainServerHealthStatus(healthServer, serviceName, isHealthy)
		}
	}()
}

func maintainServerHealthStatus(healthServer *health.Server, serviceName string, isHealthy func() bool) {
	// Check postgres health
	if isHealthy() {
		healthServer.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_SERVING)
	} else {
		healthServer.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		zap.L().Error("Postgres health check failed")
	}
}
