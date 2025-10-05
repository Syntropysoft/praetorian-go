package parsers

import (
	"context"
	"fmt"
	"strings"

	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// INIProcessor implements FileProcessor for INI files
type INIProcessor struct {
	supportedExtensions []string
}

// NewINIProcessor creates a new INI processor
func NewINIProcessor() *INIProcessor {
	return &INIProcessor{
		supportedExtensions: []string{"ini"},
	}
}

// CanProcess checks if this processor can handle the given filename
func (p *INIProcessor) CanProcess(filename string) bool {
	return ValidateFilenameAndExtension(filename, p.supportedExtensions)
}

// Process processes an INI file
func (p *INIProcessor) Process(ctx context.Context, filename string, content []byte) (*models.ConfigData, error) {
	if err := ValidateContextAndInput(ctx, filename, content); err != nil {
		return nil, err
	}

	data, err := parseINIContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse INI: %w", err)
	}

	return createConfigData(filename, "ini", data), nil
}

// GetSupportedExtensions returns supported file extensions
func (p *INIProcessor) GetSupportedExtensions() []string {
	return copyExtensions(p.supportedExtensions)
}

// parseINIContent parses INI content
func parseINIContent(content []byte) (map[string]interface{}, error) {
	// Guard clause: empty content
	if isEmptyContent(content) {
		return createEmptyResult(), nil
	}

	return parseINIContentWithSections(content), nil
}

// parseINIContentWithSections parses INI content with sections
func parseINIContentWithSections(content []byte) map[string]interface{} {
	result := createEmptyResult()
	currentSection := ""
	lines := splitLines(content)

	for _, line := range lines {
		trimmedLine := trimLine(line)
		
		// Skip empty lines and comments
		if isEmptyLine(trimmedLine) || isINIComment(trimmedLine) {
			continue
		}

		// Check for section header
		if isINISection(trimmedLine) {
			currentSection = extractINISection(trimmedLine)
			result[currentSection] = createEmptyResult()
			continue
		}

		// Parse key-value pair
		if key, value, ok := parseINIKeyValue(trimmedLine); ok {
			setINIValue(result, currentSection, key, value)
		}
	}

	return result
}

// isINIComment checks if a line is a comment in INI format
func isINIComment(line string) bool {
	return strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#")
}

// isINISection checks if a line is a section header
func isINISection(line string) bool {
	return strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")
}

// extractINISection extracts section name from section header
func extractINISection(line string) string {
	return line[1 : len(line)-1]
}

// parseINIKeyValue parses a key-value pair from a line
func parseINIKeyValue(line string) (string, string, bool) {
	separatorIndex := findKeyValueSeparator(line)
	if separatorIndex == -1 {
		return "", "", false
	}

	key, value := extractKeyValue(line, separatorIndex)
	return key, value, isValidKey(key)
}

// setINIValue sets a value in the appropriate section
func setINIValue(result map[string]interface{}, section, key, value string) {
	cleanValue := removeQuotes(value)
	
	if section != "" {
		if sectionMap, ok := result[section].(map[string]interface{}); ok {
			sectionMap[key] = cleanValue
		}
	} else {
		result[key] = cleanValue
	}
}
