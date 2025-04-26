package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreateBroker tests the function CreateBroker
func TestCreateBroker(t *testing.T) {
	// Prepare data
	validBroker := models.Broker{
		Name: "name",
	}
	validBrokerBody, _ := json.Marshal(validBroker)
	validResponse := &protogen.CreateBrokerResponse{
		Broker: &protogen.Broker{
			Id:       uuid.New().String(),
			Name:     validBroker.Name,
			Disabled: validBroker.Disabled,
		},
	}

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
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBroker(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
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
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBroker(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
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
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBroker(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
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
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBroker(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", apiBasePath+"/brokers", bytes.NewBuffer(tt.body))

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
	validBroker := models.Broker{
		Name:     "name",
		Disabled: false,
	}
	validResponse := &protogen.GetBrokerResponse{
		Broker: &protogen.Broker{
			Id:       uuid.New().String(),
			Name:     validBroker.Name,
			Disabled: validBroker.Disabled,
		},
	}

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
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().GetBroker(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to retrieve the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().GetBroker(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().GetBroker(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/brokers/{id}", nil)

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
	validBroker := models.Broker{
		Name: "name",
	}
	validBrokerBody, _ := json.Marshal(validBroker)
	validResponse := &protogen.UpdateBrokerResponse{
		Broker: &protogen.Broker{
			Id:       uuid.New().String(),
			Name:     validBroker.Name,
			Disabled: validBroker.Disabled,
		},
	}

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
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().UpdateBroker(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
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
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().UpdateBroker(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
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
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().UpdateBroker(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
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
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().UpdateBroker(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", apiBasePath+"/brokers", bytes.NewBuffer(tt.body))

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
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBroker(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to delete the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBroker(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
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
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBroker(gomock.Any(), gomock.Any()).Return(&protogen.DeleteBrokerResponse{
					Success: true,
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", apiBasePath+"/brokers/{id}", nil)

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

// TestListBrokers tests the function ListBrokers
func TestListBrokers(t *testing.T) {
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
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().ListBrokers(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to list brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamBool(gomock.Any(), gomock.Any(), "enabled").Return(false, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().ListBrokers(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamBool(gomock.Any(), gomock.Any(), "enabled").Return(false, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().ListBrokers(gomock.Any(), gomock.Any()).Return(&protogen.ListBrokersResponse{
					Brokers: []*protogen.Broker{},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/brokers", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.ListBrokers(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
