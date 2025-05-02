package repositories

//go:generate mockgen -source=broker_repository.go -destination=../../../../test/mocks/broker_repository.go --package=mocks -mock_names=BrokerRepository=BrokerRepository BrokerRepository
//go:generate mockgen -source=image_repository.go -destination=../../../../test/mocks/broker_repository_image.go --package=mocks -mock_names=ImageRepository=BrokerImageRepository ImageRepository
//go:generate mockgen -source=user_repository.go -destination=../../../../test/mocks/broker_repository_user.go --package=mocks -mock_names=UserRepository=BrokerUserRepository UserRepository
