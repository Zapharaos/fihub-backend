package clients

import (
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/brokerpb"
	"github.com/Zapharaos/fihub-backend/gen/go/healthpb"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/gen/go/transactionpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
)

type Clients struct {
	health      healthpb.HealthServiceClient
	user        userpb.UserServiceClient
	auth        authpb.AuthServiceClient
	security    securitypb.SecurityServiceClient
	broker      brokerpb.BrokerServiceClient
	transaction transactionpb.TransactionServiceClient
}

type ClientOption func(*Clients)

func WithHealthClient(health healthpb.HealthServiceClient) ClientOption {
	return func(c *Clients) { c.health = health }
}

func WithUserClient(user userpb.UserServiceClient) ClientOption {
	return func(c *Clients) { c.user = user }
}

func WithAuthClient(auth authpb.AuthServiceClient) ClientOption {
	return func(c *Clients) { c.auth = auth }
}

func WithSecurityClient(security securitypb.SecurityServiceClient) ClientOption {
	return func(c *Clients) { c.security = security }
}

func WithBrokerClient(broker brokerpb.BrokerServiceClient) ClientOption {
	return func(c *Clients) { c.broker = broker }
}

func WithTransactionClient(transaction transactionpb.TransactionServiceClient) ClientOption {
	return func(c *Clients) { c.transaction = transaction }
}

func NewClients(opts ...ClientOption) Clients {
	var c Clients
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

func (c Clients) Health() healthpb.HealthServiceClient {
	return c.health
}

func (c Clients) User() userpb.UserServiceClient {
	return c.user
}

func (c Clients) Auth() authpb.AuthServiceClient {
	return c.auth
}

func (c Clients) Security() securitypb.SecurityServiceClient {
	return c.security
}

func (c Clients) Broker() brokerpb.BrokerServiceClient {
	return c.broker
}

func (c Clients) Transaction() transactionpb.TransactionServiceClient {
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
