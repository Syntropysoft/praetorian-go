package cli

import (
	"testing"

	"github.com/syntropysoft/praetorian-go/internal/cli"
)

// TestValidateConfigPath tests config path validation
func TestValidateConfigPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid path", "config.yaml", false},
		{"valid path with dir", "config/test.yaml", false},
		{"empty path", "", true},
		{"valid absolute path", "/tmp/config.yaml", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cli.ValidateConfigPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfigPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateOutputFormat tests output format validation
func TestValidateOutputFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"valid text", "text", false},
		{"valid json", "json", false},
		{"valid yaml", "yaml", false},
		{"invalid format", "xml", true},
		{"empty format", "", true},
		{"case sensitive", "JSON", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cli.ValidateOutputFormat(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOutputFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// BenchmarkValidateConfigPath benchmarks config path validation
func BenchmarkValidateConfigPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cli.ValidateConfigPath("test-config.yaml")
	}
}

// BenchmarkValidateOutputFormat benchmarks output format validation
func BenchmarkValidateOutputFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cli.ValidateOutputFormat("json")
	}
}
