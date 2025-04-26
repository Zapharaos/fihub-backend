package service

import (
	"github.com/Zapharaos/fihub-backend/protogen"
)

// Service is the implementation of the SecurityService interface.
type Service struct {
	protogen.UnimplementedSecurityServiceServer
}
