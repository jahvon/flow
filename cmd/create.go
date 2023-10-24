package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/common"
	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/services/cache"
)

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	GroupID: DataGroup.ID,
	Short:   "Create a new set of flow configurations.",
}

// createWorkspaceCmd represents the create workspace subcommand.
var createWorkspaceCmd = &cobra.Command{
	Use:     "workspace <name>",
	Aliases: []string{"w"},
	Short:   "Create a new workspace.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a name argument")
		}

		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Panic().Msg("failed to load config")
		}

		name := args[0]
		if _, found := rootCfg.Workspaces[name]; found {
			io.PrintErrorAndExit(fmt.Errorf("workspace %s already exists at %s", name, rootCfg.Workspaces[name]))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Panic().Msg("failed to load config")
		}

		path := cmd.Flag("path").Value.String()
		if path == "" {
			path = common.ConfigDirPath()
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

		if err := config.CreateWorkspace(rootCfg, name, path); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess(fmt.Sprintf("Workspace %s created in %s", name, path))

		set, err := cmd.Flags().GetBool("set")
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		if set {
			if err := config.SetCurrentWorkspace(rootCfg, name); err != nil {
				io.PrintErrorAndExit(err)
			}
			io.PrintInfo(fmt.Sprintf("Workspace %s set as current workspace", name))
		}

		if _, err := cache.Update(); err != nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to update cache - %w", err))
		}
	},
}

func init() {
	registerFlagOrPanic(createWorkspaceCmd, *flags.SetAfterCreateFlag)
	registerFlagOrPanic(createWorkspaceCmd, *flags.WorkspacePathFlag)
	if err := createWorkspaceCmd.MarkFlagDirname(flags.WorkspacePathFlag.Name); err != nil {
		log.Panic().Err(err).Msg("Failed to mark path flag as a directory")
	}
	createCmd.AddCommand(createWorkspaceCmd)

	rootCmd.AddCommand(createCmd)
}
