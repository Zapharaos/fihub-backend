package clients

import (
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestReplaceGlobals(t *testing.T) {
	original := C()

	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockHealthServiceClient(ctrl)
	newClient := NewClients(WithHealthClient(mocks.NewMockHealthServiceClient(ctrl)))
	restore := ReplaceGlobals(newClient)

	assert.Equal(t, mockService, C().Health(), "expected global health client to be replaced")

	restore() // Revert to original

	assert.Equal(t, original.Health(), C().Health(), "expected global health client to be restored")
}

func TestNewClients_WithServiceOptions(t *testing.T) {
	t.Run("Health client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockHealthServiceClient(ctrl)
		c := NewClients(WithHealthClient(mockService))
		assert.Equal(t, mockService, c.Health())
	})

	t.Run("User client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockUserServiceClient(ctrl)
		c := NewClients(WithUserClient(mockService))
		assert.Equal(t, mockService, c.User())
	})

	t.Run("Auth client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockAuthServiceClient(ctrl)
		c := NewClients(WithAuthClient(mockService))
		assert.Equal(t, mockService, c.Auth())
	})

	t.Run("Security client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockSecurityServiceClient(ctrl)
		c := NewClients(WithSecurityClient(mockService))
		assert.Equal(t, mockService, c.Security())
	})

	t.Run("Broker client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockBrokerServiceClient(ctrl)
		c := NewClients(WithBrokerClient(mockService))
		assert.Equal(t, mockService, c.Broker())
	})

	t.Run("Transaction client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockTransactionServiceClient(ctrl)
		c := NewClients(WithTransactionClient(mockService))
		assert.Equal(t, mockService, c.Transaction())
	})
}
