package clients

//go:generate mockgen -source=../../../../protogen/health_grpc.pb.go -destination=../../../../test/mocks/health_client.go -package=mocks -mock_names=HealthServiceClient=HealthServiceClient HealthServiceClient
//go:generate mockgen -source=../../../../protogen/transaction_grpc.pb.go -destination=../../../../test/mocks/transaction_client.go -package=mocks -mock_names=TransactionServiceClient=TransactionServiceClient TransactionServiceClient
//go:generate mockgen -source=../../../../protogen/broker_grpc.pb.go -destination=../../../../test/mocks/broker_client.go -package=mocks -mock_names=BrokerServiceClient=BrokerServiceClient BrokerServiceClient
