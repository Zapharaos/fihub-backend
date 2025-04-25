package repositories

//go:generate mockgen -source=user_repository.go -destination=../../../../test/mocks/user_repository.go --package=mocks -mock_names=UserRepository=UserRepository UserRepository
//go:generate mockgen -source=role_repository.go -destination=../../../../test/mocks/user_role_repository.go --package=mocks -mock_names=RoleRepository=UserRoleRepository RoleRepository
//go:generate mockgen -source=permission_repository.go -destination=../../../../test/mocks/user_permission_repository.go --package=mocks -mock_names=PermissionRepository=UserPermissionRepository PermissionRepository
