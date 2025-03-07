package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/handlers"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreateUser tests the CreateUser handler
func TestCreateUser(t *testing.T) {
	// Prepare data
	invalidUser := users.UserInputCreate{}
	invalidUserBody, _ := json.Marshal(invalidUser)
	validUser := users.UserInputCreate{
		UserInputPassword: users.UserInputPassword{
			UserWithPassword: users.UserWithPassword{
				User: users.User{
					Email: "email@test.ut",
				},
				Password: "password",
			},
			Confirmation: "password",
		},
		Checkbox: true,
	}
	validUserBody, _ := json.Marshal(validUser)

	// Test cases
	tests := []struct {
		name           string
		user           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to decode",
			user: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fails at bad user input",
			user: invalidUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail to check existence",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Return(false, errors.New("error"))
				u.EXPECT().Create(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "User already exists",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Return(true, nil)
				u.EXPECT().Create(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail at create",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Return(false, nil)
				u.EXPECT().Create(gomock.Any()).Return(uuid.Nil, errors.New("error"))
				u.EXPECT().Get(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Fails to retrieve the user",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Return(false, nil)
				u.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				u.EXPECT().Get(gomock.Any()).Return(users.User{}, false, errors.New("error"))
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Could not find the user",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Return(false, nil)
				u.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				u.EXPECT().Get(gomock.Any()).Return(users.User{}, false, nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Exists(gomock.Any()).Return(false, nil)
				u.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				u.EXPECT().Get(gomock.Any()).Return(users.User{}, true, nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(tt.user))

			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.CreateUser(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetUserSelf tests the GetUserSelf handler
func TestGetUserSelf(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to retrieve from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, false)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/self", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetUserSelf(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestUpdateUserSelf tests the UpdateUserSelf handler
func TestUpdateUserSelf(t *testing.T) {
	// Prepare data
	invalidUser := users.User{}
	invalidUserBody, _ := json.Marshal(invalidUser)
	validUser := users.User{
		Email: "email@test.ut",
	}
	validUserBody, _ := json.Marshal(validUser)

	// Test cases
	tests := []struct {
		name           string
		user           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to retrieve from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, false)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Update(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fails to decode",
			user: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Update(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fails at bad user input",
			user: invalidUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Update(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail at update",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Update(gomock.Any()).Return(errors.New("error"))
				u.EXPECT().Get(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Fails to retrieve the user",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Update(gomock.Any()).Return(nil)
				u.EXPECT().Get(gomock.Any()).Return(users.User{}, false, errors.New("error"))
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Could not find the user",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Update(gomock.Any()).Return(nil)
				u.EXPECT().Get(gomock.Any()).Return(users.User{}, false, nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Update(gomock.Any()).Return(nil)
				u.EXPECT().Get(gomock.Any()).Return(users.User{}, true, nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/users/self", bytes.NewBuffer(tt.user))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.UpdateUserSelf(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestChangeUserPassword tests the ChangeUserPassword handler
func TestChangeUserPassword(t *testing.T) {
	// Prepare data
	invalidPassword := users.UserInputPassword{}
	invalidPasswordBody, _ := json.Marshal(invalidPassword)
	validPassword := users.UserInputPassword{
		UserWithPassword: users.UserWithPassword{
			User: users.User{
				Email: "email@test.ut",
			},
			Password: "password",
		},
		Confirmation: "password",
	}
	validPasswordBody, _ := json.Marshal(validPassword)

	// Test cases
	tests := []struct {
		name           string
		password       []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to retrieve from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, false)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithPassword(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:     "Fails to decode",
			password: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithPassword(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Fails at bad user input",
			password: invalidPasswordBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithPassword(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Fail at update password",
			password: validPasswordBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithPassword(gomock.Any()).Return(errors.New("error"))
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:     "Succeeded",
			password: validPasswordBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithPassword(gomock.Any()).Return(nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/users/password", bytes.NewBuffer(tt.password))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.ChangeUserPassword(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestDeleteUserSelf tests the DeleteUserSelf handler
func TestDeleteUserSelf(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to retrieve from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, false)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Delete(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail at delete",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Delete(gomock.Any()).Return(nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/users/self", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.DeleteUserSelf(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetUser tests the GetUser handler
func TestGetUser(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "Fails to check permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK, // should be http.StatusUnauthorized, but not with mock
		},
		{
			name: "Fails to retrieve the user",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(users.User{}, false, errors.New("error"))
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Could not find the user",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(users.User{}, false, nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Fails to retrieve user roles",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(2)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(users.User{}, true, nil)
				users.ReplaceGlobals(u)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return([]roles.Role{}, errors.New("error"))
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(2)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().Get(gomock.Any()).Return(users.User{}, true, nil)
				users.ReplaceGlobals(u)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return([]roles.Role{{Id: uuid.Nil}}, nil)
				roles.ReplaceGlobals(r)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Return([]permissions.Permission{}, nil)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/{id}", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetUser(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestSetUser tests the SetUser handler
func TestSetUser(t *testing.T) {
	// Prepare data
	validUser := users.UserWithRoles{
		User: users.User{
			Email: "email@test.ut",
		},
	}
	validUserBody, _ := json.Marshal(validUser)

	// Test cases
	tests := []struct {
		name           string
		user           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to parse param",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "Fails to check permission",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithRoles(gomock.Any(), gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK, // should be http.StatusUnauthorized, but not with mock
		},
		{
			name: "Fails to decode",
			user: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithRoles(gomock.Any(), gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail at update",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithRoles(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithRoles(gomock.Any(), gomock.Any()).Return(nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/users/{id}", bytes.NewBuffer(tt.user))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.SetUser(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestSetUserRoles tests the SetUserRoles handler
func TestSetUserRoles(t *testing.T) {
	// Prepare data
	validUser := users.UserWithRoles{
		User: users.User{
			Email: "email@test.ut",
		},
		Roles: []roles.RoleWithPermissions{
			{
				Role: roles.Role{Id: uuid.Nil},
			},
		},
	}
	validUserBody, _ := json.Marshal(validUser)

	// Test cases
	tests := []struct {
		name           string
		user           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to parse param",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "Fails to check permission",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().SetUserRoles(gomock.Any(), gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK, // should be http.StatusUnauthorized, but not with mock
		},
		{
			name: "Fails to decode",
			user: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().SetUserRoles(gomock.Any(), gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail at set roles",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().SetUserRoles(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().SetUserRoles(gomock.Any(), gomock.Any()).Return(nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/users/{id}/roles", bytes.NewBuffer(tt.user))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.SetUserRoles(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetUserRoles tests the GetUserRoles handler
func TestGetUserRoles(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "Fails to check permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Times(0)
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusOK, // should be http.StatusUnauthorized, but not with mock
		},
		{
			name: "Fail at get roles",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return([]roles.Role{}, errors.New("error"))
				roles.ReplaceGlobals(r)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				r := mocks.NewRolesRepository(ctrl)
				r.EXPECT().GetRolesByUserId(gomock.Any()).Return([]roles.Role{{Id: uuid.Nil}}, nil)
				roles.ReplaceGlobals(r)
				p := mocks.NewPermissionsRepository(ctrl)
				p.EXPECT().GetAllByRoleId(gomock.Any()).Return([]permissions.Permission{}, nil)
				permissions.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/{id}/roles", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetUserRoles(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetAllUsersWithRoles tests the GetAllUsersWithRoles handler
func TestGetAllUsersWithRoles(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to check permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), "admin.users.list").Return(false)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetAllWithRoles().Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK, // should be http.StatusUnauthorized, but not with mock
		},
		{
			name: "Fails to retrieve the users with roles",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), "admin.users.list").Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetAllWithRoles().Return([]users.UserWithRoles{}, errors.New("error"))
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), "admin.users.list").Return(true)
				handlers.ReplaceGlobals(m)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetAllWithRoles().Return([]users.UserWithRoles{}, nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetAllUsersWithRoles(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
