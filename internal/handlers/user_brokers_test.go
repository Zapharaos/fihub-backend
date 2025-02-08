package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/handlers"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreateUserBroker tests the CreateUserBroker handler
func TestCreateUserBroker(t *testing.T) {
	// Prepare data
	invalidBroker := brokers.UserInput{}
	invalidBrokerBody, _ := json.Marshal(invalidBroker)
	validBroker := brokers.UserInput{
		BrokerID: uuid.New().String(),
	}
	validBrokerBody, _ := json.Marshal(validBroker)

	// Test cases
	tests := []struct {
		name           string
		broker         []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, false)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Fails to decode",
			broker: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Fails at bad user broker input",
			broker: invalidBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Fails to verify the broker existence",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, false, errors.New("error"))
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Fails to find the broker",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, false, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, bu, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Broker is not enabled",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{Disabled: true}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, bu, nil))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Fails to verify the user broker existence",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, errors.New("error"))
				bu.EXPECT().Create(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "User broker already exists",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				bu.EXPECT().Create(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, bu, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Fails to create user broker",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bu.EXPECT().Create(gomock.Any()).Return(errors.New("error"))
				bu.EXPECT().GetAll(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Fails to retrieve all user brokers",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bu.EXPECT().Create(gomock.Any()).Return(nil)
				bu.EXPECT().GetAll(gomock.Any()).Return([]brokers.User{}, errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(bb, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Fails to retrieve new user broker",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bu.EXPECT().Create(gomock.Any()).Return(nil)
				bu.EXPECT().GetAll(gomock.Any()).Return([]brokers.User{}, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Succeeded",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bu.EXPECT().Create(gomock.Any()).Return(nil)
				bu.EXPECT().GetAll(gomock.Any()).Return([]brokers.User{{}}, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, bu, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/users/brokers", bytes.NewBuffer(tt.broker))

			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.CreateUserBroker(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestDeleteUserBroker tests the DeleteUserBroker handler
func TestDeleteUserBroker(t *testing.T) {
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
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				m.EXPECT().GetUserFromContext(gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "Fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, false)
				handlers.ReplaceGlobals(m)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fails to verify the user broker existence",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, errors.New("error"))
				bu.EXPECT().Delete(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Fails to find the user broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bu.EXPECT().Delete(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fails to delete the user broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				bu.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				bu.EXPECT().Delete(gomock.Any()).Return(nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/users/brokers/"+uuid.New().String(), nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.DeleteUserBroker(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetUserBrokers tests the GetUserBrokers handler
func TestGetUserBrokers(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, false)
				handlers.ReplaceGlobals(m)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().GetAll(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fails to retrieve all user brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().GetAll(gomock.Any()).Return([]brokers.User{}, errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().GetAll(gomock.Any()).Return([]brokers.User{}, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/brokers", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetUserBrokers(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
