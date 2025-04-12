package health

import (
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/health"
	genhealth "github.com/Zapharaos/fihub-backend/protogen/health"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func main() {

	// Setup Environment
	app.InitConfiguration()

	// Setup Logger
	app.InitLogger()

	// Start gRPC server
	port := viper.GetString("HEALTH_MICROSERVICE_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		zap.L().Error("Failed to listen Health microservice: %v", zap.Error(err))
	}

	s := grpc.NewServer()

	// Register your gRPC service here
	genhealth.RegisterHealthServiceServer(s, &health.Service{})

	zap.L().Info("gRPC Health microservice is running on port : " + port)
	if err := s.Serve(lis); err != nil {
		zap.L().Error("Failed to serve health microservice: %v", zap.Error(err))
	}
}
