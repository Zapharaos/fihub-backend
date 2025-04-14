package clients

import "github.com/Zapharaos/fihub-backend/protogen/health"

type Clients struct {
	health health.HealthServiceClient
}

func NewClients(health health.HealthServiceClient) Clients {
	return Clients{
		health: health,
	}
}

func (c Clients) Health() health.HealthServiceClient {
	return c.health
}

var _globalClients Clients

// C is used to access the global clients singleton
func C() Clients {
	return _globalClients
}

// ReplaceGlobals affect a new clients to the global clients singleton
func ReplaceGlobals(clients Clients) func() {
	prev := _globalClients
	_globalClients = clients
	return func() { ReplaceGlobals(prev) }
}
