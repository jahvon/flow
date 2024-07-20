package internal

import (
	"fmt"
	"time"

	"github.com/jahvon/tuikit/components"
	tuikitIO "github.com/jahvon/tuikit/io"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/cmd/internal/interactive"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io"
)

func RegisterLogsCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:     "logs",
		Aliases: []string{"log"},
		Short:   "List and view logs for previous flow executions.",
		Args:    cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			interactive.InitInteractiveContainer(ctx, cmd)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			interactive.WaitForExit(ctx, cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			logFunc(ctx, cmd, args)
		},
	}
	RegisterFlag(ctx, subCmd, *flags.CopyFlag)
	rootCmd.AddCommand(subCmd)
}

func logFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	lastEntry := flags.ValueFor[bool](ctx, cmd, *flags.CopyFlag, false)
	if err := filesystem.EnsureLogsDir(); err != nil {
		ctx.Logger.FatalErr(err)
	}
	if interactive.UIEnabled(ctx, cmd) {
		state := &components.TerminalState{
			Theme:  io.Theme(),
			Height: ctx.InteractiveContainer.Height(),
			Width:  ctx.InteractiveContainer.Width(),
		}
		view := components.NewLogArchiveView(state, filesystem.LogsDir(), lastEntry)
		ctx.InteractiveContainer.SetView(view)
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
