package parsers

import (
	"context"
	"testing"
)

// TestValidateFilenameAndExtension tests filename validation
func TestValidateFilenameAndExtension(t *testing.T) {
	tests := []struct {
		name               string
		filename           string
		supportedExts      []string
		expected           bool
	}{
		{"valid yaml", "test.yaml", []string{"yaml"}, true},
		{"valid json", "test.json", []string{"json"}, true},
		{"invalid extension", "test.txt", []string{"yaml"}, false},
		{"empty filename", "", []string{"yaml"}, false},
		{"multiple extensions", "test.yaml", []string{"yaml", "yml"}, true},
		{"case insensitive", "test.YAML", []string{"yaml"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateFilenameAndExtension(tt.filename, tt.supportedExts)
			if result != tt.expected {
				t.Errorf("validateFilenameAndExtension() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestValidateContextAndInput tests context validation
func TestValidateContextAndInput(t *testing.T) {
	ctx := context.Background()
	
	tests := []struct {
		name      string
		ctx       context.Context
		filename  string
		content   []byte
		wantErr   bool
	}{
		{"valid input", ctx, "test.yaml", []byte("test"), false},
		{"empty filename", ctx, "", []byte("test"), true},
		{"nil content", ctx, "test.yaml", nil, true},
		{"cancelled context", context.Background(), "test.yaml", []byte("test"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContextAndInput(tt.ctx, tt.filename, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateContextAndInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetFileExtension tests file extension extraction
func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{"yaml file", "test.yaml", "yaml"},
		{"json file", "test.json", "json"},
		{"no extension", "test", ""},
		{"multiple dots", "test.config.yaml", "yaml"},
		{"uppercase", "test.YAML", "yaml"},
		{"empty filename", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFileExtension(tt.filename)
			if result != tt.expected {
				t.Errorf("getFileExtension() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// BenchmarkValidateFilenameAndExtension benchmarks filename validation
func BenchmarkValidateFilenameAndExtension(b *testing.B) {
	exts := []string{"yaml", "json", "toml"}
	for i := 0; i < b.N; i++ {
		ValidateFilenameAndExtension("test.yaml", exts)
	}
}
