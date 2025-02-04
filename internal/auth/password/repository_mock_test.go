package password

import (
	"github.com/google/uuid"
	"time"
)

// MockRepository represents a mock Repository
type MockRepository struct {
	ID      uuid.UUID
	Found   bool
	Error   error
	Request Request
	time    time.Time
}

// NewMockRepository creates a new MockRepository of the Repository interface
func NewMockRepository() Repository {
	r := MockRepository{}
	var repo Repository = &r
	return repo
}

func (m MockRepository) Create(_ Request) (Request, error) {
	return m.Request, m.Error
}

func (m MockRepository) GetRequestID(_ uuid.UUID, _ string) (uuid.UUID, error) {
	return m.ID, m.Error
}

func (m MockRepository) GetExpiresAt(_ uuid.UUID) (time.Time, error) {
	return m.time, m.Error
}

func (m MockRepository) Delete(_ uuid.UUID) error {
	return m.Error
}

func (m MockRepository) Valid(_ uuid.UUID, _ uuid.UUID) (bool, error) {
	return m.Found, m.Error
}

func (m MockRepository) ValidForUser(_ uuid.UUID) (bool, error) {
	return m.Found, m.Error
}
