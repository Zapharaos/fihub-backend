package users

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/google/uuid"
	"time"
)

// UserWithPassword extends User with a password field for authentication purposes
type UserWithPassword struct {
	User
	Password string `json:"password"`
}

// User represents a user entity in the system
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FullName returns the full name of the user.
func (u *User) FullName() string {
	return u.LastName + " " + u.FirstName
}

// IsValid checks if a user is valid and has no missing mandatory PGFields
// * Login must not be empty
// * Login must not be shorter than 4 characters
// * LastName must not be empty
func (u *User) IsValid() (bool, error) {
	if u.Email == "" {
		return false, errors.New("missing Email")
	}
	if !email.IsValid(u.Email) {
		return false, errors.New("email is not valid")
	}
	if u.FirstName == "" {
		return false, errors.New("missing Firstname")
	}
	if u.LastName == "" {
		return false, errors.New("missing Lastname")
	}
	return true, nil
}

// IsValid checks if a user with password is valid and has no missing mandatory PGFields
// * User must be valid (see User struct)
// * Password must not be empty
// * Password must not be shorter than 8 characters
func (u *UserWithPassword) IsValid() (bool, error) {
	if ok, err := u.User.IsValid(); !ok {
		return false, err
	}
	if u.Password == "" {
		return false, errors.New("missing Password")
	}
	if len(u.Password) < 8 {
		return false, errors.New("password is too short (less than 6 characters)")
	}
	if len(u.Password) > 64 {
		return false, errors.New("password is too long (more than 64 characters)")
	}
	return true, nil
}
