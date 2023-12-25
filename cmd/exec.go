package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/exec"
	"github.com/jahvon/flow/internal/runner/launch"
	"github.com/jahvon/flow/internal/runner/parallel"
	"github.com/jahvon/flow/internal/runner/serial"
)

var execCmd = &cobra.Command{
	Use:     "exec <executable-id>",
	Aliases: config.ValidVerbs,
	Short:   "Execute a flow by ID.",
	Long: "Execute a flow where <executable-id> is the target executable's ID in the form of 'ws/ns:name'.\n" +
		"The flow subcommand used should match the target executable's verb or one of its aliases.\n\n" +
		"See" + io.DocsURL("executable-verbs") + "for more information on executable verbs." +
		"See" + io.DocsURL("executable-ids") + "for more information on executable IDs.",
	Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		runner.RegisterRunner(exec.NewRunner())
		runner.RegisterRunner(launch.NewRunner())
		runner.RegisterRunner(serial.NewRunner())
		runner.RegisterRunner(parallel.NewRunner())
	},
	Run: func(cmd *cobra.Command, args []string) {
		verbStr := cmd.CalledAs()
		verb := config.Verb(verbStr)
		if err := verb.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		idArg := args[0]
		ref := context.ExpandRef(curCtx, config.NewRef(idArg, verb))
		executable, err := curCtx.ExecutableCache.GetExecutableByRef(ref)
		if err != nil {
			io.PrintErrorAndExit(err)
		} else if executable == nil {
			io.PrintErrorAndExit(fmt.Errorf("executable %s not found", ref))
		}

		if err := executable.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		if !executable.IsExecutableFromWorkspace(curCtx.UserConfig.CurrentWorkspace) {
			io.PrintErrorAndExit(
				fmt.Errorf(
					"executable '%s' cannot be executed from workspace %s",
					ref,
					curCtx.UserConfig.CurrentWorkspace,
				),
			)
		}

		if err := runner.Exec(curCtx, executable); err != nil {
			io.PrintErrorAndExit(err)
		}
		log.Info().Msg(fmt.Sprintf("%s flow completed", ref))
	},
}

func init() {
	registerFlagOrPanic(execCmd, *flags.SyncCacheFlag)
	rootCmd.AddCommand(execCmd)
}
