package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/services/cache"
	"github.com/jahvon/flow/internal/services/git"
	"github.com/jahvon/flow/internal/workspace"
)

var syncCmd = &cobra.Command{
	Use:     "sync",
	GroupID: DataGroup.ID,
	Short:   "Sync flow cache and workspaces.",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := cache.Update()
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		log.Info().Int32("workspaces", int32(len(data.Workspaces))).
			Msg("Successfully synced flow cache.")

		// Pull git workspaces
		for wsName, ws := range data.Workspaces {
			wsCfg, err := workspace.LoadConfig(wsName, ws.Location())
			if err != nil {
				io.PrintError(err)
				continue
			} else if wsCfg == nil {
				io.PrintError(fmt.Errorf("config not found for workspace %s", ws.AssignedName()))
				continue
			}
			if wsCfg.Git != nil && wsCfg.Git.Enabled && wsCfg.Git.PullOnSync {
				if err := git.Pull(ws.Location()); err != nil {
					io.PrintError(err)
					continue
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
