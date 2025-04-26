package service

import (
	"github.com/Zapharaos/fihub-backend/protogen"
)

// Service is the implementation of the UserService interface.
type Service struct {
	protogen.UnimplementedUserServiceServer
}
