package service

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAuthService(t *testing.T) {
	tests := []struct {
		name         string
		appEnv       string
		expectDevKey bool
	}{
		{
			"development",
			"development",
			true,
		},
		{
			"production",
			"production",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("APP_ENV", tt.appEnv)
			service := NewAuthService(nil)

			if tt.expectDevKey {
				assert.Equal(t, []byte("dev-signing-key"), service.signingKey)
			} else {
				assert.NotEqual(t, []byte("dev-signing-key"), service.signingKey)
				assert.Len(t, service.signingKey, 128)
			}
			assert.Equal(t, nil, service.userClient)
		})
	}
}
