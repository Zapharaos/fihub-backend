package clients

//go:generate mockgen -source=../../../../gen/health_grpc.pb.go -destination=../../../../test/mocks/health_client.go -package=mocks -mock_names=HealthServiceClient=HealthServiceClient HealthServiceClient
//go:generate mockgen -source=../../../../gen/user_grpc.pb.go -destination=../../../../test/mocks/user_client.go -package=mocks -mock_names=UserServiceClient=UserServiceClient UserServiceClient
//go:generate mockgen -source=../../../../gen/security_grpc.pb.go -destination=../../../../test/mocks/security_client.go -package=mocks -mock_names=SecurityServiceClient=SecurityServiceClient SecurityServiceClient
//go:generate mockgen -source=../../../../gen/transaction_grpc.pb.go -destination=../../../../test/mocks/transaction_client.go -package=mocks -mock_names=TransactionServiceClient=TransactionServiceClient TransactionServiceClient
//go:generate mockgen -source=../../../../gen/broker_grpc.pb.go -destination=../../../../test/mocks/broker_client.go -package=mocks -mock_names=BrokerServiceClient=BrokerServiceClient BrokerServiceClient
