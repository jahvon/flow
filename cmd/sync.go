package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config/cache"
)

var syncCmd = &cobra.Command{
	Use:    "sync",
	Short:  "Scan workspaces and update flow cache.",
	PreRun: initInteractiveCommand,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		if err := cache.UpdateAll(curCtx.Logger); err != nil {
			logger.FatalErr(err)
		}
		logger.PlainTextSuccess("Synced flow cache")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
