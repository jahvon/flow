// Package cmd handle the cli commands
package cmd

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/cmd/version"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
)

var (
	log    = io.Log()
	curCtx *context.Context
)

var rootCmd = &cobra.Command{
	Use:   "flow",
	Short: "flow is a command line interface for managing and running machine commands.",
	Long:  `Command line interface script wrapper`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if curCtx == nil {
			log.Panic().Msg("context not initialized")
		}

		verbosityFlag, err := Flags.ValueFor(cmd, flags.VerbosityFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(fmt.Errorf("invalid verbosity flag - %w", err))
		}
		verbosity, _ := verbosityFlag.(int)

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

		syncFlag, err := Flags.ValueFor(cmd, flags.SyncCacheFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		sync, _ := syncFlag.(bool)
		if sync {
			if err := cache.UpdateAll(); err != nil {
				io.PrintErrorAndExit(err)
			}
		}
	},
	Version: version.String(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx *context.Context) {
	curCtx = ctx
	err := rootCmd.Execute()
	if err != nil {
		io.PrintErrorAndExit(fmt.Errorf("failed to execute command: %w", err))
	}
}

func init() {
	registerFlagOrPanic(rootCmd, *flags.VerbosityFlag)
	registerFlagOrPanic(rootCmd, *flags.SyncCacheFlag)
}
