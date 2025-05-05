package security

//go:generate mockgen -source=facade_public.go -destination=../../test/mocks/security_facade_public.go --package=mocks -mock_names=Repository=PasswordRepository Repository
