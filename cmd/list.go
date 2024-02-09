package cmd

import (
	"fmt"

	"github.com/jahvon/tuikit/components"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/io"
	executableio "github.com/jahvon/flow/internal/io/executable"
	workspaceio "github.com/jahvon/flow/internal/io/workspace"
	"github.com/jahvon/flow/internal/vault"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Print a list of flow entities.",
}

var workspaceList = &cobra.Command{
	Use:     "workspaces",
	Aliases: []string{"ws"},
	Short:   "Print a list of the registered flow workspaces.",
	Args:    cobra.NoArgs,
	PreRun:  initInteractiveContainer,
	PostRun: waitForExit,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		outputFormat := getFlagValue[string](cmd, *flags.OutputFormatFlag)
		tagsFilter := getFlagValue[[]string](cmd, *flags.FilterTagFlag)

		logger.Debugf("Loading workspace configs from cache")
		workspaceCache, err := curCtx.WorkspacesCache.Get()
		if err != nil {
			logger.Fatalx("failed to load workspace configs from cache", "err", err)
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

		if interactiveUIEnabled() {
			view := workspaceio.NewWorkspaceListView(
				curCtx.InteractiveContainer,
				filteredWorkspaces,
				components.Format(outputFormat),
			)
			curCtx.InteractiveContainer.SetView(view)
		} else {
			workspaceio.PrintWorkspaceList(logger, io.OutputFormat(outputFormat), filteredWorkspaces)
		}
	},
}

var executableListCmd = &cobra.Command{
	Use:     "executables",
	Aliases: []string{"execs"},
	Short:   "Print a list of executable flows.",
	Args:    cobra.NoArgs,
	PreRun:  initInteractiveContainer,
	PostRun: waitForExit,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
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
			logger.FatalErr(err)
		}
		filteredExec := allExecs
		filteredExec = filteredExec.
			FilterByWorkspace(wsFilter).
			FilterByNamespace(nsFilter).
			FilterByVerb(config.Verb(verbFilter)).
			FilterByTags(tagsFilter)

		if interactiveUIEnabled() {
			view := executableio.NewExecutableListView(
				curCtx.InteractiveContainer,
				filteredExec,
				components.Format(outputFormat),
			)
			curCtx.InteractiveContainer.SetView(view)
		} else {
			executableio.PrintExecutableList(logger, io.OutputFormat(outputFormat), filteredExec)
		}
	},
}

var vaultSecretListCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"scrt"},
	Short:   "Print a list of secrets in the flow vault.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		asPlainText := getFlagValue[bool](cmd, *flags.OutputSecretAsPlainTextFlag)

		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		v := vault.NewVault()
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
	},
}

func init() {
	registerFlagOrPanic(workspaceList, *flags.OutputFormatFlag)
	registerFlagOrPanic(workspaceList, *flags.FilterTagFlag)
	listCmd.AddCommand(workspaceList)

	registerFlagOrPanic(executableListCmd, *flags.OutputFormatFlag)
	registerFlagOrPanic(executableListCmd, *flags.FilterWorkspaceFlag)
	registerFlagOrPanic(executableListCmd, *flags.FilterNamespaceFlag)
	registerFlagOrPanic(executableListCmd, *flags.FilterVerbFlag)
	registerFlagOrPanic(executableListCmd, *flags.FilterTagFlag)
	listCmd.AddCommand(executableListCmd)

	registerFlagOrPanic(vaultSecretListCmd, *flags.OutputSecretAsPlainTextFlag)
	listCmd.AddCommand(vaultSecretListCmd)

	rootCmd.AddCommand(listCmd)
}
