package parsers

import (
	"context"

	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// MockFileProcessor for testing
type MockFileProcessor struct {
	CanProcessFunc     func(filename string) bool
	ProcessFunc        func(ctx context.Context, filename string, content []byte) (*models.ConfigData, error)
	GetExtensionsFunc  func() []string
	ProcessedFiles     []string
	SupportedExts      []string
}

func (m *MockFileProcessor) CanProcess(filename string) bool {
	m.ProcessedFiles = append(m.ProcessedFiles, filename)
	if m.CanProcessFunc != nil {
		return m.CanProcessFunc(filename)
	}
	return true
}

func (m *MockFileProcessor) Process(ctx context.Context, filename string, content []byte) (*models.ConfigData, error) {
	if m.ProcessFunc != nil {
		return m.ProcessFunc(ctx, filename, content)
	}
	return &models.ConfigData{
		Filename: filename,
		Format:   "mock",
		Data:     map[string]interface{}{"mock": "data"},
	}, nil
}

func (m *MockFileProcessor) GetSupportedExtensions() []string {
	if m.GetExtensionsFunc != nil {
		return m.GetExtensionsFunc()
	}
	if m.SupportedExts != nil {
		return m.SupportedExts
	}
	return []string{"mock"}
}

// MockFileReader for testing
type MockFileReader struct {
	ReadFileFunc   func(filename string) ([]byte, error)
	Files          map[string][]byte
	ReadFileCount  int
}

func (m *MockFileReader) ReadFile(filename string) ([]byte, error) {
	m.ReadFileCount++
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(filename)
	}
	if m.Files != nil {
		if content, exists := m.Files[filename]; exists {
			return content, nil
		}
	}
	return []byte("mock content"), nil
}

// TestHelper provides common test utilities
type TestHelper struct{}

// CreateTestContent creates test content for different formats
func (h *TestHelper) CreateTestContent(format string) []byte {
	switch format {
	case "yaml":
		return []byte("app:\n  name: test\n  version: 1.0")
	case "json":
		return []byte(`{"app": {"name": "test", "version": "1.0"}}`)
	case "toml":
		return []byte("[app]\nname = 'test'\nversion = '1.0'")
	case "properties":
		return []byte("app.name=test\napp.version=1.0")
	case "ini":
		return []byte("[app]\nname=test\nversion=1.0")
	case "env":
		return []byte("APP_NAME=test\nAPP_VERSION=1.0")
	default:
		return []byte("mock content")
	}
}

// CreateTestConfigData creates test config data
func (h *TestHelper) CreateTestConfigData(filename, format string) *models.ConfigData {
	return &models.ConfigData{
		Filename: filename,
		Format:   format,
		Data: map[string]interface{}{
			"app": map[string]interface{}{
				"name":    "test",
				"version": "1.0",
			},
		},
	}
}
