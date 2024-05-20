package interactive

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jahvon/tuikit/components"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
)

const (
	appName      = "flow"
	headerCtxKey = "ctx"
)

func UIEnabled(ctx *context.Context, cmd *cobra.Command) bool {
	disabled := flags.ValueFor[bool](ctx, cmd.Root(), *flags.NonInteractiveFlag, true)
	envDisabled, _ := strconv.ParseBool(os.Getenv("DISABLE_FLOW_INTERACTIVE"))
	return !disabled && !envDisabled && ctx.UserConfig.Interactive != nil && ctx.UserConfig.Interactive.Enabled
}

func InitInteractiveCommand(ctx *context.Context, cmd *cobra.Command) {
	if UIEnabled(ctx, cmd) {
		_, _ = fmt.Fprintln(ctx.StdOut(), io.Theme().
			RenderHeader(appName, headerCtxKey, headerCtxVal(ctx), 0))
	}
}

func InitInteractiveContainer(ctx *context.Context, cmd *cobra.Command) {
	enabled := UIEnabled(ctx, cmd)
	if enabled && ctx.InteractiveContainer == nil {
		container := components.InitalizeContainer(
			ctx.Ctx, ctx.CancelFunc, appName, headerCtxKey, headerCtxVal(ctx), io.Theme(),
		)
		ctx.InteractiveContainer = container
	}
}

func headerCtxVal(ctx *context.Context) string {
	ws := ctx.CurrentWorkspace.AssignedName()
	ns := ctx.UserConfig.CurrentNamespace
	if ws == "" {
		ws = "unk"
	}
	if ns == "" {
		ns = "*"
	}
	return fmt.Sprintf("%s/%s", ws, ns)
}

func WaitForExit(ctx *context.Context, cmd *cobra.Command) {
	if UIEnabled(ctx, cmd) && ctx.InteractiveContainer != nil {
		timeout := time.After(60 * time.Minute)
		select {
		case <-ctx.Ctx.Done():
			return
		case <-timeout:
			panic("interactive wait timeout")
		}
	}
}
