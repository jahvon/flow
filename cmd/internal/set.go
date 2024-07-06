package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jahvon/tuikit/components"
	tuiKitIO "github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/cmd/internal/interactive"
	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/vault"
)

func RegisterSetCmd(ctx *context.Context, rootCmd *cobra.Command) {
	setCmd := &cobra.Command{
		Use:     "set",
		Aliases: []string{"s"},
		Short:   "Update global or workspace configuration values.",
	}
	registerSetWorkspaceCmd(ctx, setCmd)
	registerSetNamespaceCmd(ctx, setCmd)
	registerSetWorkspaceModeCmd(ctx, setCmd)
	registerSetLogModeCmd(ctx, setCmd)
	registerSetInteractiveCmd(ctx, setCmd)
	registerSetTemplateCmd(ctx, setCmd)
	registerSetSecretCmd(ctx, setCmd)
	rootCmd.AddCommand(setCmd)
}

func registerSetWorkspaceCmd(ctx *context.Context, setCmd *cobra.Command) {
	workspaceCmd := &cobra.Command{
		Use:     "workspace NAME",
		Aliases: []string{"ws"},
		Short:   "Change the current workspace.",
		Args:    cobra.ExactArgs(1),
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { setWorkspaceFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, workspaceCmd, *flags.FixedWsModeFlag)
	setCmd.AddCommand(workspaceCmd)
}

func setWorkspaceFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	workspace := args[0]
	userConfig := ctx.UserConfig
	if _, found := userConfig.Workspaces[workspace]; !found {
		logger.Fatalf("workspace %s not found", workspace)
	}
	userConfig.CurrentWorkspace = workspace
	fixedMode := flags.ValueFor[bool](ctx, cmd, *flags.FixedWsModeFlag, false)
	if fixedMode {
		userConfig.WorkspaceMode = config.WorkspaceModeFixed
	}

	if err := file.WriteUserConfig(userConfig); err != nil {
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
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { setNamespaceFunc(ctx, cmd, args) },
	}
	setCmd.AddCommand(namespaceCmd)
}

func setNamespaceFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	namespace := args[0]
	userConfig := ctx.UserConfig
	userConfig.CurrentNamespace = namespace
	if err := file.WriteUserConfig(userConfig); err != nil {
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
		PreRun:    func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:       func(cmd *cobra.Command, args []string) { setWorkspaceModeFunc(ctx, cmd, args) },
	}
	setCmd.AddCommand(workspaceModeCmd)
}

func setWorkspaceModeFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	mode := config.WorkspaceMode(strings.ToLower(args[0]))

	userConfig := ctx.UserConfig
	if userConfig.Interactive == nil {
		userConfig.Interactive = &config.InteractiveConfig{}
	}
	userConfig.WorkspaceMode = mode
	if err := file.WriteUserConfig(userConfig); err != nil {
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
		PreRun:    func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:       func(cmd *cobra.Command, args []string) { setLogModeFunc(ctx, cmd, args) },
	}
	setCmd.AddCommand(logModeCmd)
}

func setLogModeFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	mode := tuiKitIO.LogMode(strings.ToLower(args[0]))

	userConfig := ctx.UserConfig
	userConfig.DefaultLogMode = mode
	if err := file.WriteUserConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess(fmt.Sprintf("Default log mode set to '%s'", mode))
}

func registerSetInteractiveCmd(ctx *context.Context, setCmd *cobra.Command) {
	interactiveCmd := &cobra.Command{
		Use:       "interactive [true|false]",
		Short:     "Enable or disable the interactive terminal UI experience.",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"true", "false"},
		PreRun:    func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:       func(cmd *cobra.Command, args []string) { setInteractiveFunc(ctx, cmd, args) },
	}
	setCmd.AddCommand(interactiveCmd)
}

func setInteractiveFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	enabled, err := strconv.ParseBool(args[0])
	if err != nil {
		logger.FatalErr(errors.Wrap(err, "invalid boolean value"))
	}

	userConfig := ctx.UserConfig
	if userConfig.Interactive == nil {
		userConfig.Interactive = &config.InteractiveConfig{}
	}
	userConfig.Interactive.Enabled = enabled
	if err := file.WriteUserConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}
	strVal := "disabled"
	if enabled {
		strVal = "enabled"
	}
	logger.PlainTextSuccess("Interactive UI " + strVal)
}

func registerSetTemplateCmd(ctx *context.Context, setCmd *cobra.Command) {
	templateCmd := &cobra.Command{
		Use:    "template NAME DEFINITION_TEMPLATE_PATH",
		Short:  "Set a template definition for use in flow.",
		Args:   cobra.ExactArgs(2),
		PreRun: func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:    func(cmd *cobra.Command, args []string) { setTemplateFunc(ctx, cmd, args) },
	}
	setCmd.AddCommand(templateCmd)
}

func setTemplateFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	name := args[0]
	definitionPath := args[1]
	loadedTemplates, err := file.LoadExecutableDefinitionTemplate(definitionPath)
	if err != nil {
		logger.FatalErr(err)
	}
	if err := loadedTemplates.Validate(); err != nil {
		logger.FatalErr(err)
	}
	userConfig := ctx.UserConfig
	if userConfig.Templates == nil {
		userConfig.Templates = map[string]string{}
	}
	userConfig.Templates[name] = definitionPath
	if err := file.WriteUserConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess(fmt.Sprintf("Template %s set to %s", name, definitionPath))
}

func registerSetSecretCmd(ctx *context.Context, setCmd *cobra.Command) {
	secretCmd := &cobra.Command{
		Use:     "secret NAME [VALUE]",
		Aliases: []string{"scrt"},
		Short:   "Update or create a secret in the flow secret vault.",
		Args:    cobra.ExactArgs(2),
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { setSecretFunc(ctx, cmd, args) },
	}
	setCmd.AddCommand(secretCmd)
}

func setSecretFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	reference := args[0]
	value := args[1]

	if value == "" {
		in := components.TextInput{Key: "value", Prompt: "Enter the secret value"}
		inputs, err := components.ProcessInputs(io.Theme(), &in)
		if err != nil {
			logger.FatalErr(err)
		}
		value = inputs.FindByKey("value").Value()
	}

	secret := vault.SecretValue(value)
	v := vault.NewVault(logger)
	err := v.SetSecret(reference, secret)
	if err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess(fmt.Sprintf("Secret %s set in vault", reference))
}
