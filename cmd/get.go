package cmd

import (
	"github.com/jahvon/tuikit/components"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/io"
	configio "github.com/jahvon/flow/internal/io/config"
	executableio "github.com/jahvon/flow/internal/io/executable"
	workspaceio "github.com/jahvon/flow/internal/io/workspace"
	"github.com/jahvon/flow/internal/vault"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	Short:   "Print a flow entity.",
}

var configGetCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Print the current global configuration values.",
	Args:    cobra.NoArgs,
	PreRun:  initInteractiveContainer,
	PostRun: waitForExit,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		userConfig := curCtx.UserConfig
		outputFormat := getFlagValue[string](cmd, *flags.OutputFormatFlag)
		if interactiveUIEnabled() {
			view := configio.NewUserConfigView(curCtx.InteractiveContainer, *userConfig, components.Format(outputFormat))
			curCtx.InteractiveContainer.SetView(view)
		} else {
			configio.PrintUserConfig(logger, outputFormat, userConfig)
		}
	},
}

var workspaceGetCmd = &cobra.Command{
	Use:     "workspace <name>",
	Aliases: []string{"ws"},
	Short:   "Print a workspace's configuration. If the name is omitted, the current workspace is used.",
	Args:    cobra.MaximumNArgs(1),
	PreRun:  initInteractiveContainer,
	PostRun: waitForExit,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		userConfig := curCtx.UserConfig

		var workspaceName string
		if len(args) == 1 {
			workspaceName = args[0]
		} else {
			workspaceName = userConfig.CurrentWorkspace
		}

		if _, found := userConfig.Workspaces[workspaceName]; !found {
			logger.Fatalf("workspace '%s' not found", workspaceName)
		}

		wsPath := userConfig.Workspaces[workspaceName]
		wsCfg, err := file.LoadWorkspaceConfig(workspaceName, wsPath)
		if err != nil {
			logger.FatalErr(errors.Wrap(err, "failure loading workspace config"))
		} else if wsCfg == nil {
			logger.Fatalf("config not found for workspace %s", workspaceName)
		}

		outputFormat := getFlagValue[string](cmd, *flags.OutputFormatFlag)
		if interactiveUIEnabled() {
			view := workspaceio.NewWorkspaceView(curCtx.InteractiveContainer, *wsCfg, components.Format(outputFormat))
			curCtx.InteractiveContainer.SetView(view)
		} else {
			workspaceio.PrintWorkspaceConfig(logger, outputFormat, wsCfg)
		}
	},
}

var executableGetCmd = &cobra.Command{
	Use:     "executable <verb> <id>",
	Aliases: []string{"exec"},
	Short:   "Print an executable flow by reference.",
	Long: "Print an executable by the executable's verb and ID.\nThe target executable's ID should be in the  " +
		"form of 'ws/ns:name' and the verb should match the target executable's verb or one of its aliases.\n\n" +
		"See" + io.ConfigDocsURL("executables", "Verb") + "for more information on executable verbs." +
		"See" + io.ConfigDocsURL("executable", "Ref") + "for more information on executable IDs.",
	Args:    cobra.ExactArgs(2),
	PreRun:  initInteractiveContainer,
	PostRun: waitForExit,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		verbStr := args[0]
		verb := config.Verb(verbStr)
		if err := verb.Validate(); err != nil {
			logger.FatalErr(err)
		}
		id := args[1]
		ref := config.NewRef(id, verb)

		exec, err := curCtx.ExecutableCache.GetExecutableByRef(logger, ref)
		if err != nil && errors.Is(cache.NewExecutableNotFoundError(ref.String()), err) {
			logger.Debugf("Executable %s not found in cache, syncing cache", ref)
			if err := curCtx.ExecutableCache.Update(logger); err != nil {
				logger.FatalErr(err)
			}
			exec, err = curCtx.ExecutableCache.GetExecutableByRef(logger, ref)
		}
		if err != nil {
			logger.FatalErr(err)
		} else if exec == nil {
			logger.Fatalf("executable %s not found", ref)
		}

		outputFormat := getFlagValue[string](cmd, *flags.OutputFormatFlag)
		if interactiveUIEnabled() {
			view := executableio.NewExecutableView(curCtx.InteractiveContainer, *exec, components.Format(outputFormat))
			curCtx.InteractiveContainer.SetView(view)
		} else {
			executableio.PrintExecutable(logger, outputFormat, exec)
		}
	},
}

var vaultGetCmd = &cobra.Command{
	Use:     "secret <name>",
	Aliases: []string{"scrt"},
	Short:   "Print the value of a secret in the flow secret vault.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		reference := args[0]
		asPlainText := getFlagValue[bool](cmd, *flags.OutputSecretAsPlainTextFlag)

		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		v := vault.NewVault(logger)
		secret, err := v.GetSecret(reference)
		if err != nil {
			logger.FatalErr(err)
		}

		if asPlainText {
			logger.PlainTextInfo(secret.PlainTextString())
		} else {
			logger.PlainTextInfo(secret.String())
		}
	},
}

func init() {
	registerFlagOrPanic(configGetCmd, *flags.OutputFormatFlag)
	getCmd.AddCommand(configGetCmd)

	registerFlagOrPanic(workspaceGetCmd, *flags.OutputFormatFlag)
	getCmd.AddCommand(workspaceGetCmd)

	registerFlagOrPanic(executableGetCmd, *flags.OutputFormatFlag)
	getCmd.AddCommand(executableGetCmd)

	registerFlagOrPanic(vaultGetCmd, *flags.OutputSecretAsPlainTextFlag)
	getCmd.AddCommand(vaultGetCmd)

	rootCmd.AddCommand(getCmd)
}
