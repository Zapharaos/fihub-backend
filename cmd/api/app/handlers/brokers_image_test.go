package handlers_test

import (
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
	"strings"
	"testing"
)

// TestCreateBrokerImage tests the CreateBrokerImage function
func TestCreateBrokerImage(t *testing.T) {
	// Prepare data
	fileData := []byte{0x00, 0x01, 0x02, 0x03}
	fileName := strings.Repeat("a", models.ImageNameMinLength)

	// Prepare tests
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
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to read image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(nil, "", false)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBrokerImage(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to set the image to broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBrokerImage(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBrokerImage(gomock.Any(), gomock.Any()).Return(&protogen.CreateBrokerImageResponse{
					Image: &protogen.BrokerImage{
						Id:       uuid.New().String(),
						BrokerId: uuid.New().String(),
						Name:     fileName,
						Data:     fileData,
					},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/brokers/{id}/image", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.CreateBrokerImage(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetBrokerImage tests the GetBrokerImage function
func TestGetBrokerImage(t *testing.T) {
	// Prepare tests
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to parse params",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerImage(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to retrieve broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerImage(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerImage(gomock.Any(), gomock.Any()).Return(&protogen.GetBrokerImageResponse{
					Name: "name",
					Data: []byte{0x00, 0x01},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/brokers/{id}/image/{image_id}", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetBrokerImage(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestUpdateBrokerImage tests the UpdateBrokerImage function
func TestUpdateBrokerImage(t *testing.T) {
	// Prepare data
	fileData := []byte{0x00, 0x01, 0x02, 0x03}
	fileName := strings.Repeat("a", models.ImageNameMinLength)

	// Prepare tests
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
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to parse params",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, uuid.Nil, false)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to read image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(nil, "", false)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().UpdateBrokerImage(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to update image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().UpdateBrokerImage(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().UpdateBrokerImage(gomock.Any(), gomock.Any()).Return(&protogen.UpdateBrokerImageResponse{
					Image: &protogen.BrokerImage{
						Id:       uuid.New().String(),
						BrokerId: uuid.New().String(),
						Name:     fileName,
						Data:     fileData,
					},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", apiBasePath+"/brokers/{id}/image/{image_id}", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.UpdateBrokerImage(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestDeleteBrokerImage tests the DeleteBrokerImage function
func TestDeleteBrokerImage(t *testing.T) {
	// Prepare tests
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
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be StatusUnauthorized, but not with mock
		},
		{
			name: "fails to parse params",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBrokerImage(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to delete the broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBrokerImage(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBrokerImage(gomock.Any(), gomock.Any()).Return(&protogen.DeleteBrokerImageResponse{
					Success: true,
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, bc, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", apiBasePath+"/brokers/{id}/image/{image_id}", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.DeleteBrokerImage(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
