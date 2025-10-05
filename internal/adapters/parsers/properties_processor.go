package parsers

import (
	"context"
	"fmt"
	"strings"

	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// PropertiesProcessor implements FileProcessor for Properties files
type PropertiesProcessor struct {
	supportedExtensions []string
}

// NewPropertiesProcessor creates a new Properties processor
func NewPropertiesProcessor() *PropertiesProcessor {
	return &PropertiesProcessor{
		supportedExtensions: []string{"properties"},
	}
}

// CanProcess checks if this processor can handle the given filename
func (p *PropertiesProcessor) CanProcess(filename string) bool {
	return ValidateFilenameAndExtension(filename, p.supportedExtensions)
}

// Process processes a Properties file
func (p *PropertiesProcessor) Process(ctx context.Context, filename string, content []byte) (*models.ConfigData, error) {
	if err := ValidateContextAndInput(ctx, filename, content); err != nil {
		return nil, err
	}

	data, err := parsePropertiesContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Properties: %w", err)
	}

	return createConfigData(filename, "properties", data), nil
}

// GetSupportedExtensions returns supported file extensions
func (p *PropertiesProcessor) GetSupportedExtensions() []string {
	return copyExtensions(p.supportedExtensions)
}

// parsePropertiesContent parses Properties content
func parsePropertiesContent(content []byte) (map[string]interface{}, error) {
	// Guard clause: empty content
	if isEmptyContent(content) {
		return createEmptyResult(), nil
	}

	return parseKeyValueContent(content, isPropertiesComment)
}

// isPropertiesComment checks if a line is a comment in Properties format
func isPropertiesComment(line string) bool {
	return strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!")
}
