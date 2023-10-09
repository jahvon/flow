// Package cmd handle the cli commands
package cmd

import (
	"fmt"
	"strconv"

	"github.com/jahvon/flow/internal/cmd/version"
	"github.com/jahvon/flow/internal/io"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var log = io.Log()

var rootCmd = &cobra.Command{
	Use:   "flow",
	Short: "[Alpha] CLI script wrapper",
	Long:  `Command line interface script wrapper`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbose, err := strconv.ParseBool(cmd.Flag("verbose").Value.String())
		if err != nil {
			io.PrintErrorAndExit(fmt.Errorf("invalid verbose flag - %w", err))
		}

		if verbose {
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
			io.PrintInfo("Verbose logging enabled")
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
	},
	Version: version.String(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		io.PrintErrorAndExit(fmt.Errorf("failed to execute command: %w", err))
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose log output")

	rootCmd.AddGroup(CrudGroup)
	rootCmd.AddGroup(ExecutableGroup)
}
