package cmd

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/io/library"
)

var libraryCmd = &cobra.Command{
	Use:     "library",
	Short:   "View and manage your library of workspaces and executables.",
	Aliases: []string{"lib"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		if !interactiveUIEnabled() {
			logger.FatalErr(errors.New("library command requires an interactive terminal"))
		}

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

		allExecs, err := curCtx.ExecutableCache.GetExecutableList(logger)
		if err != nil {
			logger.FatalErr(err)
		}
		allWs, err := curCtx.WorkspacesCache.GetWorkspaceConfigList(logger)
		if err != nil {
			logger.FatalErr(err)
		}

		libraryModel := library.NewLibrary(
			curCtx, allWs, allExecs,
			library.Filter{
				Workspace: wsFilter,
				Namespace: nsFilter,
				Verb:      config.Verb(verbFilter),
				Tags:      tagsFilter,
			})
		program := tea.NewProgram(
			libraryModel,
			tea.WithAltScreen(),
			tea.WithContext(curCtx.Ctx),
		)
		if _, err := program.Run(); err != nil {
			logger.FatalErr(err)
		}
	},
}

func init() {
	registerFlagOrPanic(libraryCmd, *flags.FilterWorkspaceFlag)
	registerFlagOrPanic(libraryCmd, *flags.FilterNamespaceFlag)
	registerFlagOrPanic(libraryCmd, *flags.FilterVerbFlag)
	registerFlagOrPanic(libraryCmd, *flags.FilterTagFlag)
	rootCmd.AddCommand(libraryCmd)
}
