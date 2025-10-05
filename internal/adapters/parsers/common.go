package parsers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// ValidateFilenameAndExtension validates filename and checks extension
func ValidateFilenameAndExtension(filename string, supportedExtensions []string) bool {
	// Guard clause: empty filename
	if filename == "" {
		return false
	}

	ext := GetFileExtension(filename)
	return supportsExtension(ext, supportedExtensions)
}

// ValidateContextAndInput validates context and input parameters
func ValidateContextAndInput(ctx context.Context, filename string, content []byte) error {
	// Guard clause: context cancelled
	if ctx.Err() != nil {
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	}

	// Guard clause: empty filename
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	// Guard clause: nil content
	if content == nil {
		return fmt.Errorf("content cannot be nil")
	}

	return nil
}

// GetFileExtension extracts extension from filename
func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.TrimPrefix(strings.ToLower(ext), ".")
}

// supportsExtension checks if extension is supported
func supportsExtension(extension string, supportedExtensions []string) bool {
	for _, ext := range supportedExtensions {
		if ext == extension {
			return true
		}
	}
	return false
}

// copyExtensions returns a copy of extensions slice
func copyExtensions(extensions []string) []string {
	result := make([]string, len(extensions))
	copy(result, extensions)
	return result
}

// createConfigData creates ConfigData from parsed content
func createConfigData(filename, format string, data map[string]interface{}) *models.ConfigData {
	return &models.ConfigData{
		Filename:  filename,
		Format:    format,
		Data:      data,
		Metadata:  createMetadata(filename, format, data),
		Timestamp: time.Now(),
	}
}

// createMetadata creates metadata for the config data
func createMetadata(filename, format string, data map[string]interface{}) map[string]interface{} {
	metadata := map[string]interface{}{
		"processor":  format,
		"file_size":  len(data),
		"key_count":  countKeys(data),
		"has_nested": hasNestedData(data),
	}

	// Add file-specific metadata only if filename is provided
	if filename != "" {
		_, _, _ = metadata, filepath.Base, filename
		metadata["dirname"] = filepath.Dir(filename)
	}

	return metadata
}

// countKeys counts the number of top-level keys
func countKeys(data map[string]interface{}) int {
	if data == nil {
		return 0
	}
	return len(data)
}

// hasNestedData checks if the data contains nested structures
func hasNestedData(data map[string]interface{}) bool {
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

// removeQuotes removes quotes from a string value
func removeQuotes(value string) string {
	// Guard clause: empty value
	if value == "" {
		return value
	}

	// Guard clause: value too short for quotes
	if len(value) < 2 {
		return value
	}

	// Check for double quotes
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return value[1 : len(value)-1]
	}

	// Check for single quotes
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		return value[1 : len(value)-1]
	}

	return value
}

// isEmptyContent checks if content is empty
func isEmptyContent(content []byte) bool {
	return len(content) == 0
}

// createEmptyResult creates an empty result map
func createEmptyResult() map[string]interface{} {
	return make(map[string]interface{})
}

// splitLines splits content into lines
func splitLines(content []byte) []string {
	return strings.Split(string(content), "\n")
}

// trimLine trims whitespace from a line
func trimLine(line string) string {
	return strings.TrimSpace(line)
}

// isEmptyLine checks if a line is empty or a comment
func isEmptyLine(line string) bool {
	return line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "!")
}

// findKeyValueSeparator finds the separator in a key-value line
func findKeyValueSeparator(line string) int {
	return strings.Index(line, "=")
}

// extractKeyValue extracts key and value from a line
func extractKeyValue(line string, separatorIndex int) (string, string) {
	key := strings.TrimSpace(line[:separatorIndex])
	value := strings.TrimSpace(line[separatorIndex+1:])
	return key, value
}

// isValidKey checks if a key is valid
func isValidKey(key string) bool {
	return key != ""
}

// parseKeyValueContent parses key-value content with comment detection
func parseKeyValueContent(content []byte, isComment func(string) bool) (map[string]interface{}, error) {
	// Guard clause: empty content
	if isEmptyContent(content) {
		return createEmptyResult(), nil
	}

	result := createEmptyResult()
	lines := splitLines(content)

	for _, line := range lines {
		trimmedLine := trimLine(line)

		// Skip empty lines and comments
		if isEmptyLine(trimmedLine) || isComment(trimmedLine) {
			continue
		}

		// Parse key-value pair
		if key, value, ok := parseKeyValue(trimmedLine); ok {
			result[key] = removeQuotes(value)
		}
	}

	return result, nil
}

// parseKeyValue parses a key-value pair from a line
func parseKeyValue(line string) (string, string, bool) {
	separatorIndex := findKeyValueSeparator(line)
	if separatorIndex == -1 {
		return "", "", false
	}

	key, value := extractKeyValue(line, separatorIndex)
	return key, value, isValidKey(key)
}
