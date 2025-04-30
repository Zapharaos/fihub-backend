package mappers

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserToProto converts a models.User to a protogen.User
func UserToProto(user models.User) *protogen.User {
	return &protogen.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// UserFromProto converts a protogen.User to a models.User
func UserFromProto(user *protogen.User) models.User {
	return models.User{
		ID:        uuid.MustParse(user.GetId()),
		Email:     user.GetEmail(),
		CreatedAt: user.GetCreatedAt().AsTime(),
		UpdatedAt: user.GetUpdatedAt().AsTime(),
	}
}

// UsersToProto converts a slice of models.User to a slice of protogen.User
func UsersToProto(users models.Users) []*protogen.User {
	protoUsers := make([]*protogen.User, len(users))
	for i, user := range users {
		protoUsers[i] = UserToProto(user)
	}
	return protoUsers
}

// UsersFromProto converts a slice of protogen.User to a slice of models.User
func UsersFromProto(users []*protogen.User) models.Users {
	protoUsers := make(models.Users, len(users))
	for i, user := range users {
		protoUsers[i] = UserFromProto(user)
	}
	return protoUsers
}
