package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run:   runVersion,
	}
}

// runVersion executes the version command
func runVersion(cmd *cobra.Command, args []string) {
	displayVersionInfo()
}

// displayVersionInfo displays version information
func displayVersionInfo() {
	fmt.Printf("Praetorian CLI v0.0.1-alpha\n")
	fmt.Printf("Commit: dev\n")
	fmt.Printf("Built: unknown\n")
	fmt.Printf("Go version: 1.21+\n")
}
