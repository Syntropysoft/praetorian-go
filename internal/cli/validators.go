package cli

import "fmt"

// ValidateConfigPath validates the config path
func ValidateConfigPath(path string) error {
	if path == "" {
		return fmt.Errorf("config path cannot be empty")
	}
	return nil
}

// ValidateOutputFormat validates the output format
func ValidateOutputFormat(format string) error {
	validFormats := []string{"text", "json", "yaml"}
	for _, valid := range validFormats {
		if format == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid output format: %s, must be one of: %v", format, validFormats)
}
