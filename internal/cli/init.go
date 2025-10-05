package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewInitCommand creates the init command
func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Praetorian configuration for your project",
		Long: `Initialize Praetorian configuration for your project.

This command creates a praetorian.yaml configuration file with sensible defaults
for validating your configuration files across environments.

Examples:
  praetorian init                    # Create basic configuration
  praetorian init --devsecops        # Create DevSecOps optimized configuration`,
		RunE: runInit,
	}

	// Add flags
	cmd.Flags().Bool("devsecops", false, "Initialize with DevSecOps optimizations")

	return cmd
}

// runInit executes the init command
func runInit(cmd *cobra.Command, args []string) error {
	// Guard clause: validate command
	if cmd == nil {
		return fmt.Errorf("command cannot be nil")
	}

	// Extract and validate flags
	flags, err := extractInitFlags(cmd)
	if err != nil {
		return fmt.Errorf("failed to extract flags: %w", err)
	}

	// Execute initialization
	return executeInitialization(flags)
}

// InitFlags represents init command flags
type InitFlags struct {
	DevSecOps bool
}

// extractInitFlags extracts and validates flags from command
func extractInitFlags(cmd *cobra.Command) (*InitFlags, error) {
	devsecops, err := cmd.Flags().GetBool("devsecops")
	if err != nil {
		return nil, fmt.Errorf("failed to get devsecops flag: %w", err)
	}

	return &InitFlags{
		DevSecOps: devsecops,
	}, nil
}

// executeInitialization executes the initialization process
func executeInitialization(flags *InitFlags) error {
	// Guard clause: validate flags
	if flags == nil {
		return fmt.Errorf("flags cannot be nil")
	}

	// Display initialization info
	displayInitializationInfo(flags)

	// Create configuration file
	if err := createConfigFile(flags); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	fmt.Printf("✅ Configuration initialized successfully!\n")
	return nil
}

// displayInitializationInfo displays initialization information
func displayInitializationInfo(flags *InitFlags) {
	fmt.Printf("🚀 Initializing Praetorian configuration...\n")
	fmt.Printf("🛡️  DevSecOps mode: %t\n", flags.DevSecOps)
	fmt.Printf("🔧 Creating config file...\n")
}

// createConfigFile creates the configuration file
func createConfigFile(flags *InitFlags) error {
	// Guard clause: validate flags
	if flags == nil {
		return fmt.Errorf("flags cannot be nil")
	}

	// Generate config content
	content := generateConfigContent(flags.DevSecOps)

	// Write file
	return writeConfigFileContent("praetorian.yaml", content)
}

// generateConfigContent generates the configuration content
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

// writeConfigFileContent writes content to a file
func writeConfigFileContent(filename, content string) error {
	// Guard clause: validate filename
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	// Guard clause: validate content
	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	// Write content
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write content to %s: %w", filename, err)
	}

	return nil
}
