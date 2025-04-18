package clients

//go:generate mockgen -source=../../../../protogen/health/health_grpc.pb.go -destination=../../../../test/mocks/health_client.go -package=mocks -mock_names=HealthServiceClient=HealthServiceClient HealthServiceClient
//go:generate mockgen -source=../../../../protogen/transaction/transaction_grpc.pb.go -destination=../../../../test/mocks/transaction_client.go -package=mocks -mock_names=TransactionServiceClient=TransactionServiceClient TransactionServiceClient
