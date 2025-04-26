package clients

import (
	"github.com/Zapharaos/fihub-backend/protogen"
)

type Clients struct {
	health      protogen.HealthServiceClient
	user        protogen.UserServiceClient
	security    protogen.SecurityServiceClient
	broker      protogen.BrokerServiceClient
	transaction protogen.TransactionServiceClient
}

type ClientOption func(*Clients)

func WithHealthClient(health protogen.HealthServiceClient) ClientOption {
	return func(c *Clients) { c.health = health }
}

func WithUserClient(user protogen.UserServiceClient) ClientOption {
	return func(c *Clients) { c.user = user }
}

func WithSecurityClient(security protogen.SecurityServiceClient) ClientOption {
	return func(c *Clients) { c.security = security }
}

func WithBrokerClient(broker protogen.BrokerServiceClient) ClientOption {
	return func(c *Clients) { c.broker = broker }
}

func WithTransactionClient(transaction protogen.TransactionServiceClient) ClientOption {
	return func(c *Clients) { c.transaction = transaction }
}

func NewClients(opts ...ClientOption) Clients {
	var c Clients
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

func (c Clients) Health() protogen.HealthServiceClient {
	return c.health
}

func (c Clients) User() protogen.UserServiceClient {
	return c.user
}

func (c Clients) Security() protogen.SecurityServiceClient {
	return c.security
}

func (c Clients) Broker() protogen.BrokerServiceClient {
	return c.broker
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
