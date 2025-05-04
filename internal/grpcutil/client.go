package grpcutil

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ConnectToClient creates a gRPC client connection based on service name and returns the connection
func ConnectToClient(serviceName string) *grpc.ClientConn {
	host := viper.GetString(fmt.Sprintf("%s_MICROSERVICE_HOST", serviceName))
	port := viper.GetString(fmt.Sprintf("%s_MICROSERVICE_PORT", serviceName))
	address := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Fatal("Failed to connect to gRPC service", zap.String("service", serviceName), zap.Error(err))
	}
	zap.L().Info("Connected to gRPC service", zap.String("service", serviceName), zap.String("address", conn.Target()))

	return conn
}
