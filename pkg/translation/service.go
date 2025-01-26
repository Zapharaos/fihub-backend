package translation

import (
	"golang.org/x/text/language"
	"sync"
)

// Service defines the interface for handling translation
type Service interface {
	Localizer(language language.Tag) (interface{}, error)
	Message(localizer interface{}, message Message) string
}

// Message represents a message to be translated
type Message struct {
	ID          string
	Data        interface{}
	PluralCount interface{}
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
