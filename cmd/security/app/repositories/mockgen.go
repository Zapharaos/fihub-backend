package repositories

//go:generate mockgen -source=role_repository.go -destination=../../../../test/mocks/security_role_repository.go --package=mocks -mock_names=RoleRepository=SecurityRoleRepository RoleRepository
//go:generate mockgen -source=permission_repository.go -destination=../../../../test/mocks/security_permission_repository.go --package=mocks -mock_names=PermissionRepository=SecurityPermissionRepository PermissionRepository
