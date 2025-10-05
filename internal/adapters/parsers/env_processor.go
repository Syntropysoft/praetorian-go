package parsers

import (
	"context"
	"fmt"
	"strings"

	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// ENVProcessor implements FileProcessor for ENV files
type ENVProcessor struct {
	supportedExtensions []string
}

// NewENVProcessor creates a new ENV processor
func NewENVProcessor() *ENVProcessor {
	return &ENVProcessor{
		supportedExtensions: []string{"env"},
	}
}

// CanProcess checks if this processor can handle the given filename
func (p *ENVProcessor) CanProcess(filename string) bool {
	return ValidateFilenameAndExtension(filename, p.supportedExtensions)
}

// Process processes an ENV file
func (p *ENVProcessor) Process(ctx context.Context, filename string, content []byte) (*models.ConfigData, error) {
	if err := ValidateContextAndInput(ctx, filename, content); err != nil {
		return nil, err
	}

	data, err := parseENVContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ENV: %w", err)
	}

	return createConfigData(filename, "env", data), nil
}

// GetSupportedExtensions returns supported file extensions
func (p *ENVProcessor) GetSupportedExtensions() []string {
	return copyExtensions(p.supportedExtensions)
}

// parseENVContent parses ENV content
func parseENVContent(content []byte) (map[string]interface{}, error) {
	// Guard clause: empty content
	if isEmptyContent(content) {
		return createEmptyResult(), nil
	}

	return parseKeyValueContent(content, isENVComment)
}

// isENVComment checks if a line is a comment in ENV format
func isENVComment(line string) bool {
	return strings.HasPrefix(line, "#")
}
