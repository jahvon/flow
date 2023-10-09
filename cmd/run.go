package cmd

import (
	"fmt"

	"github.com/jahvon/flow/internal/cmd/executable"
	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/executable/run"
	"github.com/jahvon/flow/internal/io"
	"github.com/spf13/cobra"
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
			log.Panic().Msg("executable is nil")
		}

		err = agent.Exec(*executable)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		log.Info().Msg(fmt.Sprintf("run %s flow completed", executable.Name))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
