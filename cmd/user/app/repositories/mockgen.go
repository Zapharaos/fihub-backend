package repositories

//go:generate mockgen -source=repository.go -destination=../../../../test/mocks/user_repository.go --package=mocks -mock_names=Repository=UserRepository Repository
