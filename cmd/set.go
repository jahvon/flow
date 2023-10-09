package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/io"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s"},
	GroupID: CrudGroup.ID,
	Short:   "Update a configuration, environment, or workspace option.",
}

// setWorkspaceCmd represents the set workspace subcommand
var setWorkspaceCmd = &cobra.Command{
	Use:     "workspace <name>",
	Aliases: []string{"w"},
	Short:   "Change the current workspace.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workspace := args[0]
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}

		if err := config.SetCurrentWorkspace(rootCfg, workspace); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess("Workspace set to " + workspace)
	},
}

func init() {
	setCmd.AddCommand(setWorkspaceCmd)
	rootCmd.AddCommand(setCmd)
}
