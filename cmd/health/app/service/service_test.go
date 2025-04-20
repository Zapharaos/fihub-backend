package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/protogen"
	"testing"
)

// TestCheckHealth tests the CheckHealth method of the Service struct.
func TestCheckHealth(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name        string
		request     *protogen.HealthRequest
		expected    *protogen.HealthResponse
		expectError bool
	}{
		{
			name: "Valid service name",
			request: &protogen.HealthRequest{
				ServiceName: "TestService",
			},
			expected: &protogen.HealthResponse{
				IsHealthy: true,
				Message:   "Service is healthy",
			},
			expectError: false,
		},
		{
			name: "Empty service name",
			request: &protogen.HealthRequest{
				ServiceName: "",
			},
			expected: &protogen.HealthResponse{
				IsHealthy: false,
				Message:   "Service name is required",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.CheckHealth(context.Background(), tt.request)
			if (err != nil) != tt.expectError {
				t.Errorf("CheckHealth() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if response.IsHealthy != tt.expected.IsHealthy || response.Message != tt.expected.Message {
				t.Errorf("CheckHealth() = %v, expected %v", response, tt.expected)
			}
		})
	}
}
