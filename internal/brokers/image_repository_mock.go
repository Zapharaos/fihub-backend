package brokers

import "github.com/google/uuid"

// MockImageRepository represents a mock ImageRepository
type MockImageRepository struct {
	image Image
	error error
	found bool
}

// NewMockImageRepository creates a new MockImageRepository of the ImageRepository interface
func NewMockImageRepository() ImageRepository {
	r := MockImageRepository{}
	var repo ImageRepository = &r
	return repo
}

func (m MockImageRepository) Create(_ Image) error {
	return m.error
}

func (m MockImageRepository) Get(_ uuid.UUID) (Image, bool, error) {
	return m.image, m.found, m.error
}

func (m MockImageRepository) Update(_ Image) error {
	return m.error
}

func (m MockImageRepository) Delete(_ uuid.UUID) error {
	return m.error
}

func (m MockImageRepository) Exists(_ uuid.UUID, _ uuid.UUID) (bool, error) {
	return m.found, m.error
}
