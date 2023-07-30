/*
Copyright © 2022 Jahvon Dockery <jahvondockery@gmail.com
*/

package cmd

import (
	"github.com/jahvon/pilotcli/internal/cmd/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the application version.",
	Long:  `Print the application version.`,
	Run: func(cmd *cobra.Command, args []string) {
		version.Print()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
