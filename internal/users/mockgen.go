package users

//go:generate mockgen -source=repository.go -destination=../../test/mocks/users_repository.go --package=mocks -mock_names=Repository=UsersRepository Repository
//go:generate mockgen -source=password/repository.go -destination=../../test/mocks/users_password_repository.go --package=mocks -mock_names=Repository=UsersPasswordRepository Repository
//go:generate mockgen -source=permissions/repository.go -destination=../../test/mocks/permissions_repository.go --package=mocks -mock_names=Repository=PermissionsRepository Repository
//go:generate mockgen -source=roles/repository.go -destination=../../test/mocks/roles_repository.go --package=mocks -mock_names=Repository=RolesRepository Repository
