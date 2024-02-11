package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/jahvon/tuikit/components"
	tuikitIO "github.com/jahvon/tuikit/io"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/io"
)

var logsCmd = &cobra.Command{
	Use:     "logs",
	Aliases: []string{"log"},
	Short:   "List and view logs for previous flow executions.",
	Args:    cobra.NoArgs,
	PreRun:  initInteractiveContainer,
	PostRun: waitForExit,
	Run: func(cmd *cobra.Command, args []string) {
		lastEntry := getFlagValue[bool](cmd, *flags.LastLogEntryFlag)
		if err := file.EnsureLogsDir(); err != nil {
			curCtx.Logger.FatalErr(err)
		}
		if interactiveUIEnabled() {
			state := &components.TerminalState{
				Theme:  io.Styles(),
				Height: curCtx.InteractiveContainer.Height(),
				Width:  curCtx.InteractiveContainer.Width(),
			}
			view := components.NewLogArchiveView(state, file.LogsDirPath, lastEntry)
			curCtx.InteractiveContainer.SetView(view)
			return
		}
		entries, err := tuikitIO.ListArchiveEntries(file.LogsDirPath)
		if err != nil {
			curCtx.Logger.FatalErr(err)
		} else if len(entries) == 0 {
			curCtx.Logger.PlainTextInfo("No logs entries found")
		}
		if lastEntry {
			data, err := entries[0].Read()
			if err != nil {
				curCtx.Logger.FatalErr(err)
			}
			_, _ = fmt.Fprint(os.Stdout, data)
		} else {
			for _, entry := range entries {
				entryStr := fmt.Sprintf(
					"%s (%s)\n%s\n\n",
					entry.Args,
					entry.Time.Local().Format(time.RFC822),
					entry.Path,
				)
				_, _ = fmt.Fprint(os.Stdout, entryStr)
			}
		}
	},
}

func init() {
	registerFlagOrPanic(logsCmd, *flags.LastLogEntryFlag)
	rootCmd.AddCommand(logsCmd)
}
