package parsers

import (
	"fmt"

	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// RegisterAllProcessors registers all available processors with the registry
func RegisterAllProcessors(registry *ParserRegistry) error {
	// Guard clause: validate registry
	if registry == nil {
		return fmt.Errorf("registry cannot be nil")
	}

	// Register all processors
	processors := []models.FileProcessor{
		NewYAMLProcessor(),
		NewJSONProcessor(),
		NewTOMLProcessor(),
		NewPropertiesProcessor(),
		NewINIProcessor(),
		NewHCLProcessor(),
		NewXMLProcessor(),
		NewENVProcessor(),
	}

	// Register each processor
	for _, processor := range processors {
		if err := registry.RegisterProcessor(processor); err != nil {
			return fmt.Errorf("failed to register processor: %w", err)
		}
	}

	return nil
}
