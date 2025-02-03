package templates

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"testing"
)

// TestTemplate_Render tests the Render method of the Template struct
func TestTemplate_Render(t *testing.T) {
	t.Run("empty template", func(t *testing.T) {
		tmpl := Template{}

		result, err := tmpl.Render()
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("template with incorrect content", func(t *testing.T) {
		tmpl := Template{
			Name:       "test_template",
			ContentRaw: "{{.Content",
			Data:       nil,
		}

		result, err := tmpl.Render()
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("template without data", func(t *testing.T) {
		tmpl := Template{
			Name:       "test_template",
			ContentRaw: "{{.Content}}",
			Data:       nil,
		}

		result, err := tmpl.Render()
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("template with incorrect data", func(t *testing.T) {
		tmpl := Template{
			Name:       "test_template",
			ContentRaw: "{{.Content}}",
			Data:       map[string]interface{}{"Content": func() {}},
		}

		result, err := tmpl.Render()
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("template with wrong data", func(t *testing.T) {
		tmpl := Template{
			Name:       "test_template",
			ContentRaw: "{{.Content}}",
			Data:       map[string]string{"Wrong": "wrong"},
		}

		result, err := tmpl.Render()
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("successful render", func(t *testing.T) {
		tmpl := Template{
			Name:       "test_template",
			ContentRaw: "{{.Content}}",
			Data:       map[string]string{"Content": "Hello World"},
		}

		result, err := tmpl.Render()
		assert.NoError(t, err)
		assert.Equal(t, "Hello World", result)
	})
}

func TestTemplate_Build(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	t.Run("content template render", func(t *testing.T) {
		tmpl := Template{
			Name:       "test_template",
			ContentRaw: "{{.Name",
			Data:       nil,
		}

		labels := LayoutLabels{
			Help:       "Help",
			Copyrights: "Copyrights",
		}

		result, err := tmpl.Build(labels)
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("successful build", func(t *testing.T) {
		tmpl := Template{
			Name:       "test_template",
			ContentRaw: "{{.Content}}",
			Data:       map[string]string{"Content": "Hello World"},
		}

		labels := LayoutLabels{
			Help:       "Help",
			Copyrights: "Copyrights",
		}

		result, err := tmpl.Build(labels)
		assert.NoError(t, err)
		assert.Contains(t, result, "Hello World")
		assert.Contains(t, result, "Help")
		assert.Contains(t, result, "Copyrights")
	})
}
