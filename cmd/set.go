package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/tbox/internal/config"
	"github.com/jahvon/tbox/internal/io"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Update a configuration, environment, or workspace option.",
}

// setWorkspaceCmd represents the set workspace subcommand
var setWorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Change the current workspace.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workspace := args[0]
		if err := config.SetCurrentWorkspace(currentConfig, workspace); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess("Workspace set to " + workspace)
	},
}

func init() {
	setCmd.AddCommand(setWorkspaceCmd)
	rootCmd.AddCommand(setCmd)
}
