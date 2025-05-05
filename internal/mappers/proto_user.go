package mappers

import (
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserToProto converts a models.User to a userpb.User
func UserToProto(user models.User) *userpb.User {
	return &userpb.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// UserFromProto converts a userpb.User to a models.User
func UserFromProto(user *userpb.User) models.User {
	return models.User{
		ID:        uuid.MustParse(user.GetId()),
		Email:     user.GetEmail(),
		CreatedAt: user.GetCreatedAt().AsTime(),
		UpdatedAt: user.GetUpdatedAt().AsTime(),
	}
}

// UsersToProto converts a slice of models.User to a slice of userpb.User
func UsersToProto(users models.Users) []*userpb.User {
	protoUsers := make([]*userpb.User, len(users))
	for i, user := range users {
		protoUsers[i] = UserToProto(user)
	}
	return protoUsers
}

// UsersFromProto converts a slice of userpb.User to a slice of models.User
func UsersFromProto(users []*userpb.User) models.Users {
	protoUsers := make(models.Users, len(users))
	for i, user := range users {
		protoUsers[i] = UserFromProto(user)
	}
	return protoUsers
}
