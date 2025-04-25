package password

//go:generate mockgen -source=repository.go -destination=../../test/mocks/password_repository.go --package=mocks -mock_names=Repository=PasswordRepository Repository
