package internal

import (
	"errors"
	"fmt"

	"github.com/jahvon/tuikit/types"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/cache"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	execIO "github.com/jahvon/flow/internal/io/executable"
	"github.com/jahvon/flow/internal/io/library"
	"github.com/jahvon/flow/types/executable"
)

func RegisterLibraryCmd(ctx *context.Context, rootCmd *cobra.Command) {
	libraryCmd := &cobra.Command{
		Use:     "library",
		Short:   "View and manage your library of workspaces and executables.",
		Aliases: []string{"lib"},
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { StartTUI(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { WaitForTUI(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { libraryFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, libraryCmd, *flags.FilterWorkspaceFlag)
	RegisterFlag(ctx, libraryCmd, *flags.FilterNamespaceFlag)
	RegisterFlag(ctx, libraryCmd, *flags.FilterVerbFlag)
	RegisterFlag(ctx, libraryCmd, *flags.FilterTagFlag)
	RegisterFlag(ctx, libraryCmd, *flags.FilterExecSubstringFlag)
	registerListExecutableCmd(ctx, libraryCmd)
	registerGetExecCmd(ctx, libraryCmd)
	rootCmd.AddCommand(libraryCmd)
}

func libraryFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	if !TUIEnabled(ctx, cmd) {
		logger.FatalErr(errors.New("library command requires an interactive terminal"))
	}

	wsFilter := flags.ValueFor[string](ctx, cmd, *flags.FilterWorkspaceFlag, false)
	if wsFilter == "." {
		wsFilter = ctx.Config.CurrentWorkspace
	}

	nsFilter := flags.ValueFor[string](ctx, cmd, *flags.FilterNamespaceFlag, false)
	if nsFilter == "." {
		nsFilter = ctx.Config.CurrentNamespace
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
	SetView(ctx, cmd, libraryModel)
}

func registerListExecutableCmd(ctx *context.Context, libraryCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "View a list of executables.",
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { StartTUI(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { WaitForTUI(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { listExecutableFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, listCmd, *flags.OutputFormatFlag)
	RegisterFlag(ctx, listCmd, *flags.FilterWorkspaceFlag)
	RegisterFlag(ctx, listCmd, *flags.FilterNamespaceFlag)
	RegisterFlag(ctx, listCmd, *flags.FilterVerbFlag)
	RegisterFlag(ctx, listCmd, *flags.FilterTagFlag)
	RegisterFlag(ctx, listCmd, *flags.FilterExecSubstringFlag)
	libraryCmd.AddCommand(listCmd)
}

func listExecutableFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	wsFilter := flags.ValueFor[string](ctx, cmd, *flags.FilterWorkspaceFlag, false)
	if wsFilter == "." {
		wsFilter = ctx.Config.CurrentWorkspace
	}

	nsFilter := flags.ValueFor[string](ctx, cmd, *flags.FilterNamespaceFlag, false)
	if nsFilter == "." {
		nsFilter = ctx.Config.CurrentNamespace
	}

	verbFilter := flags.ValueFor[string](ctx, cmd, *flags.FilterVerbFlag, false)
	tagsFilter := flags.ValueFor[[]string](ctx, cmd, *flags.FilterTagFlag, false)
	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)
	substr := flags.ValueFor[string](ctx, cmd, *flags.FilterExecSubstringFlag, false)

	allExecs, err := ctx.ExecutableCache.GetExecutableList(logger)
	if err != nil {
		logger.FatalErr(err)
	}
	filteredExec := allExecs
	filteredExec = filteredExec.
		FilterByWorkspace(wsFilter).
		FilterByNamespace(nsFilter).
		FilterByVerb(executable.Verb(verbFilter)).
		FilterByTags(tagsFilter).
		FilterBySubstring(substr)

	if TUIEnabled(ctx, cmd) {
		runFunc := func(ref string) error { return runByRef(ctx, cmd, ref) }
		view := execIO.NewExecutableListView(
			ctx,
			filteredExec,
			types.Format(outputFormat),
			runFunc,
		)
		SetView(ctx, cmd, view)
	} else {
		execIO.PrintExecutableList(logger, outputFormat, filteredExec)
	}
}

func registerGetExecCmd(ctx *context.Context, libraryCmd *cobra.Command) {
	getCmd := &cobra.Command{
		Use:     "get VERB ID",
		Aliases: []string{"find"},
		Short:   "View an executable's documentation. The executable is found by reference.",
		Long: "View an executable by the executable's verb and ID.\nThe target executable's ID should be in the  " +
			"form of 'ws/ns:name' and the verb should match the target executable's verb or one of its aliases.\n\n" +
			fmt.Sprintf(
				"See %s for more information on executable verbs.\n"+
					"See %s for more information on executable IDs.",
				io.TypesDocsURL("flowfile", "ExecutableVerb"),
				io.TypesDocsURL("flowfile", "ExecutableRef"),
			),
		Args:    cobra.ExactArgs(2),
		PreRun:  func(cmd *cobra.Command, args []string) { StartTUI(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { WaitForTUI(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { getExecFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, getCmd, *flags.OutputFormatFlag)
	libraryCmd.AddCommand(getCmd)
}

func getExecFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	verbStr := args[0]
	verb := executable.Verb(verbStr)
	if err := verb.Validate(); err != nil {
		logger.FatalErr(err)
	}
	id := args[1]
	ws, ns, name := executable.ParseExecutableID(id)
	if ws == "" {
		ws = ctx.CurrentWorkspace.AssignedName()
	}
	if ns == "" && ctx.Config.CurrentNamespace != "" {
		ns = ctx.Config.CurrentNamespace
	}
	id = executable.NewExecutableID(ws, ns, name)
	ref := executable.NewRef(id, verb)

	exec, err := ctx.ExecutableCache.GetExecutableByRef(logger, ref)
	if err != nil && errors.Is(cache.NewExecutableNotFoundError(ref.String()), err) {
		logger.Debugf("Executable %s not found in cache, syncing cache", ref)
		if err := ctx.ExecutableCache.Update(logger); err != nil {
			logger.FatalErr(err)
		}
		exec, err = ctx.ExecutableCache.GetExecutableByRef(logger, ref)
	}
	if err != nil {
		logger.FatalErr(err)
	} else if exec == nil {
		logger.Fatalf("executable %s not found", ref)
	}

	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)
	if TUIEnabled(ctx, cmd) {
		runFunc := func(ref string) error { return runByRef(ctx, cmd, ref) }
		view := execIO.NewExecutableView(ctx, exec, types.Format(outputFormat), runFunc)
		SetView(ctx, cmd, view)
	} else {
		execIO.PrintExecutable(logger, outputFormat, exec)
	}
}
