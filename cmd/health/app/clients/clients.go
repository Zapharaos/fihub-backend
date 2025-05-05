package clients

import (
	"fmt"
)

type Clients struct {
	services map[string]interface{}
}

func NewClients() Clients {
	return Clients{services: make(map[string]interface{})}
}

func (c Clients) Register(serviceName string, client interface{}) {
	c.services[serviceName] = client
}

func (c Clients) Get(serviceName string) (interface{}, bool) {
	client, ok := c.services[serviceName]
	return client, ok
}

func GetTypedClient[T any](c Clients, serviceName string) (T, error) {
	raw, ok := c.services[serviceName]
	if !ok {
		var zero T
		return zero, fmt.Errorf("service %s not registered", serviceName)
	}
	client, ok := raw.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("service %s has wrong type", serviceName)
	}
	return client, nil
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
