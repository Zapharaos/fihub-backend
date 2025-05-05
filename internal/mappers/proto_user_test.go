package mappers

import (
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

// Test_UserToProto tests the UserToProto function
func Test_UserToProto(t *testing.T) {
	userId := uuid.New()
	testDate := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	user := models.User{
		ID:        userId,
		Email:     "email@example.com",
		CreatedAt: testDate,
	}

	result := UserToProto(user)

	assert.Equal(t, userId.String(), result.Id)
	assert.Equal(t, "email@example.com", result.Email)
	assert.Equal(t, testDate.Unix(), result.CreatedAt.AsTime().Unix())
}

// Test_UserFromProto tests the UserFromProto function
func Test_UserFromProto(t *testing.T) {
	userId := uuid.New()
	testDate := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	protoUser := &userpb.User{
		Id:        userId.String(),
		Email:     "email@example.com",
		CreatedAt: timestamppb.New(testDate),
	}

	result := UserFromProto(protoUser)

	assert.Equal(t, userId, result.ID)
	assert.Equal(t, "email@example.com", result.Email)
	assert.Equal(t, testDate.Unix(), result.CreatedAt.Unix())
}

// Test_UsersToProto tests the UsersToProto function
func Test_UsersToProto(t *testing.T) {
	userId := uuid.New()
	testDate := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	users := []models.User{
		{
			ID:        userId,
			Email:     "email@example.com",
			CreatedAt: testDate,
		},
	}

	result := UsersToProto(users)

	assert.Equal(t, userId.String(), result[0].Id)
	assert.Equal(t, "email@example.com", result[0].Email)
	assert.Equal(t, testDate.Unix(), result[0].CreatedAt.AsTime().Unix())
}

// Test_UsersFromProto tests the UsersFromProto function
func Test_UsersFromProto(t *testing.T) {
	userId := uuid.New()
	testDate := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)

	protoUsers := []*userpb.User{
		{
			Id:        userId.String(),
			Email:     "email@example.com",
			CreatedAt: timestamppb.New(testDate),
		},
	}

	result := UsersFromProto(protoUsers)

	assert.Equal(t, userId, result[0].ID)
	assert.Equal(t, "email@example.com", result[0].Email)
	assert.Equal(t, testDate.Unix(), result[0].CreatedAt.Unix())
}
