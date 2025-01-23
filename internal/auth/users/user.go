package users

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/google/uuid"
	"time"
)

// UserInputCreate extends UserInputPassword with a password-confirmation and checkbox
type UserInputCreate struct {
	UserInputPassword
	Checkbox bool `json:"checkbox"`
}

// UserInputPassword extends UserWithPassword with a password-confirmation and checkbox
type UserInputPassword struct {
	UserWithPassword
	Confirmation string `json:"confirmation"`
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
func (u UserWithPassword) ToUser() User {
	return User{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToUser Returns a User struct
func (u UserInputCreate) ToUser() User {
	return User{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToUserWithPassword Returns a UserWithPassword struct
func (u UserInputPassword) ToUserWithPassword() UserWithPassword {
	return UserWithPassword{
		User:     u.ToUser(),
		Password: u.Password,
	}
}

// IsValid checks if a user is valid and has no missing mandatory PGFields
// * Email must not be empty
// * Email must not be valid
func (u User) IsValid() (bool, error) {
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
// * Password must not be valid (see isValidPassword function)
func (u UserWithPassword) IsValid() (bool, error) {
	if ok, err := u.User.IsValid(); !ok {
		return false, err
	}
	return isValidPassword(u.Password)
}

// IsValid checks if a user input is valid and has no missing mandatory PGFields
// * UserWithPassword must be valid (see UserWithPassword struct)
// * Confirmation must not be valid (see isValidConfirmation function)
func (u UserInputPassword) IsValid() (bool, error) {
	if ok, err := u.UserWithPassword.IsValid(); !ok {
		return false, err
	}
	return isValidConfirmation(u.Password, u.Confirmation)
}

// IsValid checks if a user input is valid and has no missing mandatory PGFields
// * UserInputPassword must be valid (see UserInputPassword struct)
// * Checkbox must be true
func (u UserInputCreate) IsValid() (bool, error) {
	if ok, err := u.UserInputPassword.IsValid(); !ok {
		return false, err
	}
	if !u.Checkbox {
		return false, errors.New("checkbox-invalid")
	}
	return true, nil
}

// isValidPassword checks if a password is valid
// * Password must not be empty
// * Password must not be shorter than 8 characters
// * Password must not be longer than 64 characters
func isValidPassword(password string) (bool, error) {
	if password == "" {
		return false, errors.New("password-required")
	}
	if len(password) < 8 {
		return false, errors.New("password-invalid")
	}
	if len(password) > 64 {
		return false, errors.New("password-invalid")
	}
	return true, nil
}

// isValidConfirmation checks if a password and a confirmation are valid
// * Confirmation must not be empty
// * Confirmation must be equal to password
func isValidConfirmation(password, confirmation string) (bool, error) {
	if confirmation == "" {
		return false, errors.New("confirmation-required")
	}
	if confirmation != password {
		return false, errors.New("confirmation-invalid")
	}
	return true, nil
}

// IsValidPassword checks if a user input password is valid
// * Password must not be valid (see isValidPassword function)
// * Confirmation must not be valid (see isValidConfirmation function)
func (u UserInputPassword) IsValidPassword() (bool, error) {
	if ok, err := isValidPassword(u.Password); !ok {
		return false, err
	}
	if ok, err := isValidConfirmation(u.Password, u.Confirmation); !ok {
		return false, err
	}
	return true, nil
}
