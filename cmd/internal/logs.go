package internal

import (
	"fmt"

	tuikitIO "github.com/jahvon/tuikit/io"
	"github.com/jahvon/tuikit/views"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io/logs"
)

func RegisterLogsCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:     "logs",
		Aliases: []string{"log"},
		Short:   "View execution history and logs.",
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { StartTUI(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { WaitForTUI(ctx, cmd) },
		Run: func(cmd *cobra.Command, args []string) {
			logFunc(ctx, cmd, args)
		},
	}
	RegisterFlag(ctx, subCmd, *flags.LastLogEntryFlag)
	RegisterFlag(ctx, subCmd, *flags.OutputFormatFlag)
	rootCmd.AddCommand(subCmd)
}

func logFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	lastEntry := flags.ValueFor[bool](ctx, cmd, *flags.LastLogEntryFlag, false)
	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)
	if err := filesystem.EnsureLogsDir(); err != nil {
		ctx.Logger.FatalErr(err)
	}
	if TUIEnabled(ctx, cmd) {
		view := views.NewLogArchiveView(ctx.TUIContainer.RenderState(), filesystem.LogsDir(), lastEntry)
		SetView(ctx, cmd, view)
		return
	}
	entries, err := tuikitIO.ListArchiveEntries(filesystem.LogsDir())
	if err != nil {
		ctx.Logger.FatalErr(err)
	}

	if lastEntry {
		if len(entries) == 0 {
			ctx.Logger.Fatalf("No log entries found")
		}
		data, err := entries[0].Read()
		if err != nil {
			ctx.Logger.FatalErr(err)
		}
		_, _ = fmt.Fprint(ctx.StdOut(), data)
	} else {
		logs.PrintEntries(ctx, outputFormat, entries)
	}
}
