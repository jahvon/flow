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
		ui.PrintUiFrame(
			ui.WithCurrentState(ui.IdleState),
			ui.WithCurrentWorkspace("ws"),
			ui.WithCurrentNamespace("ns"),
			ui.WithCurrentView("workspaces"),
			ui.WithCurrentFilter([]string{"tag1", "tag2"}),
			ui.WithNotice("Testing..."),
			ui.WithObjectContent(&ui.TableData{
				Headers: []string{"key", "value"},
				Rows:    [][]string{{"key1", "value1"}, {"key2", "value2"}},
			}),
		)
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}
