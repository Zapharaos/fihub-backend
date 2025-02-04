package users

import "github.com/google/uuid"

// MockRepository represents a mock Repository
type MockRepository struct {
	ID             uuid.UUID
	Err            error
	Found          bool
	User           User
	Users          []User
	UserWithRoles  UserWithRoles
	UsersWithRoles []UserWithRoles
}

// NewMockRepository creates a new MockRepository of the Repository interface
func NewMockRepository() Repository {
	r := MockRepository{}
	var repo Repository
	repo = &r
	return repo
}

func (m MockRepository) Create(_ UserWithPassword) (uuid.UUID, error) {
	return m.ID, m.Err
}

func (m MockRepository) Get(_ uuid.UUID) (User, bool, error) {
	return m.User, m.Found, m.Err
}

func (m MockRepository) GetByEmail(_ string) (User, bool, error) {
	return m.User, m.Found, m.Err
}

func (m MockRepository) Exists(_ string) (bool, error) {
	return m.Found, m.Err
}

func (m MockRepository) Authenticate(_ string, _ string) (User, bool, error) {
	return m.User, m.Found, m.Err
}

func (m MockRepository) Update(_ User) error {
	return m.Err
}

func (m MockRepository) UpdateWithPassword(_ UserWithPassword) error {
	return m.Err
}

func (m MockRepository) Delete(_ uuid.UUID) error {
	return m.Err
}

func (m MockRepository) GetWithRoles(_ uuid.UUID) (UserWithRoles, error) {
	return m.UserWithRoles, m.Err
}

func (m MockRepository) GetAllWithRoles() ([]UserWithRoles, error) {
	return m.UsersWithRoles, m.Err
}

func (m MockRepository) GetUsersByRoleID(_ uuid.UUID) ([]User, error) {
	return m.Users, m.Err
}

func (m MockRepository) UpdateWithRoles(_ UserWithRoles, _ []uuid.UUID) error {
	return m.Err
}

func (m MockRepository) SetUserRoles(_ uuid.UUID, _ []uuid.UUID) error {
	return m.Err
}

func (m MockRepository) AddUsersRole(_ []uuid.UUID, _ uuid.UUID) error {
	return m.Err
}

func (m MockRepository) RemoveUsersRole(_ []uuid.UUID, _ uuid.UUID) error {
	return m.Err
}
