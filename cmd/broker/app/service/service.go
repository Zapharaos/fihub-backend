package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/protogen"
)

// Service is the implementation of the BrokerService interface.
type Service struct {
	protogen.UnimplementedBrokerServiceServer
}

// Todo implements the Todo RPC method.
func (h *Service) Todo(ctx context.Context, req *protogen.TodoRequest) (*protogen.TodoResponse, error) {
	return &protogen.TodoResponse{}, nil
}
