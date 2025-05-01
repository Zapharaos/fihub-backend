package app

type keyContext string

const (
	// ContextKeyUserID is used as key to add the user ID in the request context
	ContextKeyUserID keyContext = "user"
)
