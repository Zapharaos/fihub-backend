package clients

import (
	"reflect"
	"testing"
)

func TestNewClients(t *testing.T) {
	c := NewClients()
	if c.services == nil {
		t.Fatal("Expected services map to be initialized, got nil")
	}
	if len(c.services) != 0 {
		t.Errorf("Expected empty services map, got map with %d items", len(c.services))
	}
}

func TestClientsRegisterAndGet(t *testing.T) {
	c := NewClients()

	// Test registering and getting a service
	c.Register("service1", "test-service")

	value, exists := c.Get("service1")
	if !exists {
		t.Error("Expected service to exist, but it doesn't")
	}

	strValue, ok := value.(string)
	if !ok {
		t.Error("Expected string type, got different type")
	}

	if strValue != "test-service" {
		t.Errorf("Expected 'test-service', got '%s'", strValue)
	}

	// Test with non-existent service
	_, exists = c.Get("non-existent")
	if exists {
		t.Error("Expected non-existent service to return false, but it returned true")
	}
}

func TestGetTypedClient(t *testing.T) {
	c := NewClients()

	// Register different types of clients
	c.Register("string-service", "test-string")
	c.Register("int-service", 42)

	// Test successful retrieval with correct type
	strClient, err := GetTypedClient[string](c, "string-service")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if strClient != "test-string" {
		t.Errorf("Expected 'test-string', got '%s'", strClient)
	}

	// Test successful retrieval with another type
	intClient, err := GetTypedClient[int](c, "int-service")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if intClient != 42 {
		t.Errorf("Expected 42, got %d", intClient)
	}

	// Test service not found
	_, err = GetTypedClient[string](c, "non-existent")
	if err == nil || err.Error() != "service non-existent not registered" {
		t.Errorf("Expected error about service not registered, got: %v", err)
	}

	// Test wrong type
	_, err = GetTypedClient[int](c, "string-service")
	if err == nil || err.Error() != "service string-service has wrong type" {
		t.Errorf("Expected error about wrong type, got: %v", err)
	}
}

func TestGlobalClients(t *testing.T) {
	// Save original state to restore later
	originalGlobals := _globalClients
	defer func() { _globalClients = originalGlobals }()

	// Initialize with new clients
	newClients := NewClients()
	newClients.Register("test-service", "test-value")

	// Replace globals and get restore function
	restore := ReplaceGlobals(newClients)

	// Verify global state changed
	updatedGlobals := C()
	value, exists := updatedGlobals.Get("test-service")
	if !exists {
		t.Error("Expected service to exist in global clients after replacement")
	}

	strValue, ok := value.(string)
	if !ok || strValue != "test-value" {
		t.Errorf("Expected 'test-value', got '%v'", value)
	}

	// Test restore function
	restore()

	// Verify original state is restored
	restoredGlobals := C()
	if !reflect.DeepEqual(restoredGlobals, originalGlobals) {
		t.Error("Expected globals to be restored to original state")
	}
}
