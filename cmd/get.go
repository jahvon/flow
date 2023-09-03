package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/io"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	Short:   "Get the current value of a configuration, environment, or workspace option.",
}

// getWorkspaceCmd represents the get workspace subcommand
var getWorkspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"w"},
	Short:   "Print the current workspace.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}
		io.PrintNotice(rootCfg.CurrentWorkspace)
	},
}

func init() {
	getCmd.AddCommand(getWorkspaceCmd)
	rootCmd.AddCommand(getCmd)
}
