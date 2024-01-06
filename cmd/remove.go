package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/vault"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove a flow object.",
}

var workspaceRemoveCmd = &cobra.Command{
	Use:     "workspace <name>",
	Aliases: []string{"ws"},
	Short:   "Remove an existing workspace from the list of known workspaces.",
	Long: "Remove an existing workspace. File contents will remain in the corresponding directory but the " +
		"workspace will be unlinked from the flow global configurations.\nNote: You cannot remove the current workspace.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		confirmed := io.AskYesNo("Are you sure you want to remove the workspace '" + name + "'?")
		if !confirmed {
			io.PrintWarning("Aborting")
			return
		}

		userConfig := file.LoadUserConfig()
		if userConfig == nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to load user config"))
		}
		if err := userConfig.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		if name == userConfig.CurrentWorkspace {
			io.PrintErrorAndExit(fmt.Errorf("cannot remove the current workspace"))
		}
		if _, found := userConfig.Workspaces[name]; !found {
			io.PrintErrorAndExit(fmt.Errorf("workspace %s was not found", name))
		}

		delete(userConfig.Workspaces, name)
		if err := file.WriteUserConfig(userConfig); err != nil {
			io.PrintErrorAndExit(err)
		}

		io.PrintWarning(fmt.Sprintf("Workspace '%s' removed", name))

		if err := cache.UpdateAll(); err != nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to update cache - %w", err))
		}
	},
}

var vaultSecretRemoveCmd = &cobra.Command{
	Use:     "secret <name>",
	Aliases: []string{"scrt"},
	Short:   "Remove a secret from the vault.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reference := args[0]

		v := vault.NewVault()
		err := v.DeleteSecret(reference)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess(fmt.Sprintf("Secret %s removed from vault", reference))
	},
}

func init() {
	removeCmd.AddCommand(workspaceRemoveCmd)
	removeCmd.AddCommand(vaultSecretRemoveCmd)

	rootCmd.AddCommand(removeCmd)
}
