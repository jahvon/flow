package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/vault"
)

var setCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s"},
	Short:   "Update global or workspace configuration values.",
}

var configWorkspaceSetCmd = &cobra.Command{
	Use:     "workspace NAME",
	Aliases: []string{"ws"},
	Short:   "Change the current workspace.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		workspace := args[0]
		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		userConfig := curCtx.UserConfig
		if _, found := userConfig.Workspaces[workspace]; !found {
			logger.Fatalf("workspace %s not found", workspace)
		}
		userConfig.CurrentWorkspace = workspace

		if err := file.WriteUserConfig(userConfig); err != nil {
			logger.FatalErr(err)
		}
		logger.PlainTextSuccess("Workspace set to " + workspace)
	},
}

var configNamespaceSetCmd = &cobra.Command{
	Use:     "namespace NAME",
	Aliases: []string{"ns"},
	Short:   "Change the current namespace.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		namespace := args[0]
		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		userConfig := curCtx.UserConfig
		userConfig.CurrentNamespace = namespace
		if err := file.WriteUserConfig(userConfig); err != nil {
			logger.FatalErr(err)
		}
		logger.PlainTextSuccess("Namespace set to " + namespace)
	},
}

var configWorkspaceModeSetCmd = &cobra.Command{
	Use:   "workspace-mode (fixed|dynamic)",
	Short: "Switch between fixed and dynamic workspace modes.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		mode := config.WorkspaceMode(strings.ToLower(args[0]))

		userConfig := curCtx.UserConfig
		if userConfig.Interactive == nil {
			userConfig.Interactive = &config.InteractiveConfig{}
		}
		userConfig.WorkspaceMode = mode
		if err := file.WriteUserConfig(userConfig); err != nil {
			logger.FatalErr(err)
		}
		logger.PlainTextSuccess(fmt.Sprintf("Workspace mode set to '%s'", string(mode)))
	},
}

var configInteractiveSetCmd = &cobra.Command{
	Use:   "interactive (true|false)",
	Short: "Enable or disable the interactive terminal UI experience.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		enabled, err := strconv.ParseBool(args[0])
		if err != nil {
			logger.FatalErr(errors.Wrap(err, "invalid boolean value"))
		}

		userConfig := curCtx.UserConfig
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
	},
}

var configTemplateSetCmd = &cobra.Command{
	Use:   "template NAME DEFINITION_TEMPLATE_PATH",
	Short: "Set a template definition for use in flow.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		name := args[0]
		definitionPath := args[1]
		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		loadedTemplates, err := file.LoadExecutableDefinitionTemplate(definitionPath)
		if err != nil {
			logger.FatalErr(err)
		}
		if err := loadedTemplates.Validate(); err != nil {
			logger.FatalErr(err)
		}
		userConfig := curCtx.UserConfig
		if userConfig.Templates == nil {
			userConfig.Templates = map[string]string{}
		}
		userConfig.Templates[name] = definitionPath
		if err := file.WriteUserConfig(userConfig); err != nil {
			logger.FatalErr(err)
		}
		logger.PlainTextSuccess(fmt.Sprintf("Template %s set to %s", name, definitionPath))
	},
}

var configPlaintextLoggerSetCmd = &cobra.Command{
	Use:   "use-plain-text-logger (true|false)",
	Short: "Enable or disable the plain text logger.",
	Long: "Enable or disable the plain text logger. " +
		"When enabled, the log output will include log level and timestamp.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		enabled, err := strconv.ParseBool(args[0])
		if err != nil {
			logger.FatalErr(errors.Wrap(err, "invalid boolean value"))
		}

		userConfig := curCtx.UserConfig
		userConfig.UsePlainTextLogger = enabled
		if err := file.WriteUserConfig(userConfig); err != nil {
			logger.FatalErr(err)
		}
		strVal := "disabled"
		if enabled {
			strVal = "enabled"
		}
		logger.PlainTextSuccess("Plain text logger " + strVal)
	},
}

var vaultSecretSetCmd = &cobra.Command{
	Use:     "secret NAME VALUE",
	Aliases: []string{"scrt"},
	Short:   "Update or create a secret in the flow secret vault.",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		reference := args[0]
		value := args[1]

		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		secret := vault.Secret(value)
		v := vault.NewVault(logger)
		err := v.SetSecret(reference, secret)
		if err != nil {
			logger.FatalErr(err)
		}
		logger.PlainTextSuccess(fmt.Sprintf("Secret %s set in vault", reference))
	},
}

func init() {
	setCmd.AddCommand(configWorkspaceSetCmd)
	setCmd.AddCommand(configNamespaceSetCmd)
	setCmd.AddCommand(configWorkspaceModeSetCmd)
	setCmd.AddCommand(configInteractiveSetCmd)
	setCmd.AddCommand(configTemplateSetCmd)
	setCmd.AddCommand(configPlaintextLoggerSetCmd)
	setCmd.AddCommand(vaultSecretSetCmd)

	rootCmd.AddCommand(setCmd)
}
