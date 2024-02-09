package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config/cache"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Scan workspaces and update flow cache.",
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		if err := cache.UpdateAll(); err != nil {
			logger.FatalErr(err)
		}
		logger.PlainTextSuccess("Synced flow cache")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
