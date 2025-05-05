package models

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TestUser_IsValid tests the IsValid method of the User struct
func TestUser_IsValid(t *testing.T) {
	// Define valid values
	validEmail := "test@example.com"

	// Define test cases
	tests := []struct {
		name     string // Test case name
		user     User   // User instance to test
		expected bool   // Expected result
		err      error  // Expected error
	}{
		{
			name: "valid User",
			user: User{
				Email: validEmail,
			},
			expected: true,
			err:      nil,
		},
		{
			name:     "missing email",
			user:     User{},
			expected: false,
			err:      ErrEmailRequired,
		},
		{
			name: "invalid email",
			user: User{
				Email: "invalid-email",
			},
			expected: false,
			err:      ErrEmailInvalid,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.user.IsValid()
			assert.Equal(t, tt.expected, got)
			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUserWithPassword_IsValid tests the IsValid method of the UserWithPassword struct
func TestUserWithPassword_IsValid(t *testing.T) {
	// Define valid values
	validPassword := "validpassword"
	validEmail := "test@example.com"
	validUser := User{
		Email: validEmail,
	}

	// Define test cases
	tests := []struct {
		name     string           // Test case name
		user     UserWithPassword // UserWithPassword instance to test
		expected bool             // Expected result
		err      error            // Expected error
	}{
		{
			name: "valid User with password",
			user: UserWithPassword{
				User:     validUser,
				Password: validPassword,
			},
			expected: true,
			err:      nil,
		},
		{
			name: "invalid User",
			user: UserWithPassword{
				User: User{},
			},
			expected: false,
			err:      ErrEmailRequired,
		},
		{
			name: "invalid password",
			user: UserWithPassword{
				User:     validUser,
				Password: strings.Repeat("a", PasswordMinLength-1),
			},
			expected: false,
			err:      ErrPasswordInvalid,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.user.IsValid()
			assert.Equal(t, tt.expected, got)
			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUserInputPassword_IsValid tests the IsValid method of the UserInputPassword struct
func TestUserInputPassword_IsValid(t *testing.T) {
	// Define valid values
	validPassword := "validpassword"
	validEmail := "test@example.com"
	validUser := User{
		Email: validEmail,
	}
	validUserWithPassword := UserWithPassword{
		User:     validUser,
		Password: validPassword,
	}

	// Define test cases
	tests := []struct {
		name     string            // Test case name
		user     UserInputPassword // UserInputPassword instance to test
		expected bool              // Expected result
		err      error             // Expected error
	}{
		{
			name: "valid User input password",
			user: UserInputPassword{
				UserWithPassword: validUserWithPassword,
				Confirmation:     validPassword,
			},
			expected: true,
			err:      nil,
		},
		{
			name: "invalid User with password",
			user: UserInputPassword{
				UserWithPassword: UserWithPassword{},
			},
			expected: false,
			err:      ErrEmailRequired,
		},
		{
			name: "password and confirmation mismatch",
			user: UserInputPassword{
				UserWithPassword: validUserWithPassword,
				Confirmation:     "differentpassword",
			},
			expected: false,
			err:      ErrConfirmationInvalid,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.user.IsValid()
			assert.Equal(t, tt.expected, got)
			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUserInputCreate_IsValid tests the IsValid method of the UserInputCreate struct
func TestUserInputCreate_IsValid(t *testing.T) {
	// Define valid values
	validPassword := "validpassword"
	validEmail := "test@example.com"
	validUser := User{
		Email: validEmail,
	}
	validUserWithPassword := UserWithPassword{
		User:     validUser,
		Password: validPassword,
	}
	validUserInputPassword := UserInputPassword{
		UserWithPassword: validUserWithPassword,
		Confirmation:     validPassword,
	}

	// Define test cases
	tests := []struct {
		name     string          // Test case name
		user     UserInputCreate // UserInputCreate instance to test
		expected bool            // Expected result
		err      error           // Expected error
	}{
		{
			name: "valid User input create",
			user: UserInputCreate{
				UserInputPassword: validUserInputPassword,
				Checkbox:          true,
			},
			expected: true,
			err:      nil,
		},
		{
			name: "checkbox not checked",
			user: UserInputCreate{
				UserInputPassword: validUserInputPassword,
				Checkbox:          false,
			},
			expected: false,
			err:      ErrCheckboxInvalid,
		},
		{
			name: "invalid User input password",
			user: UserInputCreate{
				UserInputPassword: UserInputPassword{},
			},
			expected: false,
			err:      ErrEmailRequired,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.user.IsValid()
			assert.Equal(t, tt.expected, got)
			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUserInputPassword_IsValidPassword tests the IsValidPassword method of the UserInputPassword struct
func TestUserInputPassword_IsValidPassword(t *testing.T) {
	// Define valid values
	validPassword := "validpassword"
	validEmail := "test@example.com"
	validUser := User{
		Email: validEmail,
	}
	validUserWithPassword := UserWithPassword{
		User:     validUser,
		Password: validPassword,
	}

	// Define test cases
	tests := []struct {
		name     string            // Test case name
		user     UserInputPassword // UserInputPassword instance to test
		expected bool              // Expected result
		err      error             // Expected error
	}{
		{
			name: "valid User input password",
			user: UserInputPassword{
				UserWithPassword: validUserWithPassword,
				Confirmation:     validPassword,
			},
			expected: true,
			err:      nil,
		},
		{
			name: "missing password",
			user: UserInputPassword{
				UserWithPassword: UserWithPassword{
					User:     validUser,
					Password: "",
				},
			},
			expected: false,
			err:      ErrPasswordInvalid,
		},
		{
			name: "short password",
			user: UserInputPassword{
				UserWithPassword: UserWithPassword{
					User:     validUser,
					Password: strings.Repeat("a", PasswordMinLength-1),
				},
			},
			expected: false,
			err:      ErrPasswordInvalid,
		},
		{
			name: "long password",
			user: UserInputPassword{
				UserWithPassword: UserWithPassword{
					User:     validUser,
					Password: strings.Repeat("a", PasswordMaxLength+1),
				},
			},
			expected: false,
			err:      ErrPasswordInvalid,
		},
		{
			name: "missing confirmation",
			user: UserInputPassword{
				UserWithPassword: validUserWithPassword,
				Confirmation:     "",
			},
			expected: false,
			err:      ErrConfirmationRequired,
		},
		{
			name: "password and confirmation mismatch",
			user: UserInputPassword{
				UserWithPassword: validUserWithPassword,
				Confirmation:     "differentpassword",
			},
			expected: false,
			err:      ErrConfirmationInvalid,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.user.IsValidPassword()
			assert.Equal(t, tt.expected, got)
			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
