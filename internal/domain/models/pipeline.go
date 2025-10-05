package models

import (
	"context"
	"time"
)

// FileReader defines the interface for reading files from filesystem
type FileReader interface {
	ReadFile(filename string) ([]byte, error)
	FileExists(filename string) bool
	GetFileInfo(filename string) (*FileInfo, error)
}

// FileInfo represents basic file information
type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	IsDir   bool      `json:"is_dir"`
}

// FileProcessor defines the interface for processing files in pipeline
type FileProcessor interface {
	CanProcess(filename string) bool
	Process(ctx context.Context, filename string, content []byte) (*ConfigData, error)
	GetSupportedExtensions() []string
}

// Pipeline defines the interface for file processing pipeline
type Pipeline interface {
	ProcessFiles(ctx context.Context, filenames []string) ([]*ConfigData, error)
	ProcessFile(ctx context.Context, filename string) (*ConfigData, error)
	RegisterProcessor(processor FileProcessor) error
	GetProcessors() []FileProcessor
}

// PipelineResult represents the result of pipeline processing
type PipelineResult struct {
	Success      bool                   `json:"success"`
	Processed    []*ConfigData          `json:"processed"`
	Failed       []*ProcessingError     `json:"failed,omitempty"`
	Summary      PipelineSummary        `json:"summary"`
	Duration     time.Duration          `json:"duration"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessingError represents an error during file processing
type ProcessingError struct {
	Filename string `json:"filename"`
	Error    string `json:"error"`
	Code     string `json:"code"`
}

// PipelineSummary represents pipeline processing summary
type PipelineSummary struct {
	TotalFiles    int `json:"total_files"`
	Processed     int `json:"processed"`
	Failed        int `json:"failed"`
	ByFormat      map[string]int `json:"by_format"`
}

// PipelineConfig represents pipeline configuration
type PipelineConfig struct {
	MaxWorkers    int           `json:"max_workers"`
	Timeout       time.Duration `json:"timeout"`
	RetryAttempts int           `json:"retry_attempts"`
	BufferSize    int           `json:"buffer_size"`
}
