package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/executable"
	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/executable/run"
	"github.com/jahvon/flow/internal/io"
)

var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	GroupID: ExecutableGroup.ID,
	Short:   "Execute run flow.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		agent := run.NewAgent()

		_, executable, err := executable.ArgsToExecutable(args, agent.Name(), rootCfg)
		if err != nil {
			io.PrintErrorAndExit(err)
		} else if executable == nil {
			log.Fatal().Msg("executable is nil")
		}

		err = agent.Exec(executable.Spec, nil)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		log.Info().Msg("run flow completed")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
