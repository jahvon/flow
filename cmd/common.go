package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jahvon/tuikit/components"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
)

var (
	curCtx *context.Context
)

func interactiveUIEnabled() bool {
	disabled := getPersistentFlagValue[bool](rootCmd, *flags.NonInteractiveFlag)
	envDisabled, _ := strconv.ParseBool(os.Getenv("DISABLE_FLOW_INTERACTIVE"))
	return !disabled && !envDisabled && curCtx.UserConfig.Interactive != nil && curCtx.UserConfig.Interactive.Enabled
}

func headerForCurCtx() components.Header {
	ws := curCtx.UserConfig.CurrentWorkspace
	ns := curCtx.UserConfig.CurrentNamespace
	if ws == "" {
		ws = "unk"
	}
	if ns == "" {
		ns = "*"
	}
	return components.Header{
		Name:   "flow",
		CtxKey: "ctx",
		CtxVal: fmt.Sprintf("%s/%s", ws, ns),
		Styles: io.Styles(),
	}
}

func initInteractiveContainer(_ *cobra.Command, _ []string) {
	enabled := interactiveUIEnabled()
	if enabled && curCtx.InteractiveContainer == nil {
		container := components.InitalizeContainer(curCtx.Ctx, curCtx.CancelFunc, headerForCurCtx(), io.Styles())
		curCtx.InteractiveContainer = container
	}
}

func waitForExit(_ *cobra.Command, _ []string) {
	if interactiveUIEnabled() && curCtx.InteractiveContainer != nil {
		timeout := time.After(30 * time.Minute)
		select {
		case <-curCtx.Ctx.Done():
			return
		case <-timeout:
			panic("interactive wait timeout")
		}
	}
}

func GenerateMarkdownTree(dir string) error {
	rootCmd.DisableAutoGenTag = true
	return doc.GenMarkdownTree(rootCmd, dir)
}
