// Package cmd handle the cli commands
package cmd

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/cmd/version"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/services/cache"
)

var log = io.Log()

var rootCmd = &cobra.Command{
	Use:   "flow",
	Short: "[Alpha] CLI script wrapper",
	Long:  `Command line interface script wrapper`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbosity, err := cmd.Flags().GetInt(flags.VerbosityFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(fmt.Errorf("invalid verbosity flag - %w", err))
		}

		switch verbosity {
		case 0:
			zerolog.SetGlobalLevel(zerolog.NoLevel)
		case 1:
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case 2:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case 3:
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case 4:
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		default:
			io.PrintErrorAndExit(fmt.Errorf("verbosity (%d) must be between 0 and 4", verbosity))
		}

		if verbosity != 2 {
			io.PrintInfo(fmt.Sprintf("Log level set to %d", verbosity))
		}

		syncCache, err := cmd.Flags().GetBool(flags.SyncCacheFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(fmt.Errorf("invalid sync flag - %w", err))
		}
		if syncCache {
			_, err := cache.Update()
			if err != nil {
				io.PrintErrorAndExit(err)
			}
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
	rootCmd.PersistentFlags().IntP(
		flags.VerbosityFlag.Name,
		flags.VerbosityFlag.Shorthand,
		flags.VerbosityFlag.Default.(int),
		flags.VerbosityFlag.Usage,
	)
	rootCmd.PersistentFlags().Bool(
		flags.SyncCacheFlag.Name,
		flags.SyncCacheFlag.Default.(bool),
		flags.SyncCacheFlag.Usage,
	)

	rootCmd.AddGroup(DataGroup)
	rootCmd.AddGroup(ExecutableGroup)
}
