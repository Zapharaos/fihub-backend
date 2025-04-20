package clients

import (
	"github.com/Zapharaos/fihub-backend/protogen"
)

type Clients struct {
	health      protogen.HealthServiceClient
	transaction protogen.TransactionServiceClient
}

func NewClients(health protogen.HealthServiceClient, transaction protogen.TransactionServiceClient) Clients {
	return Clients{
		health:      health,
		transaction: transaction,
	}
}

func (c Clients) Health() protogen.HealthServiceClient {
	return c.health
}

func (c Clients) Transaction() protogen.TransactionServiceClient {
	return c.transaction
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
