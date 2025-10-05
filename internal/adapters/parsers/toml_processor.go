package parsers

import (
	"context"
	"fmt"

	"github.com/BurntSushi/toml"

	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// TOMLProcessor implements FileProcessor for TOML files
type TOMLProcessor struct {
	supportedExtensions []string
}

// NewTOMLProcessor creates a new TOML processor
func NewTOMLProcessor() *TOMLProcessor {
	return &TOMLProcessor{
		supportedExtensions: []string{"toml"},
	}
}

// CanProcess checks if this processor can handle the given filename
func (p *TOMLProcessor) CanProcess(filename string) bool {
	return ValidateFilenameAndExtension(filename, p.supportedExtensions)
}

// Process processes a TOML file
func (p *TOMLProcessor) Process(ctx context.Context, filename string, content []byte) (*models.ConfigData, error) {
	if err := ValidateContextAndInput(ctx, filename, content); err != nil {
		return nil, err
	}

	data, err := parseTOMLContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	return createConfigData(filename, "toml", data), nil
}

// GetSupportedExtensions returns supported file extensions
func (p *TOMLProcessor) GetSupportedExtensions() []string {
	return copyExtensions(p.supportedExtensions)
}

// parseTOMLContent parses TOML content into a map
func parseTOMLContent(content []byte) (map[string]interface{}, error) {
	if len(content) == 0 {
		return make(map[string]interface{}), nil
	}

	var result map[string]interface{}
	if err := toml.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("TOML unmarshal failed: %w", err)
	}

	if result == nil {
		return make(map[string]interface{}), nil
	}

	return result, nil
}

