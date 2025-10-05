package cli

import "github.com/spf13/cobra"

// RegisterCommands registers all CLI commands with the root command
func RegisterCommands(rootCmd *cobra.Command) {
	// Guard clause: validate root command
	if rootCmd == nil {
		panic("root command cannot be nil")
	}

	// Add all commands
	rootCmd.AddCommand(NewValidateCommand())
	rootCmd.AddCommand(NewAuditCommand())
	rootCmd.AddCommand(NewInitCommand())
	rootCmd.AddCommand(NewVersionCommand())
}