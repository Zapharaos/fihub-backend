package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/auth/password"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/handlers"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestCreatePasswordResetRequest tests the CreatePasswordResetRequest handler
func TestCreatePasswordResetRequest(t *testing.T) {
	// Define data
	validRequest := password.InputRequest{
		Email: "test@email.tu",
	}
	validRequestBody, _ := json.Marshal(validRequest)

	// Define tests
	tests := []struct {
		name           string
		body           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to decode",
			body: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetByEmail(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to retrieve user",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetByEmail(gomock.Any()).Return(users.User{}, false, errors.New("error"))
				users.ReplaceGlobals(u)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().ValidForUser(gomock.Any()).Times(0)
				password.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "could not find user",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetByEmail(gomock.Any()).Return(users.User{}, false, nil)
				users.ReplaceGlobals(u)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().ValidForUser(gomock.Any()).Times(0)
				password.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to search for user token",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetByEmail(gomock.Any()).Return(users.User{}, true, nil)
				users.ReplaceGlobals(u)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().ValidForUser(gomock.Any()).Return(true, errors.New("error"))
				p.EXPECT().GetExpiresAt(gomock.Any()).Times(0)
				p.EXPECT().Create(gomock.Any()).Times(0)
				password.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "user already has token, get expires_at fails",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetByEmail(gomock.Any()).Return(users.User{}, true, nil)
				users.ReplaceGlobals(u)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().ValidForUser(gomock.Any()).Return(true, nil)
				p.EXPECT().GetExpiresAt(gomock.Any()).Return(time.Now(), errors.New("error"))
				p.EXPECT().Create(gomock.Any()).Times(0)
				password.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "user already has token",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetByEmail(gomock.Any()).Return(users.User{}, true, nil)
				users.ReplaceGlobals(u)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().ValidForUser(gomock.Any()).Return(true, nil)
				p.EXPECT().GetExpiresAt(gomock.Any()).Return(time.Now(), nil)
				p.EXPECT().Create(gomock.Any()).Times(0)
				password.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "fails to create token request",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().GetByEmail(gomock.Any()).Return(users.User{}, true, nil)
				users.ReplaceGlobals(u)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().ValidForUser(gomock.Any()).Return(false, nil)
				p.EXPECT().Create(gomock.Any()).Return(password.Request{}, errors.New("error"))
				password.ReplaceGlobals(p)
				// TODO : not localizer
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	/*{
	name: "fails to localizer language",
		mockSetup: func(ctrl *gomock.Controller) {
		u := mocks.NewUsersRepository(ctrl)
		u.EXPECT().GetByEmail(gomock.Any()).Return(&users.User{}, nil)
		users.ReplaceGlobals(u)
		p := mocks.NewUsersPasswordRepository(ctrl)
		p.EXPECT().GetByUserID(gomock.Any()).Return(nil, nil)
		p.EXPECT().Create(gomock.Any()).Return(nil)
		password.ReplaceGlobals(p)
		s := mocks.NewTranslationService(ctrl)
		s.EXPECT().Localizer(gomock.Any()).Return(nil, errors.New("error"))
		translation.ReplaceGlobals(s)
	},
		expectedStatus: http.StatusInternalServerError,
	},
	{
	name: "fails to send email",
		mockSetup: func(ctrl *gomock.Controller) {
		u := mocks.NewUsersRepository(ctrl)
		u.EXPECT().GetByEmail(gomock.Any()).Return(&users.User{}, nil)
		users.ReplaceGlobals(u)
		p := mocks.NewUsersPasswordRepository(ctrl)
		p.EXPECT().GetByUserID(gomock.Any()).Return(nil, nil)
		p.EXPECT().Create(gomock.Any()).Return(nil)
		password.ReplaceGlobals(p)
		s := mocks.NewTranslationService(ctrl)
		s.EXPECT().Localizer(gomock.Any()).Return(&translation.Localizer{}, nil)
		s.EXPECT().Message(gomock.Any(), gomock.Any()).Return("message")
		translation.ReplaceGlobals(s)
		e := mocks.NewEmailService(ctrl)
		e.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
		email.ReplaceGlobals(e)
	},
		expectedStatus: http.StatusInternalServerError,
	},
	{
	name: "succeeded",
		mockSetup: func(ctrl *gomock.Controller) {
		u := mocks.NewUsersRepository(ctrl)
		u.EXPECT().GetByEmail(gomock.Any()).Return(&users.User{}, nil)
		users.ReplaceGlobals(u)
		p := mocks.NewUsersPasswordRepository(ctrl)
		p.EXPECT().GetByUserID(gomock.Any()).Return(nil, nil)
		p.EXPECT().Create(gomock.Any()).Return(nil)
		password.ReplaceGlobals(p)
		s := mocks.NewTranslationService(ctrl)
		s.EXPECT().Localizer(gomock.Any()).Return(&translation.Localizer{}, nil)
		s.EXPECT().Message(gomock.Any(), gomock.Any()).Return("message")
		translation.ReplaceGlobals(s)
		e := mocks.NewEmailService(ctrl)
		e.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		email.ReplaceGlobals(e)
	},
		expectedStatus: http.StatusOK,
	},*/

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/auth/password/reset", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.CreatePasswordResetRequest(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestCreatePasswordResetRequest tests the CreatePasswordResetRequest handler
func TestGetPasswordResetRequestID(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to parse param id",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				h.EXPECT().ParseParamString(gomock.Any(), gomock.Any(), "token").Times(0)
				handlers.ReplaceGlobals(h)
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but no with mock
		},
		{
			name: "fails to parse param token",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, true)
				h.EXPECT().ParseParamString(gomock.Any(), gomock.Any(), "token").Return("", false)
				handlers.ReplaceGlobals(h)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().GetRequestID(gomock.Any(), gomock.Any()).Times(0)
				password.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but no with mock
		},
		{
			name: "fail at retrieve request",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, true)
				h.EXPECT().ParseParamString(gomock.Any(), gomock.Any(), "token").Return("", true)
				handlers.ReplaceGlobals(h)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().GetRequestID(gomock.Any(), gomock.Any()).Return(uuid.Nil, errors.New("error"))
				password.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "request not found",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, true)
				h.EXPECT().ParseParamString(gomock.Any(), gomock.Any(), "token").Return("", true)
				handlers.ReplaceGlobals(h)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().GetRequestID(gomock.Any(), gomock.Any()).Return(uuid.Nil, nil)
				password.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, true)
				h.EXPECT().ParseParamString(gomock.Any(), gomock.Any(), "token").Return("", true)
				handlers.ReplaceGlobals(h)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().GetRequestID(gomock.Any(), gomock.Any()).Return(uuid.New(), nil)
				password.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/auth/password/{id}/{token}", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetPasswordResetRequestID(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestResetPassword tests the ResetPassword handler
func TestResetPassword(t *testing.T) {
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

	// Define tests
	tests := []struct {
		name           string
		body           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to parse params",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, uuid.Nil, false)
				handlers.ReplaceGlobals(h)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().Valid(gomock.Any(), gomock.Any()).Times(0)
				password.ReplaceGlobals(p)
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but no with mock
		},
		{
			name: "fails to validate request",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), "request_id").Return(uuid.Nil, uuid.Nil, true)
				handlers.ReplaceGlobals(h)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().Valid(gomock.Any(), gomock.Any()).Return(false, errors.New("error"))
				password.ReplaceGlobals(p)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithPassword(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "request does not exist",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), "request_id").Return(uuid.Nil, uuid.Nil, true)
				handlers.ReplaceGlobals(h)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().Valid(gomock.Any(), gomock.Any()).Return(false, nil)
				password.ReplaceGlobals(p)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithPassword(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to decode",
			body: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), "request_id").Return(uuid.Nil, uuid.Nil, true)
				handlers.ReplaceGlobals(h)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().Valid(gomock.Any(), gomock.Any()).Return(true, nil)
				password.ReplaceGlobals(p)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithPassword(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at bad user password input",
			body: invalidPasswordBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), "request_id").Return(uuid.Nil, uuid.Nil, true)
				handlers.ReplaceGlobals(h)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().Valid(gomock.Any(), gomock.Any()).Return(true, nil)
				password.ReplaceGlobals(p)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithPassword(gomock.Any()).Times(0)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fail at password update",
			body: validPasswordBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), "request_id").Return(uuid.Nil, uuid.Nil, true)
				handlers.ReplaceGlobals(h)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().Valid(gomock.Any(), gomock.Any()).Return(true, nil)
				p.EXPECT().Delete(gomock.Any()).Times(0)
				password.ReplaceGlobals(p)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithPassword(gomock.Any()).Return(errors.New("error"))
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails at request removal",
			body: validPasswordBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), "request_id").Return(uuid.Nil, uuid.Nil, true)
				handlers.ReplaceGlobals(h)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().Valid(gomock.Any(), gomock.Any()).Return(true, nil)
				p.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				password.ReplaceGlobals(p)
				u := mocks.NewUsersRepository(ctrl)
				u.EXPECT().UpdateWithPassword(gomock.Any()).Return(nil)
				users.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validPasswordBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), "request_id").Return(uuid.Nil, uuid.Nil, true)
				handlers.ReplaceGlobals(h)
				p := mocks.NewUsersPasswordRepository(ctrl)
				p.EXPECT().Valid(gomock.Any(), gomock.Any()).Return(true, nil)
				p.EXPECT().Delete(gomock.Any()).Return(nil)
				password.ReplaceGlobals(p)
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
			r := httptest.NewRequest("PUT", "/api/v1/auth/password/{id}/{request_id}", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.ResetPassword(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
