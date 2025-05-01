package service

import (
	"github.com/Zapharaos/fihub-backend/gen/go/brokerpb"
)

// Service is the implementation of the BrokerService interface.
type Service struct {
	brokerpb.UnimplementedBrokerServiceServer
}
