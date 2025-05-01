package service

import (
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
)

// Service is the implementation of the SecurityService interface.
type Service struct {
	securitypb.UnimplementedSecurityServiceServer
}
