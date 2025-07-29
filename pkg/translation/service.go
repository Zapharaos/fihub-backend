package translation

import (
	"golang.org/x/text/language"
	"sync"
)

// Service defines the interface for handling translation
type Service interface {
	Localizer(language language.Tag) (interface{}, bool, error)
	Translate(localizer interface{}, message *Message) (string, bool, error)
	MustTranslate(localizer interface{}, message *Message) string
}

// Message represents a translatable message item
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

/*// Global convenience functions that use the singleton service

// T translates a message using the global service
// Returns the translated message or the messageID if translation fails
func T(language language.Tag, message *Message) string {
	service := S()
	if service == nil {
		return message.ID
	}

	localizer, _, err := service.Localizer(language)
	if err != nil {
		return message.ID // Fallback to message ID
	}

	result, _, err := service.Translate(localizer, message)
	if err != nil {
		return message.ID // Fallback to message ID
	}
	return result
}

// MustT translates a message using the global service and panics on error
func MustT(language language.Tag, message *Message) string {
	service := S()
	if service == nil {
		panic("translation service not initialized")
	}

	localizer, _, err := service.Localizer(language)
	if err != nil {
		panic(fmt.Sprintf("failed to get localizer: %v", err))
	}

	result, _, err := service.Translate(localizer, message)
	if err != nil {
		panic(fmt.Sprintf("translation failed: %v", err))
	}
	return result
}*/
