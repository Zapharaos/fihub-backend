package mock

import (
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/google/uuid"
)

// UsersRepository represents a mock users.Repository
type UsersRepository struct {
	ID             uuid.UUID
	Err            error
	Found          bool
	User           users.User
	Users          []users.User
	UserWithRoles  users.UserWithRoles
	UsersWithRoles []users.UserWithRoles
}

// NewUsersRepository creates a new UsersRepository of the users.Repository interface
func NewUsersRepository(r UsersRepository) users.Repository {
	var repo users.Repository
	repo = &r
	return repo
}

func (m UsersRepository) Create(_ users.UserWithPassword) (uuid.UUID, error) {
	return m.ID, m.Err
}

func (m UsersRepository) Get(_ uuid.UUID) (users.User, bool, error) {
	return m.User, m.Found, m.Err
}

func (m UsersRepository) GetByEmail(_ string) (users.User, bool, error) {
	return m.User, m.Found, m.Err
}

func (m UsersRepository) Exists(_ string) (bool, error) {
	return m.Found, m.Err
}

func (m UsersRepository) Authenticate(_ string, _ string) (users.User, bool, error) {
	return m.User, m.Found, m.Err
}

func (m UsersRepository) Update(_ users.User) error {
	return m.Err
}

func (m UsersRepository) UpdateWithPassword(_ users.UserWithPassword) error {
	return m.Err
}

func (m UsersRepository) Delete(_ uuid.UUID) error {
	return m.Err
}

func (m UsersRepository) GetWithRoles(_ uuid.UUID) (users.UserWithRoles, error) {
	return m.UserWithRoles, m.Err
}

func (m UsersRepository) GetAllWithRoles() ([]users.UserWithRoles, error) {
	return m.UsersWithRoles, m.Err
}

func (m UsersRepository) GetUsersByRoleID(_ uuid.UUID) ([]users.User, error) {
	return m.Users, m.Err
}

func (m UsersRepository) UpdateWithRoles(_ users.UserWithRoles, _ []uuid.UUID) error {
	return m.Err
}

func (m UsersRepository) SetUserRoles(_ uuid.UUID, _ []uuid.UUID) error {
	return m.Err
}

func (m UsersRepository) AddUsersRole(_ []uuid.UUID, _ uuid.UUID) error {
	return m.Err
}

func (m UsersRepository) RemoveUsersRole(_ []uuid.UUID, _ uuid.UUID) error {
	return m.Err
}
