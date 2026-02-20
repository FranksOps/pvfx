// Package cmd provides the CLI commands for pvfx.
//
// Cobra uses a root command with subcommands. Each subcommand
// is its own file in this package. The root command typically
// handles global flags and persistent configuration.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	// Viper is typically used with Cobra for config management,
	// though we're using flags directly here for simplicity.
)

// rootCmd represents the base command when called without any subcommands.
// cobra.Command is a struct with many optional fields - we configure
// the most commonly used ones here.
var rootCmd = &cobra.Command{
	Use:   "pvfx",
	Short: "LVM and XFS filesystem management tool",
	Long: `A CLI tool to view and manage XFS filesystems on LVM.

Examples:
  pvfx status              # Show current LVM and filesystem status
  pvfx grow --lv /dev/ubuntu-vg/ubuntu-lv  # Grow the filesystem
  pvfx grow --all          # Grow all XFS filesystems`,
	// RunE returns an error - returning nil means success.
	// This is idiomatic Cobra: return an error and let the
	// framework handle exit codes.
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main() - it only needs to happen once.
func Execute() {
	// BindFlags traverses all parents of cmd and binds flags to them.
	// This ensures that persistent flags (set on parent) are available
	// to child commands.
	if err := rootCmd.Execute(); err != nil {
		// In production CLI tools, we print errors to stderr
		// and exit with non-zero status code.
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags that apply to this command
	// and all subcommands.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pvfx.yaml)")
}