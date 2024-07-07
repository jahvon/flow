package internal

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/cmd/internal/interactive"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/library"
	"github.com/jahvon/flow/types/executable"
)

func RegisterLibraryCmd(ctx *context.Context, rootCmd *cobra.Command) {
	libraryCmd := &cobra.Command{
		Use:     "library",
		Short:   "View and manage your library of workspaces and executables.",
		Aliases: []string{"lib"},
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveContainer(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { interactive.WaitForExit(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { libraryFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, libraryCmd, *flags.FilterWorkspaceFlag)
	RegisterFlag(ctx, libraryCmd, *flags.FilterNamespaceFlag)
	RegisterFlag(ctx, libraryCmd, *flags.FilterVerbFlag)
	RegisterFlag(ctx, libraryCmd, *flags.FilterTagFlag)
	RegisterFlag(ctx, libraryCmd, *flags.FilterExecSubstringFlag)
	rootCmd.AddCommand(libraryCmd)
}

func libraryFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	if !interactive.UIEnabled(ctx, cmd) {
		logger.FatalErr(errors.New("library command requires an interactive terminal"))
	}

	wsFilter := flags.ValueFor[string](ctx, cmd, *flags.FilterWorkspaceFlag, false)
	if wsFilter == "." {
		wsFilter = ctx.UserConfig.CurrentWorkspace
	}

	nsFilter := flags.ValueFor[string](ctx, cmd, *flags.FilterNamespaceFlag, false)
	if nsFilter == "." {
		nsFilter = ctx.UserConfig.CurrentNamespace
	}

	verbFilter := flags.ValueFor[string](ctx, cmd, *flags.FilterVerbFlag, false)
	tagsFilter := flags.ValueFor[[]string](ctx, cmd, *flags.FilterTagFlag, false)
	subStr := flags.ValueFor[string](ctx, cmd, *flags.FilterExecSubstringFlag, false)

	allExecs, err := ctx.ExecutableCache.GetExecutableList(logger)
	if err != nil {
		logger.FatalErr(err)
	}
	allWs, err := ctx.WorkspacesCache.GetWorkspaceConfigList(logger)
	if err != nil {
		logger.FatalErr(err)
	}

	runFunc := func(ref string) error { return runByRef(ctx, cmd, ref) }
	libraryModel := library.NewLibraryView(
		ctx, allWs, allExecs,
		library.Filter{
			Workspace: wsFilter,
			Namespace: nsFilter,
			Verb:      executable.Verb(verbFilter),
			Tags:      tagsFilter,
			Substring: subStr,
		},
		io.Theme(),
		runFunc,
	)
	ctx.InteractiveContainer.SetView(libraryModel)
}
