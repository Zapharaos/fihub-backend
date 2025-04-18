package handlers_test

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails at bad image input",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(nil, "", true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to verify broker image existence",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(true, errors.New("error"))
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "broker already has an image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(true, nil)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to create an image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(false, nil)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Return(errors.New("error"))
				bi.EXPECT().Get(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, bi))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve the broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(false, nil)
				bb.EXPECT().SetImage(gomock.Any(), gomock.Any()).Times(0)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, false, errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, bi))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find the broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(false, nil)
				bb.EXPECT().SetImage(gomock.Any(), gomock.Any()).Times(0)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, false, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, bi))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to set the image to broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(false, nil)
				bb.EXPECT().SetImage(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, bi))
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
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(false, nil)
				bb.EXPECT().SetImage(gomock.Any(), gomock.Any()).Return(nil)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(bb, nil, bi))
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
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Get(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to retrieve broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, false, errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, false, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{Data: []byte{0x00, 0x01}}, true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
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
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "image input is invalid",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(nil, "", true)
				handlers.ReplaceGlobals(m)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to verify image broker existence",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, errors.New("error"))
				bi.EXPECT().Update(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find image broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
				bi.EXPECT().Update(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "fails to update image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
				bi.EXPECT().Update(gomock.Any()).Return(errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
				bi.EXPECT().Update(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, false, errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), uuid.New(), true)
				m.EXPECT().ReadImage(gomock.Any(), gomock.Any()).Return(fileData, fileName, true)
				handlers.ReplaceGlobals(m)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
				bi.EXPECT().Update(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, false, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
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
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
				bi.EXPECT().Update(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
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
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
			},
			expectedStatus: http.StatusOK, // should be StatusBadRequest, but not with mock
		},
		{
			name: "fails to verify the image broker existence",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, errors.New("error"))
				bi.EXPECT().Delete(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find the image broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
				bi.EXPECT().Delete(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "fails to delete the broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.EXPECT().ParseUUIDPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, uuid.Nil, true)
				handlers.ReplaceGlobals(m)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
				bi.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
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
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
				bi.EXPECT().Delete(gomock.Any()).Return(nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, nil, bi))
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
