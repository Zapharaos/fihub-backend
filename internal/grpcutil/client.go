package grpcutil

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

// GetClientAddress retrieves the address of a gRPC service based on its name
func GetClientAddress(serviceName string) string {
	host := viper.GetString(fmt.Sprintf("%s_MICROSERVICE_HOST", serviceName))
	port := viper.GetString(fmt.Sprintf("%s_MICROSERVICE_PORT", serviceName))
	return fmt.Sprintf("%s:%s", host, port)
}

// ConnectToClient creates a gRPC client connection based on service name and returns the connection
func ConnectToClient(serviceName string) *grpc.ClientConn {
	address := GetClientAddress(serviceName)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Fatal("Failed to connect to gRPC service", zap.String("service", serviceName), zap.Error(err))
	}
	zap.L().Info("Connected to gRPC service", zap.String("service", serviceName), zap.String("address", conn.Target()))

	return conn
}

// CheckClientHealth checks the health of a gRPC service using the health check protocol
func CheckClientHealth(conn *grpc.ClientConn, serviceName string) (*grpc_health_v1.HealthCheckResponse, error) {
	healthClient := grpc_health_v1.NewHealthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: serviceName,
	})
}
