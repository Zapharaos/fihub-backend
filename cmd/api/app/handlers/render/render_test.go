package render

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestError tests the Error function
func TestError(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		err           error
		message       string
		expectMessage string
	}{
		{
			name:    "err is nil",
			err:     nil,
			message: "",
		},
		{
			name:          "message is empty",
			err:           errors.New("test error"),
			message:       "",
			expectMessage: "test error",
		},
		{
			name:          "err and message are set",
			err:           errors.New("test error"),
			message:       "test message",
			expectMessage: "test message: test error",
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			// Call the function
			Error(w, r, tt.err, tt.message)
			resp := w.Result()
			defer resp.Body.Close()

			// Retrieve the response body
			var response ErrorResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)

			// Check the response
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			assert.Equal(t, TitleInternalServerError, response.Title)
			assert.Equal(t, tt.expectMessage, response.Message)
		})
	}
}

// TestBadRequest tests the BadRequest function
func TestBadRequest(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		err           error
		expectMessage string
	}{
		{
			name: "err is nil",
			err:  nil,
		},
		{
			name:          "err is set",
			err:           errors.New("bad request"),
			expectMessage: "bad request",
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			// Call the function
			BadRequest(w, r, tt.err)
			resp := w.Result()
			defer resp.Body.Close()

			// Retrieve the response body
			var response ErrorResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)

			// Check the response
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, TitleBadRequest, response.Title)
			assert.Equal(t, tt.expectMessage, response.Message)
		})
	}
}

// TestOK tests the OK function
func TestOK(t *testing.T) {
	// Create a new recorder
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	// Call the function
	OK(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	// Check the response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestNotImplemented tests the NotImplemented function
func TestNotImplemented(t *testing.T) {
	// Create a new recorder
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	// Call the function
	NotImplemented(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	// Check the response
	assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		assert.Fail(t, "failed to decode response")
	}
	assert.Equal(t, "Not Implemented", response["message"])
}

// TestJSON tests the JSON function
func TestJSON(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		data           interface{}
		expectedStatus int
	}{
		{
			name:           "success",
			data:           map[string]string{"key": "value"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "encode error",
			data:           func() {}, // invalid data type for JSON encoding
			expectedStatus: http.StatusInternalServerError,
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			// Call the function
			JSON(w, r, tt.data)
			resp := w.Result()
			defer resp.Body.Close()

			// Check the response
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.data, response)
			}
		})
	}
}

// TestNotFound tests the NotFound function
func TestNotFound(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		err           error
		expectMessage string
	}{
		{
			name:          "err is nil",
			err:           nil,
			expectMessage: "",
		},
		{
			name:          "err is set",
			err:           errors.New("not found"),
			expectMessage: "not found",
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			// Call the function
			NotFound(w, r, tt.err)
			resp := w.Result()
			defer resp.Body.Close()

			// Check the response
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)

			var response ErrorResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, TitleNotFound, response.Title)
			assert.Equal(t, tt.expectMessage, response.Message)
		})
	}
}

// TestCount tests the Count function
func TestCount(t *testing.T) {
	// Create a new recorder
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	count := int64(42)

	// Call the function
	Count(w, r, count)
	resp := w.Result()
	defer resp.Body.Close()

	// Check the response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response CountResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		assert.Fail(t, "failed to decode response")
	}

	assert.Equal(t, count, response.Count)
}

// TestErrorCodesCodeToHttpCode tests the ErrorCodesCodeToHttpCode function
func TestErrorCodesCodeToHttpCode(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		statusCode     codes.Code
		expectedStatus int
	}{
		{
			name:           "FailedPrecondition",
			statusCode:     codes.FailedPrecondition,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "AlreadyExists",
			statusCode:     codes.AlreadyExists,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "InvalidArgument",
			statusCode:     codes.InvalidArgument,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "NotFound",
			statusCode:     codes.NotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "PermissionDenied",
			statusCode:     codes.PermissionDenied,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Internal",
			statusCode:     codes.Internal,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Unimplemented",
			statusCode:     codes.Unimplemented,
			expectedStatus: http.StatusInternalServerError, // Fallback case
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			// Create a status error with the specified code
			err := status.Error(tt.statusCode, "test error")

			// Call the function
			ErrorCodesCodeToHttpCode(w, r, err)
			resp := w.Result()
			defer resp.Body.Close()

			// Check the response
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// For BadRequest, we should also check the error response
			if tt.statusCode == codes.InvalidArgument {
				var response ErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, TitleBadRequest, response.Title)
				assert.Contains(t, response.Message, "test error")
			}
		})
	}
}
