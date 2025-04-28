package app

type keyContext string

const (
	// ContextKeyUser is used as key to add the user data in the request context
	ContextKeyUser                     keyContext = "user"
	ContextKeyUserRolesWithPermissions keyContext = "userRoles"
)
