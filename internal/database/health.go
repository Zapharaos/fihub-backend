package database

import (
	"go.uber.org/zap"
	"time"
)

// HealthCheckable represents any database connection that can be health-checked
type HealthCheckable interface {
	IsHealthy() bool
}

// HealthTarget represents a registered health target with its settings
type HealthTarget struct {
	Name    string
	Checker HealthCheckable
	Recover func()
}

// HealthMonitor manages multiple health targets with a single goroutine
type HealthMonitor struct {
	interval time.Duration
	targets  []HealthTarget
	ticker   *time.Ticker
	done     chan bool
}

// NewHealthMonitor creates a new health check manager with the given interval
func NewHealthMonitor(interval time.Duration) *HealthMonitor {
	if interval <= 0 {
		interval = 30 * time.Second
	}

	return &HealthMonitor{
		interval: interval,
		targets:  make([]HealthTarget, 0),
		done:     make(chan bool),
	}
}

// AddTarget registers a new health target with the monitor
func (m *HealthMonitor) AddTarget(name string, checker HealthCheckable, recover func()) {
	m.targets = append(m.targets, HealthTarget{
		Name:    name,
		Checker: checker,
		Recover: recover,
	})
}

// Start begins the health check monitoring for all registered checkers
func (m *HealthMonitor) Start() func() {
	m.ticker = time.NewTicker(m.interval)

	go func() {
		for {
			select {
			case <-m.done:
				m.ticker.Stop()
				return
			case <-m.ticker.C:
				for _, entry := range m.targets {
					if !entry.Checker.IsHealthy() {
						zap.L().Info("Attempting recovery for target", zap.String("target", entry.Name))
						entry.Recover()
					}
				}
			}
		}
	}()

	// Return function to stop the health check
	return func() {
		m.done <- true
	}
}

// StartHealthMonitoring is a convenience function to create a manager, add a target, and start monitoring
func StartHealthMonitoring(name string, interval time.Duration, checker HealthCheckable, recover func()) func() {
	manager := NewHealthMonitor(interval)
	manager.AddTarget(name, checker, recover)
	return manager.Start()
}
