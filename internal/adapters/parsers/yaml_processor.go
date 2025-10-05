package parsers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// YAMLProcessor implements FileProcessor for YAML files
type YAMLProcessor struct {
	supportedExtensions []string
}

// NewYAMLProcessor creates a new YAML processor
func NewYAMLProcessor() *YAMLProcessor {
	return &YAMLProcessor{
		supportedExtensions: []string{"yaml", "yml"},
	}
}

// CanProcess checks if this processor can handle the given filename
func (p *YAMLProcessor) CanProcess(filename string) bool {
	// Guard clause: validate filename
	if filename == "" {
		return false
	}

	ext := p.getFileExtension(filename)
	return p.supportsExtension(ext)
}

// Process processes a YAML file
func (p *YAMLProcessor) Process(ctx context.Context, filename string, content []byte) (*models.ConfigData, error) {
	// Guard clause: validate context
	if ctx.Err() != nil {
		return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
	}

	// Guard clause: validate input
	if err := p.validateInput(filename, content); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Parse YAML content
	data, err := p.parseYAML(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Create config data
	return p.createConfigData(filename, data), nil
}

// GetSupportedExtensions returns supported file extensions
func (p *YAMLProcessor) GetSupportedExtensions() []string {
	// Return a copy to prevent external modification
	extensions := make([]string, len(p.supportedExtensions))
	copy(extensions, p.supportedExtensions)
	return extensions
}

// parseYAML parses YAML content into a map
func (p *YAMLProcessor) parseYAML(content []byte) (map[string]interface{}, error) {
	// Guard clause: empty content
	if len(content) == 0 {
		return make(map[string]interface{}), nil
	}

	var result map[string]interface{}
	if err := yaml.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("YAML unmarshal failed: %w", err)
	}

	// Handle null result
	if result == nil {
		return make(map[string]interface{}), nil
	}

	return result, nil
}

// createConfigData creates ConfigData from parsed content
func (p *YAMLProcessor) createConfigData(filename string, data map[string]interface{}) *models.ConfigData {
	return &models.ConfigData{
		Filename:  filename,
		Format:    "yaml",
		Data:      data,
		Metadata:  p.createMetadata(filename, data),
		Timestamp: time.Now(),
	}
}

// createMetadata creates metadata for the config data
func (p *YAMLProcessor) createMetadata(filename string, data map[string]interface{}) map[string]interface{} {
	metadata := map[string]interface{}{
		"processor":    "yaml",
		"file_size":    len(data),
		"key_count":    p.countKeys(data),
		"has_nested":   p.hasNestedData(data),
	}

	// Add file-specific metadata
	if filename != "" {
		metadata["basename"] = filepath.Base(filename)
		metadata["dirname"] = filepath.Dir(filename)
	}

	return metadata
}

// validateInput validates input parameters
func (p *YAMLProcessor) validateInput(filename string, content []byte) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if content == nil {
		return fmt.Errorf("content cannot be nil")
	}
	return nil
}

// getFileExtension extracts extension from filename
func (p *YAMLProcessor) getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.TrimPrefix(strings.ToLower(ext), ".")
}

// supportsExtension checks if extension is supported
func (p *YAMLProcessor) supportsExtension(extension string) bool {
	for _, ext := range p.supportedExtensions {
		if ext == extension {
			return true
		}
	}
	return false
}

// countKeys counts the number of top-level keys
func (p *YAMLProcessor) countKeys(data map[string]interface{}) int {
	if data == nil {
		return 0
	}
	return len(data)
}

// hasNestedData checks if the data contains nested structures
func (p *YAMLProcessor) hasNestedData(data map[string]interface{}) bool {
	if data == nil {
		return false
	}

	for _, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			return true
		case []interface{}:
			if len(v) > 0 {
				return true
			}
		}
	}

	return false
}
