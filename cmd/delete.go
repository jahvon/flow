package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/io"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a configuration, environment, or workspace option.",
}

// deleteWorkspaceCmd represents the delete workspace subcommand
var deleteWorkspaceCmd = &cobra.Command{
	Use:   "workspace <name>",
	Short: "Delete an existing workspace.",
	Long: "Delete an existing workspace. File contents will remain in the corresponding directory but the workspace will be " +
		"unlinked from flow's conv. Note: You cannot delete the current workspace.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a name argument")
		}

		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}

		name := args[0]
		if name == rootCfg.CurrentWorkspace {
			io.PrintErrorAndExit(fmt.Errorf("cannot delete the current workspace"))
		}
		if _, found := rootCfg.Workspaces[name]; !found {
			io.PrintErrorAndExit(fmt.Errorf("workspace %s was not found", name))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		confirmed := io.AskYesNo("Are you sure you want to delete the workspace '" + name + "'?")
		if !confirmed {
			io.PrintInfo("Aborting")
			return
		}

		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}

		if err := config.DeleteWorkspace(rootCfg, name); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintWarning(fmt.Sprintf("Workspace %s deleted", name))
	},
}

func init() {
	deleteCmd.AddCommand(deleteWorkspaceCmd)
	rootCmd.AddCommand(deleteCmd)
}
