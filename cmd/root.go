// Package cmd handle the cli commands
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/cmd/version"
	"github.com/jahvon/flow/internal/context"
)

var rootCmd = &cobra.Command{
	Use:   "flow",
	Short: "flow is a command line interface designed to make managing and running development workflows easier.",
	Long: "flow is a command line interface designed to make managing and running development workflows easier." +
		"It's driven by executables organized across workspaces and namespaces defined in a workspace.\n\n" +
		"See github.com/jahvon/flow for more information.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbosity := getPersistentFlagValue[int](cmd.Root(), *flags.VerbosityFlag)
		curCtx.Logger.SetLevel(verbosity)
		if verbosity != 0 {
			curCtx.Logger.Infof("Log level set to %d", verbosity)
		}

		sync := getPersistentFlagValue[bool](cmd.Root(), *flags.SyncCacheFlag)
		if sync {
			if err := cache.UpdateAll(); err != nil {
				curCtx.Logger.FatalErr(err)
			}
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		curCtx.Finalize()
	},
	Version: version.String(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx *context.Context) {
	curCtx = ctx
	if curCtx == nil {
		panic("current context is not initialized")
	}

	err := rootCmd.Execute()
	if err != nil {
		curCtx.Logger.FatalErr(fmt.Errorf("failed to execute command: %w", err))
	}
}

func init() {
	registerPersistentFlagOrPanic(rootCmd, *flags.VerbosityFlag)
	registerPersistentFlagOrPanic(rootCmd, *flags.SyncCacheFlag)
	registerPersistentFlagOrPanic(rootCmd, *flags.NonInteractiveFlag)
}
