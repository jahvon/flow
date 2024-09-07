package internal

import (
	"os"
	"strconv"

	"github.com/jahvon/tuikit/components"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/types/executable"
	"github.com/jahvon/flow/types/workspace"
)

func RegisterFlag(ctx *context.Context, cmd *cobra.Command, flag flags.Metadata) {
	flagSet, err := flags.ToPflag(cmd, flag, false)
	if err != nil {
		ctx.Logger.FatalErr(err)
	}
	cmd.Flags().AddFlagSet(flagSet)
	if flag.Required {
		MarkFlagRequired(ctx, cmd, flag.Name)
	}
}

func RegisterPersistentFlag(ctx *context.Context, cmd *cobra.Command, flag flags.Metadata) {
	flagSet, err := flags.ToPflag(cmd, flag, true)
	if err != nil {
		ctx.Logger.FatalErr(err)
	}
	cmd.PersistentFlags().AddFlagSet(flagSet)
}

func MarkFlagRequired(ctx *context.Context, cmd *cobra.Command, name string) {
	if err := cmd.MarkFlagRequired(name); err != nil {
		ctx.Logger.FatalErr(err)
	}
}

func MarkFlagMutuallyExclusive(cmd *cobra.Command, names ...string) {
	cmd.MarkFlagsMutuallyExclusive(names...)
}

func MarkOneFlagRequired(cmd *cobra.Command, names ...string) {
	cmd.MarkFlagsOneRequired(names...)
}

func MarkFlagFilename(ctx *context.Context, cmd *cobra.Command, name string) {
	if err := cmd.MarkFlagFilename(name); err != nil {
		ctx.Logger.FatalErr(err)
	}
}

func UIEnabled(ctx *context.Context, cmd *cobra.Command) bool {
	disabled := flags.ValueFor[bool](ctx, cmd.Root(), *flags.NonInteractiveFlag, true)
	envDisabled, _ := strconv.ParseBool(os.Getenv("DISABLE_FLOW_INTERACTIVE"))
	return !disabled && !envDisabled && ctx.Config.ShowTUI()
}

func SetView(ctx *context.Context, cmd *cobra.Command, view components.TeaModel) {
	if UIEnabled(ctx, cmd) {
		ctx.SetView(view)
	} else {
		ctx.Logger.Errorx("interactive mode is disabled", "view", view.Type())
	}
}

func SetLoadingView(ctx *context.Context, cmd *cobra.Command) {
	if UIEnabled(ctx, cmd) {
		view := components.NewLoadingView("thinking...", io.Theme())
		SetView(ctx, cmd, view)
	}
}

func printContext(ctx *context.Context, cmd *cobra.Command) {
	if UIEnabled(ctx, cmd) {
		ctx.Logger.Println(io.Theme().RenderHeader(context.AppName, context.HeaderCtxKey, ctx.String(), 0))
	}
}

func workspaceOrCurrent(ctx *context.Context, workspaceName string) *workspace.Workspace {
	var ws *workspace.Workspace
	if workspaceName == "" {
		ws = ctx.CurrentWorkspace
		workspaceName = ws.AssignedName()
	} else {
		wsPath, wsFound := ctx.Config.Workspaces[workspaceName]
		if !wsFound {
			return nil
		}
		var err error
		ws, err = filesystem.LoadWorkspaceConfig(workspaceName, wsPath)
		if err != nil {
			ctx.Logger.Error(err, "unable to load workspace config")
		}
		ws.SetContext(workspaceName, wsPath)
	}
	ctx.Logger.Debugf("'%s' workspace set", workspaceName)
	return ws
}

func loadFlowfileTemplate(ctx *context.Context, name, path string) *executable.Template {
	if name != "" {
		if ctx.Config.Templates == nil {
			ctx.Logger.Errorf("template %s not found", name)
			return nil
		}
		var found bool
		if path, found = ctx.Config.Templates[name]; !found {
			ctx.Logger.Errorf("template %s not found", name)
			return nil
		}
	} else {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			ctx.Logger.Errorf("flowfile template at %s not found", path)
			return nil
		}
	}
	tmpl, err := filesystem.LoadFlowFileTemplate(name, path)
	if err != nil {
		ctx.Logger.Error(err, "unable to load flowfile template")
		return nil
	}
	return tmpl
}
