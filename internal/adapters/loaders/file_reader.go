package loaders

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/syntropysoft/praetorian-go/internal/domain/models"
)

// LocalFileReader implements FileReader interface for local filesystem
type LocalFileReader struct {
	basePath string
}

// NewLocalFileReader creates a new local file reader
func NewLocalFileReader(basePath string) *LocalFileReader {
	// Guard clause: validate and normalize base path
	normalizedPath := normalizeBasePath(basePath)
	
	return &LocalFileReader{
		basePath: normalizedPath,
	}
}

// normalizeBasePath normalizes and validates a base path
func normalizeBasePath(basePath string) string {
	// Guard clause: empty path defaults to current directory
	if basePath == "" {
		return getCurrentDirectory()
	}
	
	// Guard clause: already absolute path
	if filepath.IsAbs(basePath) {
		return basePath
	}
	
	// Convert relative path to absolute
	return convertToAbsolutePath(basePath)
}

// getCurrentDirectory returns the current working directory
func getCurrentDirectory() string {
	dir, err := filepath.Abs(".")
	if err != nil {
		return "." // Fallback
	}
	return dir
}

// convertToAbsolutePath converts a relative path to absolute
func convertToAbsolutePath(relativePath string) string {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return getCurrentDirectory() // Fallback to current directory
	}
	return absPath
}

// ReadFile reads a file from the local filesystem
func (r *LocalFileReader) ReadFile(filename string) ([]byte, error) {
	// Guard clause: validate filename
	if err := r.validateFilename(filename); err != nil {
		return nil, fmt.Errorf("invalid filename: %w", err)
	}

	// Guard clause: check if file exists
	if !r.FileExists(filename) {
		return nil, fmt.Errorf("file does not exist: %s", filename)
	}

	// Read file content
	fullPath := r.getFullPath(filename)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return content, nil
}

// FileExists checks if a file exists
func (r *LocalFileReader) FileExists(filename string) bool {
	// Guard clause: validate filename
	if err := r.validateFilename(filename); err != nil {
		return false
	}

	fullPath := r.getFullPath(filename)
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}

// GetFileInfo returns file information
func (r *LocalFileReader) GetFileInfo(filename string) (*models.FileInfo, error) {
	// Guard clause: validate filename
	if err := r.validateFilename(filename); err != nil {
		return nil, fmt.Errorf("invalid filename: %w", err)
	}

	// Guard clause: check if file exists
	if !r.FileExists(filename) {
		return nil, fmt.Errorf("file does not exist: %s", filename)
	}

	fullPath := r.getFullPath(filename)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for %s: %w", filename, err)
	}

	return &models.FileInfo{
		Name:    stat.Name(),
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
		IsDir:   stat.IsDir(),
	}, nil
}

// validateFilename validates a filename
func (r *LocalFileReader) validateFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	
	// Check for path traversal attacks
	if strings.Contains(filename, "..") {
		return fmt.Errorf("filename contains invalid path traversal: %s", filename)
	}
	
	// Check for absolute paths outside base path
	if filepath.IsAbs(filename) {
		// For absolute paths, ensure they're within base path
		if !strings.HasPrefix(filename, r.basePath) {
			return fmt.Errorf("filename outside base path: %s", filename)
		}
	}
	
	return nil
}

// getFullPath returns the full path for a filename
func (r *LocalFileReader) getFullPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(r.basePath, filename)
}

// GetBasePath returns the base path
func (r *LocalFileReader) GetBasePath() string {
	return r.basePath
}

// SetBasePath sets a new base path
func (r *LocalFileReader) SetBasePath(basePath string) error {
	// Guard clause: validate base path
	if err := validateBasePath(basePath); err != nil {
		return fmt.Errorf("invalid base path: %w", err)
	}

	// Guard clause: check if path exists and is directory
	if err := validateDirectoryExists(basePath); err != nil {
		return fmt.Errorf("directory validation failed: %w", err)
	}

	// Set normalized absolute path
	r.basePath = normalizeBasePath(basePath)
	return nil
}

// validateBasePath validates a base path
func validateBasePath(basePath string) error {
	if basePath == "" {
		return fmt.Errorf("base path cannot be empty")
	}
	return nil
}

// validateDirectoryExists validates that a path exists and is a directory
func validateDirectoryExists(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path does not exist: %w", err)
	}
	
	if !stat.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}
	
	return nil
}

// ListFiles lists files matching a pattern
func (r *LocalFileReader) ListFiles(pattern string) ([]string, error) {
	// Guard clause: validate pattern
	if err := validatePattern(pattern); err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	// Find matching files
	matches, err := findMatchingFiles(r.getFullPath(pattern))
	if err != nil {
		return nil, fmt.Errorf("failed to find files: %w", err)
	}

	// Convert to relative paths
	return convertToRelativePaths(matches, r.basePath), nil
}

// validatePattern validates a glob pattern
func validatePattern(pattern string) error {
	if pattern == "" {
		return fmt.Errorf("pattern cannot be empty")
	}
	return nil
}

// findMatchingFiles finds files matching a pattern
func findMatchingFiles(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob failed for pattern %s: %w", pattern, err)
	}
	return matches, nil
}

// convertToRelativePaths converts absolute paths to relative paths
func convertToRelativePaths(absolutePaths []string, basePath string) []string {
	relativePaths := make([]string, 0, len(absolutePaths))
	
	for _, absPath := range absolutePaths {
		if relativePath, err := filepath.Rel(basePath, absPath); err == nil {
			relativePaths = append(relativePaths, relativePath)
		}
		// Skip files that can't be made relative (silently)
	}
	
	return relativePaths
}

// IsDirectory checks if a path is a directory
func (r *LocalFileReader) IsDirectory(path string) bool {
	// Guard clause: validate path
	if err := r.validateFilename(path); err != nil {
		return false
	}

	fullPath := r.getFullPath(path)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	
	return stat.IsDir()
}

// GetFileSize returns the size of a file
func (r *LocalFileReader) GetFileSize(filename string) (int64, error) {
	// Guard clause: validate filename
	if err := r.validateFilename(filename); err != nil {
		return 0, fmt.Errorf("invalid filename: %w", err)
	}

	// Guard clause: check if file exists
	if !r.FileExists(filename) {
		return 0, fmt.Errorf("file does not exist: %s", filename)
	}

	fullPath := r.getFullPath(filename)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file size for %s: %w", filename, err)
	}

	return stat.Size(), nil
}

// GetFileModTime returns the modification time of a file
func (r *LocalFileReader) GetFileModTime(filename string) (time.Time, error) {
	// Guard clause: validate filename
	if err := r.validateFilename(filename); err != nil {
		return time.Time{}, fmt.Errorf("invalid filename: %w", err)
	}

	// Guard clause: check if file exists
	if !r.FileExists(filename) {
		return time.Time{}, fmt.Errorf("file does not exist: %s", filename)
	}

	fullPath := r.getFullPath(filename)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get file mod time for %s: %w", filename, err)
	}

	return stat.ModTime(), nil
}
