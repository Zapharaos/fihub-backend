package proto

//go:generate mockgen -source=../gen/go/healthpb/health_grpc.pb.go -destination=../test/mocks/health_client.go -package=mocks HealthServiceClient
//go:generate mockgen -source=../gen/go/userpb/user_grpc.pb.go -destination=../test/mocks/user_client.go -package=mocks UserServiceClient
//go:generate mockgen -source=../gen/go/securitypb/security_grpc.pb.go -destination=../test/mocks/security_client.go -package=mocks SecurityServiceClient
//go:generate mockgen -source=../gen/go/securitypb/security_public_grpc.pb.go -destination=../test/mocks/security_client_public.go -package=mocks PublicSecurityServiceClient
//go:generate mockgen -source=../gen/go/transactionpb/transaction_grpc.pb.go -destination=../test/mocks/transaction_client.go -package=mocks TransactionServiceClient
//go:generate mockgen -source=../gen/go/brokerpb/broker_grpc.pb.go -destination=../test/mocks/broker_client.go -package=mocks BrokerServiceClient
