package users

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/google/uuid"
	"time"
)

// UserInput extends UserWithPassword with a password-confirmation and checkbox
type UserInput struct {
	UserWithPassword
	Confirmation string `json:"confirmation"`
	Checkbox     bool   `json:"checkbox"`
}

// UserWithPassword extends User with a password field for authentication purposes
type UserWithPassword struct {
	User
	Password string `json:"password"`
}

// User represents a user entity in the system
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUser Returns a User struct without the password hash
func (u *UserWithPassword) ToUser() *User {
	return &User{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToUser Returns a User struct
func (u *UserInput) ToUser() User {
	return User{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToUserWithPassword Returns a UserWithPassword struct
func (u *UserInput) ToUserWithPassword() *UserWithPassword {
	return &UserWithPassword{
		User:     u.ToUser(),
		Password: u.Password,
	}
}

// IsValid checks if a user is valid and has no missing mandatory PGFields
// * Email must not be empty
// * Email must not be valid
func (u *User) IsValid() (bool, error) {
	if u.Email == "" {
		return false, errors.New("email-required")
	}
	if !email.IsValid(u.Email) {
		return false, errors.New("email-invalid")
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
		return false, errors.New("password-required")
	}
	if len(u.Password) < 8 {
		return false, errors.New("password-invalid")
	}
	if len(u.Password) > 64 {
		return false, errors.New("password-invalid")
	}
	return true, nil
}

// IsValid checks if a user input is valid and has no missing mandatory PGFields
// * UserWithPassword must be valid (see UserWithPassword struct)
// * Confirmation must not be empty
// * Confirmation must be equal to UserWithPassword.Password
// * Checkbox must be true
func (u *UserInput) IsValid() (bool, error) {
	if ok, err := u.UserWithPassword.IsValid(); !ok {
		return false, err
	}
	if u.Confirmation == "" {
		return false, errors.New("confirmation-required")
	}
	if u.Confirmation != u.Password {
		return false, errors.New("confirmation-invalid")
	}
	if u.Checkbox != true {
		return false, errors.New("checkbox-invalid")
	}
	return true, nil
}
