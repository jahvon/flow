package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/io"
	configio "github.com/jahvon/flow/internal/io/config"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"c", "cfg"},
	Short:   "Manage flow user configurations.",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the flow user configuration.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		resp := io.AskYesNo("This will overwrite your current flow user configuration. Are you sure you want to continue?")
		if !resp {
			io.PrintWarning("Aborting")
			return
		}

		if err := file.InitUserConfig(); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess("Initialized flow user configuration")
	},
}

var configSetCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s"},
	Short:   "Update user configuration values.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		userConfig := file.LoadUserConfig()
		if userConfig == nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to load user config"))
		}
		if err := userConfig.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		ws := getFlagValue[string](cmd, *flags.SetWorkspaceFlag)
		ns := getFlagValue[string](cmd, *flags.SetNamespaceFlag)
		uiEnabled := getFlagValue[bool](cmd, *flags.SetUIEnabledFlag)
		uiEnabledChanged := uiEnabled != userConfig.UIEnabled
		if ws == "" && ns == "" && !uiEnabledChanged {
			io.PrintErrorAndExit(fmt.Errorf("no flags provided"))
		}

		if ws != "" {
			if _, found := userConfig.Workspaces[ws]; !found {
				io.PrintErrorAndExit(fmt.Errorf("workspace %s not found", ws))
			}
			userConfig.CurrentWorkspace = ws
		}

		if ns != "" {
			userConfig.CurrentNamespace = ns
		}

		if uiEnabledChanged {
			userConfig.UIEnabled = uiEnabled
		}

		if err := file.WriteUserConfig(userConfig); err != nil {
			io.PrintErrorAndExit(err)
		}

		io.PrintSuccess("Updated flow user configuration")
	},
}

var configGetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	Short:   "Get the current user configurations.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		userConfig := file.LoadUserConfig()
		if userConfig == nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to load user config"))
		}
		if err := userConfig.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		outputFormat := getFlagValue[string](cmd, *flags.OutputFormatFlag)
		configio.PrintUserConfig(io.OutputFormat(outputFormat), userConfig)
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)

	registerFlagOrPanic(configSetCmd, *flags.SetWorkspaceFlag)
	registerFlagOrPanic(configSetCmd, *flags.SetNamespaceFlag)
	registerFlagOrPanic(configSetCmd, *flags.SetUIEnabledFlag)
	configCmd.AddCommand(configSetCmd)

	registerFlagOrPanic(configGetCmd, *flags.OutputFormatFlag)
	configCmd.AddCommand(configGetCmd)

	rootCmd.AddCommand(configCmd)
}
