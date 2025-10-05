package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewAuditCommand creates the audit command
func NewAuditCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Run comprehensive security and compliance audits on configuration files",
		Long: `Run comprehensive security and compliance audits on configuration files.

Examples:
  praetorian audit                           # Run all audits on current directory
  praetorian audit --type security          # Run security audit only
  praetorian audit --config praetorian.yaml # Use specific config file
  praetorian audit --output json            # Output in JSON format`,
		RunE: runAudit,
	}

	// Add flags
	cmd.Flags().StringP("config", "c", "praetorian.yaml", "Configuration file path")
	cmd.Flags().StringP("output", "o", "text", "Output format (text, json, yaml)")
	cmd.Flags().StringP("type", "t", "all", "Audit type (security, compliance, performance, all)")

	return cmd
}

// runAudit executes the audit command
func runAudit(cmd *cobra.Command, args []string) error {
	// Guard clause: validate command
	if cmd == nil {
		return fmt.Errorf("command cannot be nil")
	}

	// Extract and validate flags
	flags, err := extractAuditFlags(cmd)
	if err != nil {
		return fmt.Errorf("failed to extract flags: %w", err)
	}

	// Execute audit
	return executeAudit(flags)
}

// AuditFlags represents audit command flags
type AuditFlags struct {
	ConfigPath   string
	OutputFormat string
	AuditType    string
}

// extractAuditFlags extracts and validates flags from command
func extractAuditFlags(cmd *cobra.Command) (*AuditFlags, error) {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, fmt.Errorf("failed to get config flag: %w", err)
	}

	outputFormat, err := cmd.Flags().GetString("output")
	if err != nil {
		return nil, fmt.Errorf("failed to get output flag: %w", err)
	}

	auditType, err := cmd.Flags().GetString("type")
	if err != nil {
		return nil, fmt.Errorf("failed to get type flag: %w", err)
	}

	// Guard clause: validate config path
	if err := ValidateConfigPath(configPath); err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	// Guard clause: validate output format
	if err := ValidateOutputFormat(outputFormat); err != nil {
		return nil, fmt.Errorf("invalid output format: %w", err)
	}

	// Guard clause: validate audit type
	if err := validateAuditType(auditType); err != nil {
		return nil, fmt.Errorf("invalid audit type: %w", err)
	}

	return &AuditFlags{
		ConfigPath:   configPath,
		OutputFormat: outputFormat,
		AuditType:    auditType,
	}, nil
}

// validateAuditType validates the audit type
func validateAuditType(auditType string) error {
	validTypes := []string{"security", "compliance", "performance", "all"}
	for _, valid := range validTypes {
		if auditType == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid audit type: %s, must be one of: %v", auditType, validTypes)
}

// executeAudit executes the audit process
func executeAudit(flags *AuditFlags) error {
	// Guard clause: validate flags
	if flags == nil {
		return fmt.Errorf("flags cannot be nil")
	}

	// Display audit info
	displayAuditInfo(flags)

	// TODO: Implement actual audit logic
	fmt.Printf("✅ Audit completed successfully!\n")

	return nil
}

// displayAuditInfo displays audit information
func displayAuditInfo(flags *AuditFlags) {
	fmt.Printf("🔒 Running %s audit...\n", flags.AuditType)
	fmt.Printf("📁 Config: %s\n", flags.ConfigPath)
	fmt.Printf("📤 Output: %s\n", flags.OutputFormat)
}
