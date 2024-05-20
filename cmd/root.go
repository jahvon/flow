// Package cmd handle the cli commands
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal"
	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/cmd/internal/version"
	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/internal/context"
)

func NewRootCmd(ctx *context.Context) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "flow",
		Short: "flow is a command line interface designed to make managing and running development workflows easier.",
		Long: "flow is a command line interface designed to make managing and running development workflows easier." +
			"It's driven by executables organized across workspaces and namespaces defined in a workspace.\n\n" +
			"See github.com/jahvon/flow for more information.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			verbosity := flags.ValueFor[int](ctx, cmd.Root(), *flags.VerbosityFlag, true)
			ctx.Logger.SetLevel(verbosity)
			sync := flags.ValueFor[bool](ctx, cmd.Root(), *flags.SyncCacheFlag, true)
			if sync {
				if err := cache.UpdateAll(ctx.Logger); err != nil {
					ctx.Logger.FatalErr(err)
				}
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) { ctx.Finalize() },
		Version:           version.String(),
	}
	internal.RegisterPersistentFlag(ctx, rootCmd, *flags.VerbosityFlag)
	internal.RegisterPersistentFlag(ctx, rootCmd, *flags.SyncCacheFlag)
	internal.RegisterPersistentFlag(ctx, rootCmd, *flags.NonInteractiveFlag)
	return rootCmd
}

func Execute(ctx *context.Context, rootCmd *cobra.Command) error {
	if ctx == nil {
		panic("current context is not initialized")
	} else if rootCmd == nil {
		panic("root command is not initialized")
	}

	rootCmd.SetOut(ctx.StdOut())
	rootCmd.SetErr(ctx.StdOut())
	rootCmd.SetIn(ctx.StdIn())
	RegisterSubCommands(ctx, rootCmd)

	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	return nil
}

func RegisterSubCommands(ctx *context.Context, rootCmd *cobra.Command) {
	if ctx == nil {
		panic("current context is not initialized")
	} else if rootCmd == nil {
		panic("root command is not initialized")
	}

	internal.RegisterExecCmd(ctx, rootCmd)
	internal.RegisterGetCmd(ctx, rootCmd)
	internal.RegisterInitCmd(ctx, rootCmd)
	internal.RegisterLibraryCmd(ctx, rootCmd)
	internal.RegisterListCmd(ctx, rootCmd)
	internal.RegisterLogsCmd(ctx, rootCmd)
	internal.RegisterRemoveCmd(ctx, rootCmd)
	internal.RegisterSetCmd(ctx, rootCmd)
	internal.RegisterSyncCmd(ctx, rootCmd)
}
