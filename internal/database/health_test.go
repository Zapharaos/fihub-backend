package database

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockHealthChecker implements the HealthCheckable interface for testing
type MockHealthChecker struct {
	healthy bool
	mu      sync.Mutex
}

// IsHealthy returns the configured health status
func (m *MockHealthChecker) IsHealthy() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.healthy
}

// SetHealth allows tests to change the health status
func (m *MockHealthChecker) SetHealth(healthy bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.healthy = healthy
}

func TestNewHealthMonitor(t *testing.T) {
	// Test with normal interval
	monitor := NewHealthMonitor(10 * time.Second)
	assert.Equal(t, 10*time.Second, monitor.interval)
	assert.NotNil(t, monitor.targets)
	assert.NotNil(t, monitor.done)
	assert.Len(t, monitor.targets, 0)

	// Test with zero interval (should default to 30s)
	monitor = NewHealthMonitor(0)
	assert.Equal(t, 30*time.Second, monitor.interval)

	// Test with negative interval (should default to 30s)
	monitor = NewHealthMonitor(-5 * time.Second)
	assert.Equal(t, 30*time.Second, monitor.interval)
}

func TestAddTarget(t *testing.T) {
	monitor := NewHealthMonitor(time.Second)
	checker := &MockHealthChecker{healthy: true}

	monitor.AddTarget("test-db", checker, func() {})

	assert.Len(t, monitor.targets, 1)
	assert.Equal(t, "test-db", monitor.targets[0].Name)
	assert.Equal(t, checker, monitor.targets[0].Checker)
	assert.NotNil(t, monitor.targets[0].Recover)
}

func TestHealthMonitorStart(t *testing.T) {
	monitor := NewHealthMonitor(100 * time.Millisecond)

	checker := &MockHealthChecker{healthy: true}
	recoverCalled := false
	recoverFn := func() { recoverCalled = true }

	monitor.AddTarget("test-db", checker, recoverFn)

	// Start monitoring
	stop := monitor.Start()

	// Initially healthy, recover shouldn't be called
	time.Sleep(150 * time.Millisecond)
	assert.False(t, recoverCalled)

	// Make it unhealthy, recover should be called
	checker.SetHealth(false)
	time.Sleep(150 * time.Millisecond)
	assert.True(t, recoverCalled)

	// Stop monitoring
	stop()
}

func TestStartHealthMonitoring(t *testing.T) {
	checker := &MockHealthChecker{healthy: true}
	recoverCalled := false
	recoverFn := func() { recoverCalled = true }

	// Start monitoring with convenience function
	stop := StartHealthMonitoring("test-db", 100*time.Millisecond, checker, recoverFn)

	// Initially healthy, recover shouldn't be called
	time.Sleep(150 * time.Millisecond)
	assert.False(t, recoverCalled)

	// Make it unhealthy, recover should be called
	checker.SetHealth(false)
	time.Sleep(150 * time.Millisecond)
	assert.True(t, recoverCalled)

	// Stop monitoring
	stop()
}
