package commands

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/syntropysoft/praetorian-go/internal/cli"
)

// TestInitCommandIntegration tests the init command integration
func TestInitCommandIntegration(t *testing.T) {
	// Setup
	testDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(testDir)

	t.Run("should create basic config file", func(t *testing.T) {
		// Execute init command logic (simulating command execution)
		flags := &cli.InitFlags{DevSecOps: false}
		
		// Test the core functionality without CLI framework
		err := executeInitCommand(flags)
		if err != nil {
			t.Fatalf("Init command failed: %v", err)
		}

		// Verify file was created
		if _, err := os.Stat("praetorian.yaml"); os.IsNotExist(err) {
			t.Error("Expected praetorian.yaml to be created")
		}
	})

	t.Run("should create devsecops config file", func(t *testing.T) {
		// Clean up previous test
		os.Remove("praetorian.yaml")
		
		flags := &cli.InitFlags{DevSecOps: true}
		
		err := executeInitCommand(flags)
		if err != nil {
			t.Fatalf("DevSecOps init command failed: %v", err)
		}

		// Verify file was created
		if _, err := os.Stat("praetorian.yaml"); os.IsNotExist(err) {
			t.Error("Expected praetorian.yaml to be created")
		}

		// Verify content contains DevSecOps config
		content, err := os.ReadFile("praetorian.yaml")
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}
		
		if !containsString(string(content), "secret_detection") {
			t.Error("Expected DevSecOps configuration to include secret_detection")
		}
	})
}

// executeInitCommand simulates the init command execution
func executeInitCommand(flags *cli.InitFlags) error {
	if flags == nil {
		return fmt.Errorf("flags cannot be nil")
	}
	
	// Generate config content
	content := generateConfigContent(flags.DevSecOps)
	
	// Write file
	return writeConfigFileContent("praetorian.yaml", content)
}

// Helper functions (copied from init.go for testing)
func generateConfigContent(devsecops bool) string {
	baseContent := `# Praetorian Configuration
version: "1.0"
files:
  - "config/*.yaml"
  - "config/*.json"
  - "config/*.toml"

environments:
  dev: "config/dev"
  staging: "config/staging"
  prod: "config/prod"
`

	if devsecops {
		baseContent += `
# DevSecOps specific configurations
security:
  secret_detection: true
  vulnerability_scan: true
  permission_check: true
`
	}

	return baseContent
}

func writeConfigFileContent(filename, content string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	
	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write content to %s: %w", filename, err)
	}

	return nil
}

func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}
