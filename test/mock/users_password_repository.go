package mock

import (
	"github.com/Zapharaos/fihub-backend/internal/auth/password"
	"github.com/google/uuid"
	"time"
)

// PasswordRepository represents a mock password.Repository
type PasswordRepository struct {
	ID      uuid.UUID
	Found   bool
	Error   error
	Request password.Request
	time    time.Time
}

// NewPasswordRepository creates a new PasswordRepository of the password.Repository interface
func NewPasswordRepository() password.Repository {
	r := PasswordRepository{}
	var repo password.Repository = &r
	return repo
}

func (m PasswordRepository) Create(_ password.Request) (password.Request, error) {
	return m.Request, m.Error
}

func (m PasswordRepository) GetRequestID(_ uuid.UUID, _ string) (uuid.UUID, error) {
	return m.ID, m.Error
}

func (m PasswordRepository) GetExpiresAt(_ uuid.UUID) (time.Time, error) {
	return m.time, m.Error
}

func (m PasswordRepository) Delete(_ uuid.UUID) error {
	return m.Error
}

func (m PasswordRepository) Valid(_ uuid.UUID, _ uuid.UUID) (bool, error) {
	return m.Found, m.Error
}

func (m PasswordRepository) ValidForUser(_ uuid.UUID) (bool, error) {
	return m.Found, m.Error
}
