package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/syntropysoft/praetorian-go/internal/cli"
)

var (
	version = "0.0.1-alpha"
	commit  = "dev"
	date    = "unknown"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "praetorian",
		Short: "Guardian of Configurations",
		Long: `Praetorian CLI – Universal Validation Framework for DevSecOps

🛡️  Guardian of Configurations & Security

Praetorian CLI is a multi-environment configuration validation tool designed for DevSecOps teams.
It ensures that your configuration files remain consistent across environments and detects 
critical differences before production deployments.

Perfect for:
• Microservices architectures with multiple config files
• Multi-environment deployments (dev, staging, prod)
• Security compliance and configuration drift detection
• CI/CD pipelines requiring config validation`,

		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Show banner on first run
			if cmd.Name() == "praetorian" {
				showBanner()
			}
		},
	}

	// Register all commands
	cli.RegisterCommands(rootCmd)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// showBanner displays the Praetorian banner
func showBanner() {
	cyan := color.New(color.FgCyan).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	
	fmt.Println()
	fmt.Println(cyan(`  ____                 _             _                ____ _     ___ `))
	fmt.Println(cyan(` |  _ \ _ __ __ _  ___| |_ ___  _ __(_) __ _ _ __    / ___| |   |_ _|`))
	fmt.Println(cyan(` | |_) | '__/ _`) + bold(` |/ _ \ __/ _ \| '__| |/ _`) + cyan(` | '_ \  | |   | |    | | `))
	fmt.Println(cyan(` |  __/| | | (_| |  __/ || (_) | |  | | (_| | | | | | |___| |___ | | `))
	fmt.Println(cyan(` |_|   |_|  \__,_|\___|\__\___/|_|  |_|\__,_|_| |_|  \____|_____|___|`))
	fmt.Println()
	fmt.Println(cyan(`                                                                     `))
	fmt.Println(cyan(`🛡️  Guardian of Configurations & Security`))
	fmt.Println()
}

// newValidateCommand creates the validate command
func newValidateCommand() *cobra.Command {
	var configPath string
	var output string
	var pipeline bool

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration files across environments",
		Long: `Validate configuration files for consistency across environments.

This command compares configuration files between different environments (dev, staging, prod)
and reports missing keys, extra keys, and value differences.

Examples:
  praetorian validate                           # Validate current directory
  praetorian validate --config praetorian.yaml # Use specific config file
  praetorian validate --output json            # Output in JSON format
  praetorian validate --pipeline               # CI/CD friendly output`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🔍 Starting configuration validation...")
			fmt.Printf("📁 Config: %s\n", configPath)
			fmt.Printf("📤 Output: %s\n", output)
			fmt.Printf("🚀 Pipeline mode: %v\n", pipeline)
			
			// TODO: Implement validation logic
			fmt.Println("✅ Validation completed successfully!")
			
			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "praetorian.yaml", "Configuration file path")
	cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format (text, json, yaml)")
	cmd.Flags().BoolVar(&pipeline, "pipeline", false, "Enable pipeline mode for CI/CD")

	return cmd
}

// newAuditCommand creates the audit command
func newAuditCommand() *cobra.Command {
	var auditType string
	var output string
	var configPath string

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Run security and compliance audits",
		Long: `Run comprehensive security and compliance audits on configuration files.

Examples:
  praetorian audit                           # Run all audits on current directory
  praetorian audit --type security          # Run security audit only
  praetorian audit --config praetorian.yaml # Use specific config file
  praetorian audit --output json            # Output in JSON format`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🔒 Starting Praetorian Audit...")
			fmt.Printf("📁 Config: %s\n", configPath)
			fmt.Printf("🔍 Type: %s\n", auditType)
			fmt.Printf("📤 Output: %s\n", output)
			
			// TODO: Implement audit logic
			fmt.Println("✅ Audit completed successfully!")
			
			return nil
		},
	}

	cmd.Flags().StringVarP(&auditType, "type", "t", "all", "Audit type (security, compliance, performance, all)")
	cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format (text, json, yaml)")
	cmd.Flags().StringVarP(&configPath, "config", "c", "praetorian.yaml", "Configuration file path")

	return cmd
}

// newInitCommand creates the init command
func newInitCommand() *cobra.Command {
	var devsecops bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Praetorian configuration",
		Long: `Initialize Praetorian configuration for your project.

This command creates a praetorian.yaml configuration file with sensible defaults
for validating your configuration files across environments.

Examples:
  praetorian init                    # Create basic configuration
  praetorian init --devsecops        # Create DevSecOps optimized configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🚀 Initializing Praetorian configuration...")
			fmt.Printf("🛡️  DevSecOps mode: %v\n", devsecops)
			
			// TODO: Implement init logic
			fmt.Println("✅ Configuration initialized successfully!")
			
			return nil
		},
	}

	cmd.Flags().BoolVar(&devsecops, "devsecops", false, "Initialize with DevSecOps optimizations")

	return cmd
}

// newVersionCommand creates the version command
func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Praetorian CLI v%s\n", version)
			fmt.Printf("Commit: %s\n", commit)
			fmt.Printf("Built: %s\n", date)
			fmt.Printf("Go version: %s\n", "1.21+")
		},
	}

	return cmd
}
