package clients

//go:generate mockgen -source=../../../../protogen/health/health_grpc.pb.go -destination=../../../../test/mocks/health_client.go -package=mocks -mock_names=HealthServiceClient=HealthServiceClient HealthServiceClient
