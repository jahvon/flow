package internal

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/cache"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	execIO "github.com/jahvon/flow/internal/io/executable"
	"github.com/jahvon/flow/internal/io/library"
	"github.com/jahvon/flow/types/executable"
)

func RegisterBrowseCmd(ctx *context.Context, rootCmd *cobra.Command) {
	browseCmd := &cobra.Command{
		Use:     "browse [EXECUTABLE-REFERENCE]",
		Short:   "Discover and explore available executables.",
		Aliases: []string{"ls", "library"},
		Long: "Browse executables across workspaces.\n\n" +
			"  flow browse                # Interactive multi-pane executable browser\n" +
			"  flow browse --list         # Simple list view of executables\n" +
			"  flow browse VERB [ID]      # Detailed view of specific executable\n\n" +
			fmt.Sprintf(
				"See %s for more information on executable verbs and "+
					"%s for more information on executable references.",
				io.TypesDocsURL("flowfile", "executableverb"),
				io.TypesDocsURL("flowfile", "executableref"),
			),
		Args: cobra.MaximumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			execList, err := ctx.ExecutableCache.GetExecutableList(ctx.Logger)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			execIDs := make([]string, 0, len(execList))
			for _, e := range execList {
				execIDs = append(execIDs, e.ID())
			}
			return execIDs, cobra.ShellCompDirectiveNoFileComp
		},
		PreRun:  func(cmd *cobra.Command, args []string) { StartTUI(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { WaitForTUI(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { browseFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, browseCmd, *flags.ListFlag)
	RegisterFlag(ctx, browseCmd, *flags.OutputFormatFlag)
	RegisterFlag(ctx, browseCmd, *flags.FilterWorkspaceFlag)
	RegisterFlag(ctx, browseCmd, *flags.FilterNamespaceFlag)
	RegisterFlag(ctx, browseCmd, *flags.FilterVerbFlag)
	RegisterFlag(ctx, browseCmd, *flags.FilterTagFlag)
	RegisterFlag(ctx, browseCmd, *flags.FilterExecSubstringFlag)
	RegisterFlag(ctx, browseCmd, *flags.AllNamespacesFlag)
	rootCmd.AddCommand(browseCmd)
}

func browseFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	if len(args) >= 1 {
		viewExecutable(ctx, cmd, args)
		return
	}

	listMode := flags.ValueFor[bool](ctx, cmd, *flags.ListFlag, false)
	if listMode || !TUIEnabled(ctx, cmd) {
		listExecutables(ctx, cmd, args)
		return
	}

	executableLibrary(ctx, cmd, args)
}

func executableLibrary(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	if !TUIEnabled(ctx, cmd) {
		logger.FatalErr(errors.New("interactive discovery requires an interactive terminal"))
	}

	wsFilter := flags.ValueFor[string](ctx, cmd, *flags.FilterWorkspaceFlag, false)
	switch wsFilter {
	case ".":
		wsFilter = ctx.Config.CurrentWorkspace
	case executable.WildcardWorkspace:
		wsFilter = ""
	}

	nsFilter := flags.ValueFor[string](ctx, cmd, *flags.FilterNamespaceFlag, false)
	allNs := flags.ValueFor[bool](ctx, cmd, *flags.AllNamespacesFlag, false)
	switch {
	case allNs && nsFilter != "":
		logger.PlainTextWarn("cannot use both --all and --namespace flags, ignoring --namespace")
		fallthrough
	case allNs:
		nsFilter = executable.WildcardNamespace
	case nsFilter == ".":
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
		io.Theme(ctx.Config.Theme.String()),
		runFunc,
	)
	SetView(ctx, cmd, libraryModel)
}

func listExecutables(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	wsFilter := flags.ValueFor[string](ctx, cmd, *flags.FilterWorkspaceFlag, false)
	if wsFilter == "." {
		wsFilter = ctx.Config.CurrentWorkspace
	}

	nsFilter := flags.ValueFor[string](ctx, cmd, *flags.FilterNamespaceFlag, false)
	allNs := flags.ValueFor[bool](ctx, cmd, *flags.AllNamespacesFlag, false)
	switch {
	case allNs && nsFilter != "":
		logger.PlainTextWarn("cannot use both --all and --namespace flags, ignoring --namespace")
		fallthrough
	case allNs:
		nsFilter = executable.WildcardNamespace
	case nsFilter == ".":
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
		view := execIO.NewExecutableListView(ctx, filteredExec, runFunc)
		SetView(ctx, cmd, view)
	} else {
		execIO.PrintExecutableList(logger, outputFormat, filteredExec)
	}
}

func viewExecutable(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger

	verbStr := args[0]
	verb := executable.Verb(verbStr)
	if err := verb.Validate(); err != nil {
		logger.FatalErr(err)
	}

	var execID string
	if len(args) > 1 {
		id := args[1]
		ws, ns, name := executable.MustParseExecutableID(id)
		if ws == executable.WildcardWorkspace {
			ws = ctx.CurrentWorkspace.AssignedName()
		}
		if ns == executable.WildcardNamespace && ctx.Config.CurrentNamespace != "" {
			ns = ctx.Config.CurrentNamespace
		}
		execID = executable.NewExecutableID(ws, ns, name)
	}
	ref := executable.NewRef(execID, verb)

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
		view := execIO.NewExecutableView(ctx, exec, runFunc)
		SetView(ctx, cmd, view)
	} else {
		execIO.PrintExecutable(logger, outputFormat, exec)
	}
}
