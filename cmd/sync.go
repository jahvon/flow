package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config/cache"

	"github.com/jahvon/flow/internal/io"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Scan workspaces and update flow cache.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cache.UpdateAll(); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintInfo("Synced flow cache")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
