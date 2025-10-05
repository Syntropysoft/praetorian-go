package parsers

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// XMLProcessor implements FileProcessor for XML files
type XMLProcessor struct {
	supportedExtensions []string
}

// NewXMLProcessor creates a new XML processor
func NewXMLProcessor() *XMLProcessor {
	return &XMLProcessor{
		supportedExtensions: []string{"xml"},
	}
}

// CanProcess checks if this processor can handle the given filename
func (p *XMLProcessor) CanProcess(filename string) bool {
	return ValidateFilenameAndExtension(filename, p.supportedExtensions)
}

// Process processes an XML file
func (p *XMLProcessor) Process(ctx context.Context, filename string, content []byte) (*models.ConfigData, error) {
	if err := ValidateContextAndInput(ctx, filename, content); err != nil {
		return nil, err
	}

	data, err := parseXMLContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return createConfigData(filename, "xml", data), nil
}

// GetSupportedExtensions returns supported file extensions
func (p *XMLProcessor) GetSupportedExtensions() []string {
	return copyExtensions(p.supportedExtensions)
}

// parseXMLContent parses XML content
func parseXMLContent(content []byte) (map[string]interface{}, error) {
	// Guard clause: empty content
	if isEmptyContent(content) {
		return createEmptyResult(), nil
	}

	var data map[string]interface{}
	decoder := xml.NewDecoder(strings.NewReader(string(content)))
	err := decoder.Decode(&data) // This will decode into a generic map
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return data, nil
}
