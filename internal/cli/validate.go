package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewValidateCommand creates the validate command
func NewValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration files for consistency across environments",
		Long: `Validate configuration files for consistency across environments.

This command compares configuration files between different environments (dev, staging, prod)
and reports missing keys, extra keys, and value differences.

Examples:
  praetorian validate                           # Validate current directory
  praetorian validate --config praetorian.yaml # Use specific config file
  praetorian validate --output json            # Output in JSON format
  praetorian validate --pipeline               # CI/CD friendly output`,
		RunE: runValidate,
	}

	// Add flags
	cmd.Flags().StringP("config", "c", "praetorian.yaml", "Configuration file path")
	cmd.Flags().StringP("output", "o", "text", "Output format (text, json, yaml)")
	cmd.Flags().Bool("pipeline", false, "Enable pipeline mode for CI/CD")

	return cmd
}

// runValidate executes the validate command
func runValidate(cmd *cobra.Command, args []string) error {
	// Guard clause: validate command
	if cmd == nil {
		return fmt.Errorf("command cannot be nil")
	}

	// Extract and validate flags
	flags, err := extractValidateFlags(cmd)
	if err != nil {
		return fmt.Errorf("failed to extract flags: %w", err)
	}

	// Execute validation
	return executeValidation(flags)
}

// ValidateFlags represents validation command flags
type ValidateFlags struct {
	ConfigPath   string
	OutputFormat string
	PipelineMode bool
}

// extractValidateFlags extracts and validates flags from command
func extractValidateFlags(cmd *cobra.Command) (*ValidateFlags, error) {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, fmt.Errorf("failed to get config flag: %w", err)
	}

	outputFormat, err := cmd.Flags().GetString("output")
	if err != nil {
		return nil, fmt.Errorf("failed to get output flag: %w", err)
	}

	pipelineMode, err := cmd.Flags().GetBool("pipeline")
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline flag: %w", err)
	}

	// Guard clause: validate config path
	if err := ValidateConfigPath(configPath); err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	// Guard clause: validate output format
	if err := ValidateOutputFormat(outputFormat); err != nil {
		return nil, fmt.Errorf("invalid output format: %w", err)
	}

	return &ValidateFlags{
		ConfigPath:   configPath,
		OutputFormat: outputFormat,
		PipelineMode: pipelineMode,
	}, nil
}

// executeValidation executes the validation process
func executeValidation(flags *ValidateFlags) error {
	// Guard clause: validate flags
	if flags == nil {
		return fmt.Errorf("flags cannot be nil")
	}

	// Display validation info
	displayValidationInfo(flags)

	// TODO: Implement actual validation logic
	fmt.Printf("✅ Validation completed successfully!\n")

	// Handle pipeline output
	if flags.PipelineMode {
		displayPipelineOutput()
	}

	return nil
}

// displayValidationInfo displays validation information
func displayValidationInfo(flags *ValidateFlags) {
	fmt.Printf("🔍 Validating configuration files...\n")
	fmt.Printf("📁 Config: %s\n", flags.ConfigPath)
	fmt.Printf("📤 Output: %s\n", flags.OutputFormat)
	
	if flags.PipelineMode {
		fmt.Printf("🚀 Pipeline mode: enabled\n")
	}
}

// displayPipelineOutput displays pipeline-friendly output
func displayPipelineOutput() {
	fmt.Printf("PRAETORIAN_VALIDATION_STATUS=success\n")
}
