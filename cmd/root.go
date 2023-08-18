// Package cmd handle the cli commands
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/tbox/internal/cmd/version"
	"github.com/jahvon/tbox/internal/io"
)

var rootCmd = &cobra.Command{
	Use:     "tbox",
	Short:   "[Alpha] CLI script wrapper",
	Long:    `Command line interface script wrapper`,
	Version: version.String(),
}

var log = io.Log()

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to execute the root command")
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose log output")
}
