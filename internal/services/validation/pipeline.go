package validation

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// FilePipeline implements the Pipeline interface
type FilePipeline struct {
	readers    []models.FileReader
	processors []models.FileProcessor
	config     models.PipelineConfig
}

// NewFilePipeline creates a new file pipeline
func NewFilePipeline(config models.PipelineConfig) *FilePipeline {
	return &FilePipeline{
		readers:    make([]models.FileReader, 0),
		processors: make([]models.FileProcessor, 0),
		config:     config,
	}
}

// ProcessFiles processes multiple files concurrently
func (p *FilePipeline) ProcessFiles(ctx context.Context, filenames []string) ([]*models.ConfigData, error) {
	// Guard clause: validate input
	if err := p.validateFilenames(filenames); err != nil {
		return nil, fmt.Errorf("invalid filenames: %w", err)
	}

	// Guard clause: check if we have processors
	if len(p.processors) == 0 {
		return nil, fmt.Errorf("no processors registered")
	}

	// Process files concurrently
	return p.processFilesConcurrently(ctx, filenames)
}

// ProcessFile processes a single file
func (p *FilePipeline) ProcessFile(ctx context.Context, filename string) (*models.ConfigData, error) {
	// Guard clause: validate input
	if err := p.validateFilename(filename); err != nil {
		return nil, fmt.Errorf("invalid filename: %w", err)
	}

	// Find appropriate processor
	processor := p.findProcessor(filename)
	if processor == nil {
		return nil, fmt.Errorf("no processor found for file: %s", filename)
	}

	// Read file content
	content, err := p.readFileContent(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Process file
	return processor.Process(ctx, filename, content)
}

// RegisterProcessor registers a new processor
func (p *FilePipeline) RegisterProcessor(processor models.FileProcessor) error {
	// Guard clause: validate processor
	if processor == nil {
		return fmt.Errorf("processor cannot be nil")
	}

	// Guard clause: check for duplicates
	if p.isProcessorRegistered(processor) {
		return fmt.Errorf("processor already registered")
	}

	p.processors = append(p.processors, processor)
	return nil
}

// GetProcessors returns all registered processors
func (p *FilePipeline) GetProcessors() []models.FileProcessor {
	// Return a copy to prevent external modification
	processors := make([]models.FileProcessor, len(p.processors))
	copy(processors, p.processors)
	return processors
}

// processFilesConcurrently processes files using worker pool pattern
func (p *FilePipeline) processFilesConcurrently(ctx context.Context, filenames []string) ([]*models.ConfigData, error) {
	// Create channels for work distribution
	jobs := make(chan string, len(filenames))
	results := make(chan *fileProcessingResult, len(filenames))
	errors := make(chan error, len(filenames))

	// Start workers
	var wg sync.WaitGroup
	workerCount := p.getWorkerCount()

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go p.worker(ctx, &wg, jobs, results, errors)
	}

	// Send jobs
	for _, filename := range filenames {
		jobs <- filename
	}
	close(jobs)

	// Wait for workers to complete
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results
	return p.collectResults(results, errors)
}

// worker processes files from the jobs channel
func (p *FilePipeline) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- *fileProcessingResult, errors chan<- error) {
	defer wg.Done()

	for filename := range jobs {
		select {
		case <-ctx.Done():
			errors <- ctx.Err()
			return
		default:
			result := p.processFileWithResult(ctx, filename)
			if result.Error != nil {
				errors <- result.Error
			} else {
				results <- result
			}
		}
	}
}

// processFileWithResult processes a file and returns a result
func (p *FilePipeline) processFileWithResult(ctx context.Context, filename string) *fileProcessingResult {
	configData, err := p.ProcessFile(ctx, filename)
	return &fileProcessingResult{
		Filename: filename,
		Data:     configData,
		Error:    err,
	}
}

// collectResults collects results from workers
func (p *FilePipeline) collectResults(results <-chan *fileProcessingResult, errors <-chan error) ([]*models.ConfigData, error) {
	var processed []*models.ConfigData
	var processingErrors []error

	// Collect successful results
	for result := range results {
		if result.Data != nil {
			processed = append(processed, result.Data)
		}
	}

	// Collect errors
	for err := range errors {
		processingErrors = append(processingErrors, err)
	}

	// Return error if any processing failed
	if len(processingErrors) > 0 {
		return processed, fmt.Errorf("processing errors: %v", processingErrors)
	}

	return processed, nil
}

// findProcessor finds the appropriate processor for a filename
func (p *FilePipeline) findProcessor(filename string) models.FileProcessor {
	// Guard clause: no processors registered
	if len(p.processors) == 0 {
		return nil
	}

	// Find processor that can handle this file
	for _, processor := range p.processors {
		if processor.CanProcess(filename) {
			return processor
		}
	}

	return nil
}

// readFileContent reads file content using registered readers
func (p *FilePipeline) readFileContent(filename string) ([]byte, error) {
	// Guard clause: no readers registered
	if len(p.readers) == 0 {
		return nil, fmt.Errorf("no file readers registered")
	}

	// Try each reader until one succeeds
	for _, reader := range p.readers {
		if reader.FileExists(filename) {
			return reader.ReadFile(filename)
		}
	}

	return nil, fmt.Errorf("file not found: %s", filename)
}

// validateFilenames validates a list of filenames
func (p *FilePipeline) validateFilenames(filenames []string) error {
	if filenames == nil {
		return fmt.Errorf("filenames cannot be nil")
	}
	if len(filenames) == 0 {
		return fmt.Errorf("filenames cannot be empty")
	}
	
	for _, filename := range filenames {
		if err := p.validateFilename(filename); err != nil {
			return fmt.Errorf("invalid filename %s: %w", filename, err)
		}
	}
	
	return nil
}

// validateFilename validates a single filename
func (p *FilePipeline) validateFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if !filepath.IsAbs(filename) && !strings.HasPrefix(filename, "./") {
		return fmt.Errorf("filename must be absolute path or start with ./")
	}
	return nil
}

// isProcessorRegistered checks if a processor is already registered
func (p *FilePipeline) isProcessorRegistered(newProcessor models.FileProcessor) bool {
	newExtensions := newProcessor.GetSupportedExtensions()
	
	for _, existing := range p.processors {
		existingExtensions := existing.GetSupportedExtensions()
		if p.hasExtensionOverlap(newExtensions, existingExtensions) {
			return true
		}
	}
	
	return false
}

// hasExtensionOverlap checks if two extension lists overlap
func (p *FilePipeline) hasExtensionOverlap(ext1, ext2 []string) bool {
	for _, e1 := range ext1 {
		for _, e2 := range ext2 {
			if e1 == e2 {
				return true
			}
		}
	}
	return false
}

// getWorkerCount returns the number of workers to use
func (p *FilePipeline) getWorkerCount() int {
	if p.config.MaxWorkers > 0 {
		return p.config.MaxWorkers
	}
	return 4 // Default worker count
}

// fileProcessingResult represents the result of processing a single file
type fileProcessingResult struct {
	Filename string
	Data     *models.ConfigData
	Error    error
}
