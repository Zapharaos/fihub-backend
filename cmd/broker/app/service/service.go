package service

import (
	"github.com/Zapharaos/fihub-backend/protogen"
)

// Service is the implementation of the BrokerService interface.
type Service struct {
	protogen.UnimplementedBrokerServiceServer
}
