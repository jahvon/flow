package internal

import (
	"fmt"

	"github.com/jahvon/tuikit/components"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/cmd/internal/interactive"
	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	executableio "github.com/jahvon/flow/internal/io/executable"
	workspaceio "github.com/jahvon/flow/internal/io/workspace"
	"github.com/jahvon/flow/internal/vault"
)

func RegisterListCmd(ctx *context.Context, rootCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Print a list of flow entities.",
	}
	registerListWorkspaceCmd(ctx, listCmd)
	registerListExecutableCmd(ctx, listCmd)
	registerListSecretCmd(ctx, listCmd)
	rootCmd.AddCommand(listCmd)
}

func registerListWorkspaceCmd(ctx *context.Context, listCmd *cobra.Command) {
	workspaceCmd := &cobra.Command{
		Use:     "workspaces",
		Aliases: []string{"ws"},
		Short:   "Print a list of the registered flow workspaces.",
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveContainer(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { interactive.WaitForExit(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { listWorkspaceFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, workspaceCmd, *flags.OutputFormatFlag)
	RegisterFlag(ctx, workspaceCmd, *flags.FilterTagFlag)
	listCmd.AddCommand(workspaceCmd)
}

func listWorkspaceFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)
	tagsFilter := flags.ValueFor[[]string](ctx, cmd, *flags.FilterTagFlag, false)

	logger.Debugf("Loading workspace configs from cache")
	workspaceCache, err := ctx.WorkspacesCache.GetLatestData(logger)
	if err != nil {
		logger.Fatalx("failure loading workspace configs from cache", "err", err)
	}

	filteredWorkspaces := make([]config.WorkspaceConfig, 0)
	for name, ws := range workspaceCache.Workspaces {
		location := workspaceCache.WorkspaceLocations[name]
		ws.SetContext(name, location)
		if !ws.Tags.HasAnyTag(tagsFilter) {
			continue
		}
		filteredWorkspaces = append(filteredWorkspaces, *ws)
	}

	if len(filteredWorkspaces) == 0 {
		logger.Fatalf("no workspaces found")
	}

	if interactive.UIEnabled(ctx, cmd) {
		view := workspaceio.NewWorkspaceListView(
			ctx,
			filteredWorkspaces,
			components.Format(outputFormat),
		)
		ctx.InteractiveContainer.SetView(view)
	} else {
		workspaceio.PrintWorkspaceList(logger, outputFormat, filteredWorkspaces)
	}
}

func registerListExecutableCmd(ctx *context.Context, listCmd *cobra.Command) {
	executableCmd := &cobra.Command{
		Use:     "executables",
		Aliases: []string{"execs"},
		Short:   "Print a list of executable flows.",
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveContainer(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { interactive.WaitForExit(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { listExecutableFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, executableCmd, *flags.OutputFormatFlag)
	RegisterFlag(ctx, executableCmd, *flags.FilterWorkspaceFlag)
	RegisterFlag(ctx, executableCmd, *flags.FilterNamespaceFlag)
	RegisterFlag(ctx, executableCmd, *flags.FilterVerbFlag)
	RegisterFlag(ctx, executableCmd, *flags.FilterTagFlag)
	RegisterFlag(ctx, executableCmd, *flags.FilterExecSubstringFlag)
	listCmd.AddCommand(executableCmd)
}

func listExecutableFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
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
		FilterByVerb(config.Verb(verbFilter)).
		FilterByTags(tagsFilter).
		FilterBySubstring(substr)

	if interactive.UIEnabled(ctx, cmd) {
		runFunc := func(ref string) error { return runByRef(ctx, cmd, ref) }
		view := executableio.NewExecutableListView(
			ctx,
			filteredExec,
			components.Format(outputFormat),
			runFunc,
		)
		ctx.InteractiveContainer.SetView(view)
	} else {
		executableio.PrintExecutableList(logger, outputFormat, filteredExec)
	}
}

func registerListSecretCmd(ctx *context.Context, listCmd *cobra.Command) {
	vaultSecretListCmd := &cobra.Command{
		Use:     "secrets",
		Aliases: []string{"scrt"},
		Short:   "Print a list of secrets in the flow vault.",
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { listSecretFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, vaultSecretListCmd, *flags.OutputSecretAsPlainTextFlag)
	listCmd.AddCommand(vaultSecretListCmd)
}

func listSecretFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	asPlainText := flags.ValueFor[bool](ctx, cmd, *flags.OutputSecretAsPlainTextFlag, false)

	v := vault.NewVault(logger)
	secrets, err := v.GetAllSecrets()
	if err != nil {
		logger.FatalErr(err)
	}

	for ref, secret := range secrets {
		if asPlainText {
			logger.PlainTextInfo(fmt.Sprintf("%s: %s", ref, secret.PlainTextString()))
		} else {
			logger.PlainTextInfo(fmt.Sprintf("%s: %s", ref, secret.String()))
		}
	}
}
