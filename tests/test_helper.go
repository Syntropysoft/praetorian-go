package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestHelper provides common testing utilities
type TestHelper struct {
	TempDir string
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	tempDir := t.TempDir()
	return &TestHelper{TempDir: tempDir}
}

// CreateTestFile creates a test file with content
func (h *TestHelper) CreateTestFile(filename, content string) error {
	fullPath := filepath.Join(h.TempDir, filename)
	dir := filepath.Dir(fullPath)
	
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// GetTestFilePath returns the full path for a test file
func (h *TestHelper) GetTestFilePath(filename string) string {
	return filepath.Join(h.TempDir, filename)
}

// FileExists checks if a test file exists
func (h *TestHelper) FileExists(filename string) bool {
	fullPath := h.GetTestFilePath(filename)
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}

// ReadTestFile reads content from a test file
func (h *TestHelper) ReadTestFile(filename string) (string, error) {
	fullPath := h.GetTestFilePath(filename)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

// AssertFileContains asserts that a test file contains specific content
func (h *TestHelper) AssertFileContains(t *testing.T, filename, expectedContent string) {
	t.Helper()
	
	content, err := h.ReadTestFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", filename, err)
	}
	
	if !strings.Contains(content, expectedContent) {
		t.Errorf("File %s does not contain expected content '%s'", filename, expectedContent)
	}
}

// AssertFileExists asserts that a test file exists
func (h *TestHelper) AssertFileExists(t *testing.T, filename string) {
	t.Helper()
	
	if !h.FileExists(filename) {
		t.Errorf("Expected file %s to exist", filename)
	}
}

// AssertFileNotExists asserts that a test file does not exist
func (h *TestHelper) AssertFileNotExists(t *testing.T, filename string) {
	t.Helper()
	
	if h.FileExists(filename) {
		t.Errorf("Expected file %s to not exist", filename)
	}
}

// Cleanup removes all test files
func (h *TestHelper) Cleanup() error {
	return os.RemoveAll(h.TempDir)
}

// Common test content generators
func (h *TestHelper) GenerateYAMLConfig() string {
	return `version: "1.0"
app:
  name: "test-app"
  version: "1.0.0"
database:
  host: "localhost"
  port: 5432
`
}

func (h *TestHelper) GenerateJSONConfig() string {
	return `{
  "version": "1.0",
  "app": {
    "name": "test-app",
    "version": "1.0.0"
  },
  "database": {
    "host": "localhost",
    "port": 5432
  }
}`
}

func (h *TestHelper) GeneratePropertiesConfig() string {
	return `app.name=test-app
app.version=1.0.0
database.host=localhost
database.port=5432
`
}
