package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var growCmd = &cobra.Command{
	Use:   "grow",
	Short: "Grow LVM logical volume and XFS filesystem",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Grow command")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(growCmd)
}