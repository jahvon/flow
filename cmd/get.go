package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/io"
	configio "github.com/jahvon/flow/internal/io/config"
	executableio "github.com/jahvon/flow/internal/io/executable"
	"github.com/jahvon/flow/internal/io/ui/views"
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
	PreRun:  startApp,
	PostRun: waitForExit,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		userConfig := curCtx.UserConfig
		if userConfig == nil {
			logger.FatalErr(fmt.Errorf("failed to load user config"))
		}
		if err := userConfig.Validate(); err != nil {
			logger.FatalErr(err)
		}

		outputFormat := getFlagValue[string](cmd, *flags.OutputFormatFlag)
		if interactiveUIEnabled() {
			viewBuilder := views.NewUserConfigView(curCtx.App, *userConfig, config.OutputFormat(outputFormat))
			curCtx.App.BuildAndSetView(viewBuilder)
		} else {
			configio.PrintUserConfig(io.OutputFormat(outputFormat), userConfig)
		}
	},
}

var workspaceGetCmd = &cobra.Command{
	Use:     "workspace <name>",
	Aliases: []string{"ws"},
	Short:   "Print a workspace's configuration. If the name is omitted, the current workspace is used.",
	Args:    cobra.MaximumNArgs(1),
	PreRun:  startApp,
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
			logger.Fatalf("failed loading workspace config: %v", err)
		} else if wsCfg == nil {
			logger.Fatalf("config not found for workspace %s", workspaceName)
		}

		outputFormat := getFlagValue[string](cmd, *flags.OutputFormatFlag)
		if interactiveUIEnabled() {
			viewBuilder := views.NewWorkspaceView(curCtx.App, *wsCfg, config.OutputFormat(outputFormat))
			curCtx.App.BuildAndSetView(viewBuilder)
		} else {
			workspaceio.PrintWorkspaceConfig(io.OutputFormat(outputFormat), wsCfg)
		}
	},
}

var executableGetCmd = &cobra.Command{
	Use:     "executable <verb> <id>",
	Aliases: []string{"exec"},
	Short:   "Print an executable flow by reference.",
	Long: "Print an executable by the executable's verb and ID.\nThe target executable's ID should be in the  " +
		"form of 'ws/ns:name' and the verb should match the target executable's verb or one of its aliases.\n\n" +
		"See" + io.DocsURL("executable-verbs") + "for more information on executable verbs." +
		"See" + io.DocsURL("executable-ids") + "for more information on executable IDs.",
	Args:    cobra.ExactArgs(2),
	PreRun:  startApp,
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

		exec, err := curCtx.ExecutableCache.GetExecutableByRef(ref)
		if err != nil {
			logger.FatalErr(err)
		} else if exec == nil {
			logger.Fatalf("executable %s not found", ref)
		}

		outputFormat := getFlagValue[string](cmd, *flags.OutputFormatFlag)
		if interactiveUIEnabled() {
			viewBuilder := views.NewExecutableView(curCtx.App, *exec, config.OutputFormat(outputFormat))
			curCtx.App.BuildAndSetView(viewBuilder)
		} else {
			executableio.PrintExecutable(io.OutputFormat(outputFormat), exec)
		}
	},
}

var vaultGetCmd = &cobra.Command{
	Use:     "secret <name>",
	Aliases: []string{"scrt"},
	Short:   "Print the value of a secret in the flow secret vault.",
	Args:    cobra.ExactArgs(1),
	PreRun:  setTermView,
	PostRun: exitApp,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		reference := args[0]
		asPlainText := getFlagValue[bool](cmd, *flags.OutputSecretAsPlainTextFlag)

		v := vault.NewVault()
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
