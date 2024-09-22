package internal

import (
	"fmt"
	"time"

	tuikitIO "github.com/jahvon/tuikit/io"
	"github.com/jahvon/tuikit/views"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
)

func RegisterLogsCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:     "logs",
		Aliases: []string{"log"},
		Short:   "List and view logs for previous flow executions.",
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { StartTUI(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { WaitForTUI(ctx, cmd) },
		Run: func(cmd *cobra.Command, args []string) {
			logFunc(ctx, cmd, args)
		},
	}
	RegisterFlag(ctx, subCmd, *flags.LastLogEntryFlag)
	rootCmd.AddCommand(subCmd)
}

func logFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	lastEntry := flags.ValueFor[bool](ctx, cmd, *flags.LastLogEntryFlag, false)
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
	} else if len(entries) == 0 {
		ctx.Logger.PlainTextInfo("No logs entries found")
	}
	if lastEntry {
		data, err := entries[0].Read()
		if err != nil {
			ctx.Logger.FatalErr(err)
		}
		_, _ = fmt.Fprint(ctx.StdOut(), data)
	} else {
		for _, entry := range entries {
			entryStr := fmt.Sprintf(
				"%s (%s)\n%s\n\n",
				entry.Args,
				entry.Time.Local().Format(time.RFC822),
				entry.Path,
			)
			_, _ = fmt.Fprint(ctx.StdOut(), entryStr)
		}
	}
}
