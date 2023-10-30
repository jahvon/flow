package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/io/ui"
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Interact with flow via the UI",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintUiFrame()
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}
