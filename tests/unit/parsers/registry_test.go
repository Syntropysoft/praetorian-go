package parsers

import (
	"testing"

	"github.com/syntropysoft/praetorian-go/internal/adapters/parsers"
)

// TestParserRegistry tests the parser registry functionality
func TestParserRegistry(t *testing.T) {
	registry := parsers.NewParserRegistry()

	t.Run("should have processors registered", func(t *testing.T) {
		// Test by trying to get a processor for a known file type
		processor, err := registry.GetProcessor("test.yaml")
		if err != nil {
			t.Errorf("Expected to find YAML processor, got error: %v", err)
		}
		if processor == nil {
			t.Error("Expected processor to be non-nil")
		}
	})

	t.Run("should get processor for valid extension", func(t *testing.T) {
		processor, err := registry.GetProcessor("test.yaml")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if processor == nil {
			t.Error("Expected processor, got nil")
		}
	})

	t.Run("should return error for invalid extension", func(t *testing.T) {
		_, err := registry.GetProcessor("test.invalid")
		if err == nil {
			t.Error("Expected error for invalid extension")
		}
	})

	t.Run("should get processor for file", func(t *testing.T) {
		processor, err := registry.GetProcessor("test.yaml")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if processor == nil {
			t.Error("Expected processor, got nil")
		}
	})
}

// TestRegistryIntegration tests registry integration with real processors
func TestRegistryIntegration(t *testing.T) {
	registry := parsers.NewParserRegistry()

	t.Run("should support multiple file types", func(t *testing.T) {
		fileTypes := []string{"test.yaml", "test.json", "test.toml", "test.properties", "test.ini", "test.hcl", "test.xml", "test.env"}
		
		for _, filename := range fileTypes {
			processor, err := registry.GetProcessor(filename)
			if err != nil {
				t.Errorf("Expected to find processor for %s, got error: %v", filename, err)
			}
			if processor == nil {
				t.Errorf("Expected processor for %s to be non-nil", filename)
			}
		}
	})
}

// BenchmarkRegistryGetProcessor benchmarks processor retrieval
func BenchmarkRegistryGetProcessor(b *testing.B) {
	registry := parsers.NewParserRegistry()
	for i := 0; i < b.N; i++ {
		registry.GetProcessor("test.yaml")
	}
}
