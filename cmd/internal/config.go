package internal

import (
	"fmt"
	"strconv"
	"strings"

	io2 "github.com/jahvon/tuikit/io"
	"github.com/jahvon/tuikit/types"
	"github.com/jahvon/tuikit/views"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io"
	config2 "github.com/jahvon/flow/internal/io/config"
	"github.com/jahvon/flow/types/config"
)

func RegisterConfigCmd(ctx *context.Context, rootCmd *cobra.Command) {
	setCmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"c", "cfg"},
		Short:   "Update flow configuration values.",
	}
	registerConfigResetCmd(ctx, setCmd)
	registerSetConfigCmd(ctx, setCmd)
	registerViewConfigCmd(ctx, setCmd)
	registerSetWorkspaceCmd(ctx, setCmd)
	registerSetNamespaceCmd(ctx, setCmd)
	registerSetWorkspaceModeCmd(ctx, setCmd)
	registerSetLogModeCmd(ctx, setCmd)
	registerSetTUICmd(ctx, setCmd)
	registerRegisterTemplateCmd(ctx, setCmd)
	registerSetSecretCmd(ctx, setCmd)
	rootCmd.AddCommand(setCmd)
}

func registerConfigResetCmd(ctx *context.Context, configCmd *cobra.Command) {
	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Restore the default flow configuration values. This will overwrite the current configuration.",
		Args:  cobra.NoArgs,
		Run:   func(cmd *cobra.Command, args []string) { resetConfigFunc(ctx, cmd, args) },
	}
	configCmd.AddCommand(resetCmd)
}

func resetConfigFunc(ctx *context.Context, _ *cobra.Command, _ []string) {
	logger := ctx.Logger
	form, err := views.NewForm(
		io.Theme(),
		ctx.StdIn(),
		ctx.StdOut(),
		&views.FormField{
			Key:   "confirm",
			Type:  views.PromptTypeConfirm,
			Title: "This will overwrite your current flow configurations. Are you sure you want to continue?",
		})
	if err != nil {
		logger.FatalErr(err)
	}
	if err := form.Run(ctx.Ctx); err != nil {
		logger.FatalErr(err)
	}
	resp := form.FindByKey("confirm").Value()
	if truthy, _ := strconv.ParseBool(resp); !truthy {
		logger.Warnf("Aborting")
		return
	}

	if err := filesystem.InitConfig(); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess("Restored flow configurations")
}

func registerSetConfigCmd(ctx *context.Context, configCmd *cobra.Command) {
	setCmd := &cobra.Command{
		Use:     "set",
		Aliases: []string{"s", "update"},
		Short:   "Update flow configuration values.",
	}
	registerSetWorkspaceCmd(ctx, setCmd)
	registerSetNamespaceCmd(ctx, setCmd)
	registerSetWorkspaceModeCmd(ctx, setCmd)
	registerSetLogModeCmd(ctx, setCmd)
	registerSetTUICmd(ctx, setCmd)
	configCmd.AddCommand(setCmd)
}

func registerSetWorkspaceCmd(ctx *context.Context, setCmd *cobra.Command) {
	workspaceCmd := &cobra.Command{
		Use:     "workspace NAME",
		Aliases: []string{"ws"},
		Short:   "Change the current workspace.",
		Args:    cobra.ExactArgs(1),
		PreRun:  func(cmd *cobra.Command, args []string) { printContext(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { setWorkspaceFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, workspaceCmd, *flags.FixedWsModeFlag)
	setCmd.AddCommand(workspaceCmd)
}

func setWorkspaceFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	workspace := args[0]
	userConfig := ctx.Config
	if _, found := userConfig.Workspaces[workspace]; !found {
		logger.Fatalf("workspace %s not found", workspace)
	}
	userConfig.CurrentWorkspace = workspace
	fixedMode := flags.ValueFor[bool](ctx, cmd, *flags.FixedWsModeFlag, false)
	if fixedMode {
		userConfig.WorkspaceMode = config.ConfigWorkspaceModeFixed
	}

	if err := filesystem.WriteConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess("Workspace set to " + workspace)
}

func registerSetNamespaceCmd(ctx *context.Context, setCmd *cobra.Command) {
	namespaceCmd := &cobra.Command{
		Use:     "namespace NAME",
		Aliases: []string{"ns"},
		Short:   "Change the current namespace.",
		Args:    cobra.ExactArgs(1),
		PreRun:  func(cmd *cobra.Command, args []string) { printContext(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { setNamespaceFunc(ctx, cmd, args) },
	}
	setCmd.AddCommand(namespaceCmd)
}

func setNamespaceFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	namespace := args[0]
	userConfig := ctx.Config
	userConfig.CurrentNamespace = namespace
	if err := filesystem.WriteConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess("Namespace set to " + namespace)
}

func registerSetWorkspaceModeCmd(ctx *context.Context, setCmd *cobra.Command) {
	workspaceModeCmd := &cobra.Command{
		Use:       "workspace-mode [fixed|dynamic]",
		Short:     "Switch between fixed and dynamic workspace modes.",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"fixed", "dynamic"},
		PreRun:    func(cmd *cobra.Command, args []string) { printContext(ctx, cmd) },
		Run:       func(cmd *cobra.Command, args []string) { setWorkspaceModeFunc(ctx, cmd, args) },
	}
	setCmd.AddCommand(workspaceModeCmd)
}

func setWorkspaceModeFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	mode := config.ConfigWorkspaceMode(strings.ToLower(args[0]))

	userConfig := ctx.Config
	if userConfig.Interactive == nil {
		userConfig.Interactive = &config.Interactive{}
	}
	userConfig.WorkspaceMode = mode
	if err := filesystem.WriteConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess(fmt.Sprintf("Workspace mode set to '%s'", string(mode)))
}

func registerSetLogModeCmd(ctx *context.Context, setCmd *cobra.Command) {
	logModeCmd := &cobra.Command{
		Use:       "log-mode [logfmt|json|text|hidden]",
		Short:     "Set the default log mode.",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"logfmt", "json", "text", "hidden"},
		PreRun:    func(cmd *cobra.Command, args []string) { printContext(ctx, cmd) },
		Run:       func(cmd *cobra.Command, args []string) { setLogModeFunc(ctx, cmd, args) },
	}
	setCmd.AddCommand(logModeCmd)
}

func setLogModeFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	mode := io2.LogMode(strings.ToLower(args[0]))

	userConfig := ctx.Config
	userConfig.DefaultLogMode = mode
	if err := filesystem.WriteConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess(fmt.Sprintf("Default log mode set to '%s'", mode))
}

func registerSetTUICmd(ctx *context.Context, setCmd *cobra.Command) {
	tuiCmd := &cobra.Command{
		Use:       "tui [true|false]",
		Short:     "Enable or disable the interactive terminal UI experience.",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"true", "false"},
		Run:       func(cmd *cobra.Command, args []string) { setInteractiveFunc(ctx, cmd, args) },
	}
	setCmd.AddCommand(tuiCmd)
}

func setInteractiveFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	enabled, err := strconv.ParseBool(args[0])
	if err != nil {
		logger.FatalErr(errors.Wrap(err, "invalid boolean value"))
	}

	userConfig := ctx.Config
	if userConfig.Interactive == nil {
		userConfig.Interactive = &config.Interactive{}
	}
	userConfig.Interactive.Enabled = enabled
	if err := filesystem.WriteConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}
	strVal := "disabled"
	if enabled {
		strVal = "enabled"
	}
	logger.PlainTextSuccess("Interactive UI " + strVal)
}

func registerViewConfigCmd(ctx *context.Context, configCmd *cobra.Command) {
	viewCmd := &cobra.Command{
		Use:     "view",
		Aliases: []string{"show", "current"},
		Short:   "View the current global configuration values.",
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { StartTUI(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { WaitForTUI(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { viewConfigFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, viewCmd, *flags.OutputFormatFlag)
	configCmd.AddCommand(viewCmd)
}

func viewConfigFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	userConfig := ctx.Config
	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)
	if TUIEnabled(ctx, cmd) {
		view := config2.NewUserConfigView(ctx.TUIContainer, *userConfig, types.Format(outputFormat))
		SetView(ctx, cmd, view)
	} else {
		config2.PrintUserConfig(logger, outputFormat, userConfig)
	}
}
