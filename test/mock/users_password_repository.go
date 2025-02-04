package mock

import (
	"github.com/Zapharaos/fihub-backend/internal/auth/password"
	"github.com/google/uuid"
	"time"
)

// UsersPasswordRepository represents a mock password.Repository
type UsersPasswordRepository struct {
	ID      uuid.UUID
	Found   bool
	Error   error
	Request password.Request
	time    time.Time
}

// NewUsersPasswordRepository creates a new UsersPasswordRepository of the password.Repository interface
func NewUsersPasswordRepository() password.Repository {
	r := UsersPasswordRepository{}
	var repo password.Repository = &r
	return repo
}

func (m UsersPasswordRepository) Create(_ password.Request) (password.Request, error) {
	return m.Request, m.Error
}

func (m UsersPasswordRepository) GetRequestID(_ uuid.UUID, _ string) (uuid.UUID, error) {
	return m.ID, m.Error
}

func (m UsersPasswordRepository) GetExpiresAt(_ uuid.UUID) (time.Time, error) {
	return m.time, m.Error
}

func (m UsersPasswordRepository) Delete(_ uuid.UUID) error {
	return m.Error
}

func (m UsersPasswordRepository) Valid(_ uuid.UUID, _ uuid.UUID) (bool, error) {
	return m.Found, m.Error
}

func (m UsersPasswordRepository) ValidForUser(_ uuid.UUID) (bool, error) {
	return m.Found, m.Error
}
