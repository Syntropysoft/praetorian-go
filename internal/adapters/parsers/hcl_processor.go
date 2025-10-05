package parsers

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// HCLProcessor implements FileProcessor for HCL files
type HCLProcessor struct {
	supportedExtensions []string
}

// NewHCLProcessor creates a new HCL processor
func NewHCLProcessor() *HCLProcessor {
	return &HCLProcessor{
		supportedExtensions: []string{"hcl"},
	}
}

// CanProcess checks if this processor can handle the given filename
func (p *HCLProcessor) CanProcess(filename string) bool {
	return ValidateFilenameAndExtension(filename, p.supportedExtensions)
}

// Process processes an HCL file
func (p *HCLProcessor) Process(ctx context.Context, filename string, content []byte) (*models.ConfigData, error) {
	if err := ValidateContextAndInput(ctx, filename, content); err != nil {
		return nil, err
	}

	data, err := parseHCLContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HCL: %w", err)
	}

	return createConfigData(filename, "hcl", data), nil
}

// GetSupportedExtensions returns supported file extensions
func (p *HCLProcessor) GetSupportedExtensions() []string {
	return copyExtensions(p.supportedExtensions)
}

// parseHCLContent parses HCL content
func parseHCLContent(content []byte) (map[string]interface{}, error) {
	// Guard clause: empty content
	if isEmptyContent(content) {
		return createEmptyResult(), nil
	}

	var result map[string]interface{}
	// hclsimple.Decode requires a filename for context, even if it's a dummy
	if err := hclsimple.Decode("config.hcl", content, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to parse HCL: %w", err)
	}

	return result, nil
}
