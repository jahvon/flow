package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/crypto"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/vault"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize or restore the flow application state.",
}

var configInitCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Initialize the flow global configuration.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		resp := io.AskYesNo("This will overwrite your current flow configurations. Are you sure you want to continue?")
		if !resp {
			io.PrintWarning("Aborting")
			return
		}

		if err := file.InitUserConfig(); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess("Initialized flow global configurations")
	},
}

var workspaceInitCmd = &cobra.Command{
	Use:     "workspace <name> <path>",
	Aliases: []string{"ws"},
	Short:   "Initialize and add a workspace to the list of known workspaces.",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		path := args[1]

		userConfig := file.LoadUserConfig()
		if userConfig == nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to load user config"))
		}
		if err := userConfig.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		if _, found := userConfig.Workspaces[name]; found {
			io.PrintErrorAndExit(fmt.Errorf("workspace %s already exists at %s", name, userConfig.Workspaces[name]))
		}

		if path == "" {
			path = filepath.Join(file.CachedDataDirPath(), name)
		} else if path == "." || strings.HasPrefix(path, "./") {
			wd, err := os.Getwd()
			if err != nil {
				io.PrintErrorAndExit(err)
			}
			if path == "." {
				path = wd
			} else {
				path = fmt.Sprintf("%s/%s", wd, path[2:])
			}
		} else if path == "~" || strings.HasPrefix(path, "~/") {
			hd, err := os.UserHomeDir()
			if err != nil {
				io.PrintErrorAndExit(err)
			}
			if path == "~" {
				path = hd
			} else {
				path = fmt.Sprintf("%s/%s", hd, path[2:])
			}
		}

		if err := file.InitWorkspaceConfig(name, path); err != nil {
			io.PrintErrorAndExit(err)
		}
		userConfig.Workspaces[name] = path

		set := getFlagValue[bool](cmd, *flags.SetAfterCreateFlag)
		if set {
			userConfig.CurrentWorkspace = name
			io.PrintInfo(fmt.Sprintf("Workspace '%s' set as current workspace", name))
		}

		if err := file.WriteUserConfig(userConfig); err != nil {
			io.PrintErrorAndExit(err)
		}

		if err := cache.UpdateAll(); err != nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to update cache - %w", err))
		}

		io.PrintSuccess(fmt.Sprintf("Workspace '%s' created in %s", name, path))
	},
}

var vaultInitCmd = &cobra.Command{
	Use:   "vault",
	Short: "Create a new flow secret vault.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		generatedKey, err := crypto.GenerateKey()
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		if err = vault.RegisterEncryptionKey(generatedKey); err != nil {
			io.PrintErrorAndExit(err)
		}

		io.PrintSuccess(fmt.Sprintf("Your vault encryption key is: %s", generatedKey))
		newKeyMsg := fmt.Sprintf(
			"You will need this key to modify your vault data. Store it somewhere safe!\n"+
				"Set this value to the %s environment variable if you do not want to be prompted for it every time.",
			vault.EncryptionKeyEnvVar,
		)
		io.PrintNotice(newKeyMsg)
	},
}

func init() {
	initCmd.AddCommand(configInitCmd)

	registerFlagOrPanic(workspaceInitCmd, *flags.SetAfterCreateFlag)
	initCmd.AddCommand(workspaceInitCmd)

	initCmd.AddCommand(vaultInitCmd)

	rootCmd.AddCommand(initCmd)
}
