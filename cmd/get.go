package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/tbox/internal/io"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the current value of a configuration, environment, or workspace option.",
}

// getWorkspaceCmd represents the get workspace subcommand
var getWorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Print the current workspace.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		io.PrintInfo(currentConfig.CurrentWorkspace)
	},
}

func init() {
	getCmd.AddCommand(getWorkspaceCmd)
	rootCmd.AddCommand(getCmd)
}
