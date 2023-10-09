package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/executable"
	"github.com/jahvon/flow/internal/config"
	openagent "github.com/jahvon/flow/internal/executable/open"
	"github.com/jahvon/flow/internal/io"
)

var openCmd = &cobra.Command{
	Use:     "open",
	Aliases: []string{"o"},
	GroupID: ExecutableGroup.ID,
	Short:   "Execute open uri flow.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		agent := openagent.NewAgent()

		_, executable, err := executable.ArgsToExecutable(args, agent.Name(), rootCfg)
		if err != nil {
			io.PrintErrorAndExit(err)
		} else if executable == nil {
			log.Fatal().Msg("executable is nil")
		}

		err = agent.Exec(*executable)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		log.Info().Msg("open uri flow completed")
	},
}

func init() {
	openCmd.AddGroup(ExecutableGroup)
	rootCmd.AddCommand(openCmd)
}
