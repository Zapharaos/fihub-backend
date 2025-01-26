package email

import (
	"sync"
)

// Service defines the interface for handling emails
type Service interface {
	Send(emailTo, subject, plainTextContent, htmlContent string) error
}

var (
	_globalServiceMu sync.RWMutex
	_globalService   Service
)

// S is used to access the global service singleton
func S() Service {
	_globalServiceMu.RLock()
	defer _globalServiceMu.RUnlock()

	service := _globalService
	return service
}

// ReplaceGlobals affect a new repository to the global service singleton
func ReplaceGlobals(service Service) func() {
	_globalServiceMu.Lock()
	defer _globalServiceMu.Unlock()

	prev := _globalService
	_globalService = service
	return func() { ReplaceGlobals(prev) }
}
