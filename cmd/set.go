package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/vault"
)

var setCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s"},
	Short:   "Update global or workspace configuration values.",
}

var configWorkspaceSetCmd = &cobra.Command{
	Use:     "workspace <name>",
	Aliases: []string{"ws"},
	Short:   "Change the current workspace.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workspace := args[0]
		userConfig := file.LoadUserConfig()
		if userConfig == nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to load user config"))
		}
		if err := userConfig.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		if _, found := userConfig.Workspaces[workspace]; !found {
			io.PrintErrorAndExit(fmt.Errorf("workspace %s not found", workspace))
		}
		userConfig.CurrentWorkspace = workspace

		if err := file.WriteUserConfig(userConfig); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess("Workspace set to " + workspace)
	},
}

var configNamespaceSetCmd = &cobra.Command{
	Use:     "namespace <name>",
	Aliases: []string{"ns"},
	Short:   "Change the current namespace.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := args[0]
		userConfig := file.LoadUserConfig()
		if userConfig == nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to load user config"))
		}
		if err := userConfig.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		userConfig.CurrentNamespace = namespace
		if err := file.WriteUserConfig(userConfig); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess("Namespace set to " + namespace)
	},
}

var configInteractiveSetCmd = &cobra.Command{
	Use:   "interactive <true|false>",
	Short: "Enable or disable the interactive terminal UI experience.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		enabled, err := strconv.ParseBool(args[0])
		if err != nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to parse boolean value: %v", err))
		}

		userConfig := file.LoadUserConfig()
		if userConfig == nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to load user config"))
		}
		if err := userConfig.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		userConfig.InteractiveUI = enabled
		if err := file.WriteUserConfig(userConfig); err != nil {
			io.PrintErrorAndExit(err)
		}
		strVal := "disabled"
		if enabled {
			strVal = "enabled"
		}
		io.PrintSuccess("Interactive UI " + strVal)
	},
}

var vaultSecretSetCmd = &cobra.Command{
	Use:     "secret <name> <value>",
	Aliases: []string{"scrt"},
	Short:   "Update or create a secret in the flow secret vault.",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		reference := args[0]
		value := args[1]

		secret := vault.Secret(value)
		v := vault.NewVault()
		err := v.SetSecret(reference, secret)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess(fmt.Sprintf("Secret %s set in vault", reference))
	},
}

func init() {
	setCmd.AddCommand(configWorkspaceSetCmd)
	setCmd.AddCommand(configNamespaceSetCmd)
	setCmd.AddCommand(configInteractiveSetCmd)
	setCmd.AddCommand(vaultSecretSetCmd)

	rootCmd.AddCommand(setCmd)
}
