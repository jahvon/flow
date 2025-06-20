package workspace

import (
	"fmt"
	"path/filepath"

	"github.com/jahvon/tuikit"
	"github.com/jahvon/tuikit/themes"
	"github.com/jahvon/tuikit/types"
	"github.com/jahvon/tuikit/views"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io/common"
	"github.com/jahvon/flow/internal/services/open"
	"github.com/jahvon/flow/types/workspace"
)

func NewWorkspaceView(
	ctx *context.Context,
	ws *workspace.Workspace,
) tuikit.View {
	container := ctx.TUIContainer
	var workspaceKeyCallbacks = []types.KeyCallback{
		{
			Key: "o", Label: "open",
			Callback: func() error {
				if err := open.Open(ws.Location(), false); err != nil {
					container.HandleError(fmt.Errorf("unable to open workspace: %w", err))
				}
				return nil
			},
		},
		{
			Key: "e", Label: "edit",
			Callback: func() error {
				fullPath := filepath.Join(ws.Location(), filesystem.WorkspaceConfigFileName)
				if err := common.OpenInEditor(fullPath, ctx.StdIn(), ctx.StdOut()); err != nil {
					container.HandleError(fmt.Errorf("unable to edit workspace: %w", err))
				}
				return nil
			},
		},
		{
			Key: "s", Label: "set",
			Callback: func() error {
				curCfg, err := filesystem.LoadConfig()
				if err != nil {
					container.HandleError(err)
					return nil
				}
				curCfg.CurrentWorkspace = ws.AssignedName()
				if err := filesystem.WriteConfig(curCfg); err != nil {
					container.HandleError(err)
				}
				container.SetState(common.HeaderContextKey, fmt.Sprintf("%s/*", ws.AssignedName()))
				container.SetNotice("workspace updated", themes.OutputLevelInfo)
				return nil
			},
		},
	}

	return views.NewEntityView(container.RenderState(), ws, types.EntityFormatDocument, workspaceKeyCallbacks...)
}

func NewWorkspaceListView(
	ctx *context.Context,
	workspaces workspace.WorkspaceList,
) tuikit.View {
	container := ctx.TUIContainer
	if len(workspaces.Items()) == 0 {
		container.HandleError(fmt.Errorf("no workspaces found"))
	}

	selectFunc := func(filterVal string) error {
		var ws *workspace.Workspace
		for _, s := range workspaces {
			if s.AssignedName() == filterVal || s.DisplayName == filterVal {
				ws = s
				break
			}
		}
		if ws == nil {
			return fmt.Errorf("workspace not found")
		}

		return ctx.SetView(NewWorkspaceView(ctx, ws))
	}

	return views.NewCollectionView(container.RenderState(), workspaces, types.CollectionFormatList, selectFunc)
}
