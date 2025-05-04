package grpcutil

import (
	"github.com/spf13/viper"
	"testing"
)

func TestConnectGRPCService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		viper.Set("TEST_MICROSERVICE_HOST", "localhost")
		viper.Set("TEST_MICROSERVICE_PORT", "50051")

		conn := ConnectToClient("TEST")
		if conn == nil {
			t.Fatal("Expected a valid gRPC connection, got nil")
		}

		if conn.Target() != "localhost:50051" {
			t.Errorf("Expected connection target to be 'localhost:50051', got '%s'", conn.Target())
		}

		_ = conn.Close()
	})
}
