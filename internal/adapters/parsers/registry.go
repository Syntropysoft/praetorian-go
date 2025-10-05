package parsers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// ParserRegistry manages file processors
type ParserRegistry struct {
	processors []models.FileProcessor
	byExtension map[string]models.FileProcessor
}

// NewParserRegistry creates a new parser registry with all processors registered
func NewParserRegistry() *ParserRegistry {
	registry := &ParserRegistry{
		processors:  make([]models.FileProcessor, 0),
		byExtension: make(map[string]models.FileProcessor),
	}

	// Register all processors
	if err := RegisterAllProcessors(registry); err != nil {
		panic(fmt.Sprintf("Failed to register processors: %v", err))
	}

	return registry
}

// RegisterProcessor registers a new processor
func (r *ParserRegistry) RegisterProcessor(processor models.FileProcessor) error {
	// Guard clause: validate processor
	if processor == nil {
		return fmt.Errorf("processor cannot be nil")
	}

	// Guard clause: check for conflicts
	if err := r.validateNoConflicts(processor); err != nil {
		return fmt.Errorf("processor conflicts with existing: %w", err)
	}

	// Register processor
	r.processors = append(r.processors, processor)

	// Register by extension for fast lookup
	for _, ext := range processor.GetSupportedExtensions() {
		r.byExtension[ext] = processor
	}

	return nil
}

// GetProcessor returns the appropriate processor for a filename
func (r *ParserRegistry) GetProcessor(filename string) (models.FileProcessor, error) {
	// Guard clause: validate filename
	if filename == "" {
		return nil, fmt.Errorf("filename cannot be empty")
	}

	// Get file extension
	ext := r.getFileExtension(filename)
	if ext == "" {
		return nil, fmt.Errorf("no extension found in filename: %s", filename)
	}

	// Look up processor by extension
	processor, exists := r.byExtension[ext]
	if !exists {
		return nil, fmt.Errorf("no processor found for extension: %s", ext)
	}

	return processor, nil
}

// GetProcessorByExtension returns processor by extension
func (r *ParserRegistry) GetProcessorByExtension(extension string) (models.FileProcessor, error) {
	// Guard clause: validate extension
	if extension == "" {
		return nil, fmt.Errorf("extension cannot be empty")
	}

	// Normalize extension (remove leading dot if present)
	ext := r.normalizeExtension(extension)

	processor, exists := r.byExtension[ext]
	if !exists {
		return nil, fmt.Errorf("no processor found for extension: %s", ext)
	}

	return processor, nil
}

// GetAllProcessors returns all registered processors
func (r *ParserRegistry) GetAllProcessors() []models.FileProcessor {
	// Return a copy to prevent external modification
	processors := make([]models.FileProcessor, len(r.processors))
	copy(processors, r.processors)
	return processors
}

// GetSupportedExtensions returns all supported extensions
func (r *ParserRegistry) GetSupportedExtensions() []string {
	extensions := make([]string, 0, len(r.byExtension))
	for ext := range r.byExtension {
		extensions = append(extensions, ext)
	}
	return extensions
}

// HasProcessor checks if a processor is registered for an extension
func (r *ParserRegistry) HasProcessor(extension string) bool {
	ext := r.normalizeExtension(extension)
	_, exists := r.byExtension[ext]
	return exists
}

// GetProcessorCount returns the number of registered processors
func (r *ParserRegistry) GetProcessorCount() int {
	return len(r.processors)
}

// validateNoConflicts checks if a processor conflicts with existing ones
func (r *ParserRegistry) validateNoConflicts(newProcessor models.FileProcessor) error {
	newExtensions := newProcessor.GetSupportedExtensions()
	
	for _, existing := range r.processors {
		existingExtensions := existing.GetSupportedExtensions()
		if r.hasExtensionOverlap(newExtensions, existingExtensions) {
			return fmt.Errorf("extension conflict between processors")
		}
	}
	
	return nil
}

// hasExtensionOverlap checks if two extension lists overlap
func (r *ParserRegistry) hasExtensionOverlap(ext1, ext2 []string) bool {
	for _, e1 := range ext1 {
		for _, e2 := range ext2 {
			if e1 == e2 {
				return true
			}
		}
	}
	return false
}

// getFileExtension extracts extension from filename
func (r *ParserRegistry) getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	return r.normalizeExtension(ext)
}

// normalizeExtension normalizes an extension (removes leading dot, converts to lowercase)
func (r *ParserRegistry) normalizeExtension(extension string) string {
	ext := strings.TrimPrefix(extension, ".")
	return strings.ToLower(ext)
}

// CreateDefaultRegistry creates a registry with default processors
func CreateDefaultRegistry() (*ParserRegistry, error) {
	registry := NewParserRegistry()

	// Register default processors
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

	for _, processor := range processors {
		if err := registry.RegisterProcessor(processor); err != nil {
			return nil, fmt.Errorf("failed to register processor: %w", err)
		}
	}

	return registry, nil
}
