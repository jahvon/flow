package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/io"
	executableio "github.com/jahvon/flow/internal/io/executable"
)

var executablesCmd = &cobra.Command{
	Use:     "executables",
	Aliases: []string{"execs"},
	Short:   "Manage flow executables.",
}

var executableGetCmd = &cobra.Command{
	Use:     "get <verb> <id>",
	Aliases: []string{"g"},
	Short:   "Get an executable by reference.",
	Long: "Get an executable by the executable's verb and ID.\nThe target executable's ID should be in the form of  " +
		"'ws/ns:name' and the verb should match the target executable's verb or one of its aliases.\n\n" +
		"See" + io.DocsURL("executable-verbs") + "for more information on executable verbs." +
		"See" + io.DocsURL("executable-ids") + "for more information on executable IDs.",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		verbStr := args[0]
		verb := config.Verb(verbStr)
		if err := verb.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}
		id := args[1]
		ref := config.NewRef(id, verb)

		exec, err := curCtx.ExecutableCache.GetExecutableByRef(ref)
		if err != nil {
			io.PrintErrorAndExit(err)
		} else if exec == nil {
			io.PrintErrorAndExit(fmt.Errorf("executable %s not found", ref))
		}

		outputFormat := getFlagValue[string](cmd, *flags.OutputFormatFlag)
		executableio.PrintExecutable(io.OutputFormat(outputFormat), exec)
	},
}

var executablesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "ParameterList executables.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		wsFilter := getFlagValue[string](cmd, *flags.FilterWorkspaceFlag)
		if wsFilter == "." {
			wsFilter = curCtx.UserConfig.CurrentWorkspace
		}

		nsFilter := getFlagValue[string](cmd, *flags.FilterNamespaceFlag)
		if nsFilter == "." {
			nsFilter = curCtx.UserConfig.CurrentNamespace
		}

		verbFilter := getFlagValue[string](cmd, *flags.FilterVerbFlag)
		tagsFilter := getFlagValue[[]string](cmd, *flags.FilterTagFlag)
		outputFormat := getFlagValue[string](cmd, *flags.OutputFormatFlag)

		allExecs, err := curCtx.ExecutableCache.GetExecutableList()
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		filteredExec := allExecs
		filteredExec = filteredExec.
			FilterByWorkspace(wsFilter).
			FilterByNamespace(nsFilter).
			FilterByVerb(config.Verb(verbFilter)).
			FilterByTags(tagsFilter)
		executableio.PrintExecutableList(io.OutputFormat(outputFormat), filteredExec)
	},
}

func init() {
	registerFlagOrPanic(executableGetCmd, *flags.SyncCacheFlag)
	registerFlagOrPanic(executableGetCmd, *flags.OutputFormatFlag)
	executablesCmd.AddCommand(executableGetCmd)

	registerFlagOrPanic(executablesListCmd, *flags.SyncCacheFlag)
	registerFlagOrPanic(executablesListCmd, *flags.OutputFormatFlag)
	registerFlagOrPanic(executablesListCmd, *flags.FilterWorkspaceFlag)
	registerFlagOrPanic(executablesListCmd, *flags.FilterNamespaceFlag)
	registerFlagOrPanic(executablesListCmd, *flags.FilterVerbFlag)
	registerFlagOrPanic(executablesListCmd, *flags.FilterTagFlag)
	executablesCmd.AddCommand(executablesListCmd)

	rootCmd.AddCommand(executablesCmd)
}
