package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/tbox/internal/cmd/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current tbox version.",
	Run: func(cmd *cobra.Command, args []string) {
		version.Print()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
