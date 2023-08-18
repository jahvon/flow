package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/tbox/internal/config"
	"github.com/jahvon/tbox/internal/io"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a configuration, environment, or workspace option.",
}

// setWorkspaceCmd represents the set workspace subcommand
var createWorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Create a new workspace.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a name argument")
		}

		name := args[0]
		if _, found := currentConfig.Workspaces[name]; found {
			io.PrintErrorAndExit(fmt.Errorf("workspace %s already exists at %s", name, currentConfig.Workspaces[name]))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		path := cmd.Flag("path").Value.String()
		if err := config.AddWorkspace(currentConfig, name, path); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess(fmt.Sprintf("Workspace %s created in %s", name, path))

		if cmd.Flag("set").Value.String() == "true" {
			if err := config.SetCurrentWorkspace(currentConfig, name); err != nil {
				io.PrintErrorAndExit(err)
			}
			io.PrintInfo(fmt.Sprintf("Workspace %s set as current workspace", name))
		}

	},
}

func init() {
	createWorkspaceCmd.Flags().StringP("path", "p", config.DataDirPath(), "Path to the directory where the workspace should be created")
	createWorkspaceCmd.Flags().BoolP("set", "s", false, "Set the newly created workspace as the current workspace")
	createCmd.AddCommand(createWorkspaceCmd)
	rootCmd.AddCommand(createCmd)
}
