package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreateBroker tests the function CreateBroker
func TestCreateBroker(t *testing.T) {
	// Prepare data
	invalidBroker := brokers.Broker{}
	invalidBrokerBody, _ := json.Marshal(invalidBroker)
	validBroker := brokers.Broker{
		Name: "name",
	}
	validBrokerBody, _ := json.Marshal(validBroker)

	tests := []struct {
		name           string
		body           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to check permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to decode",
			body: []byte("invalid"),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at bad broker input",
			body: invalidBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to verify broker existence",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, errors.New("error"))
				bb.EXPECT().Create(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "broker already exists",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(true, nil)
				bb.EXPECT().Create(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to create broker",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Create(gomock.Any()).Return(uuid.Nil, errors.New("error"))
				bb.EXPECT().Get(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve the broker",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Create(gomock.Any()).Return(uuid.Nil, nil)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, false, errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find broker",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Create(gomock.Any()).Return(uuid.Nil, nil)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, false, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Create(gomock.Any()).Return(uuid.Nil, nil)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/brokers", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.CreateBroker(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetBroker tests the function GetBroker
func TestGetBroker(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to retrieve the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, false, errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, false, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/brokers/{id}", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetBroker(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestUpdateBroker tests the function UpdateBroker
func TestUpdateBroker(t *testing.T) {
	// Prepare data
	invalidBroker := brokers.Broker{}
	invalidBrokerBody, _ := json.Marshal(invalidBroker)
	validBroker := brokers.Broker{
		Name: "name",
	}
	validBrokerBody, _ := json.Marshal(validBroker)

	tests := []struct {
		name           string
		body           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to check permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to decode",
			body: []byte("invalid"),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at bad broker input",
			body: invalidBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to retrieve broker",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, false, errors.New("error"))
				bb.EXPECT().ExistsByName(gomock.Any()).Times(0)
				bb.EXPECT().Update(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find broker",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, false, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Times(0)
				bb.EXPECT().Update(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "fails to verify new name usage",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(true, errors.New("error"))
				bb.EXPECT().Update(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "new name already in use",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(true, nil)
				bb.EXPECT().Update(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to update broker",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Update(gomock.Any()).Return(errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve the broker",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Update(gomock.Any()).Return(nil)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, false, errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find broker",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Update(gomock.Any()).Return(nil)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, false, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Update(gomock.Any()).Return(nil)
				bb.EXPECT().Get(gomock.Any()).Return(brokers.Broker{}, true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/brokers", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.UpdateBroker(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestDeleteBroker tests the function DeleteBroker
func TestDeleteBroker(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to check permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to retrieve the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Return(false, errors.New("error"))
				bb.EXPECT().Delete(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bb.EXPECT().Delete(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "fails to delete the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Return(true, nil)
				bb.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Return(true, nil)
				bb.EXPECT().Delete(gomock.Any()).Return(nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/brokers/{id}", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.DeleteBroker(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetBrokers tests the function GetBrokers
func TestGetBrokers(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamBool(gomock.Any(), gomock.Any(), "enabled").Return(false, false)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAllEnabled().Times(0)
				bb.EXPECT().GetAll().Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to retrieve all enabled brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamBool(gomock.Any(), gomock.Any(), "enabled").Return(true, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAllEnabled().Return(nil, errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve all brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamBool(gomock.Any(), gomock.Any(), "enabled").Return(false, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAll().Return(nil, errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamBool(gomock.Any(), gomock.Any(), "enabled").Return(false, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAll().Return([]brokers.Broker{}, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/brokers", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetBrokers(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
