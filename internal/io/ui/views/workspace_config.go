package views

import (
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/io/ui/types"
	"github.com/jahvon/flow/internal/services/open"
)

func NewWorkspaceView(
	parent types.ParentView,
	ws config.WorkspaceConfig,
	format config.OutputFormat,
) types.ViewBuilder {
	var workspaceKeyCallbacks = []types.KeyCallback{
		{
			Key: "o", Label: "open",
			Callback: func() error {
				if err := open.Open(ws.Location(), false); err != nil {
					parent.HandleInternalError(fmt.Errorf("unable to open workspace: %w", err))
				}
				return nil
			},
		},
		{
			Key: "e", Label: "edit",
			Callback: func() error {
				openInEditor(parent, filepath.Join(ws.Location(), file.WorkspaceConfigFileName))
				return nil
			},
		},
		{
			Key: "s", Label: "set",
			Callback: func() error {
				curCfg := file.LoadUserConfig()
				curCfg.CurrentWorkspace = ws.AssignedName()
				if err := file.WriteUserConfig(curCfg); err != nil {
					log.Err(err).Msg("failed to write user config")
					parent.HandleInternalError(err)
				}
				parent.SetContext(ws.AssignedName(), "")
				parent.SetNotice("workspace updated", types.NoticeLevelInfo)
				return nil
			},
		},
	}

	return NewEntityView(parent, &ws, format, workspaceKeyCallbacks...)
}

func NewWorkspaceListView(
	parent types.ParentView,
	workspaces config.WorkspaceConfigList,
	format config.OutputFormat,
) types.ViewBuilder {
	if len(workspaces.Items()) == 0 {
		parent.HandleInternalError(fmt.Errorf("no workspaces found"))
	}

	selectFunc := func(filterVal string) error {
		ws, found := lo.Find(workspaces, func(s config.WorkspaceConfig) bool {
			return s.AssignedName() == filterVal || s.DisplayName == filterVal
		})
		if !found {
			return fmt.Errorf("workspace not found")
		}
		parent.BuildAndSetView(NewWorkspaceView(parent, ws, format))
		return nil
	}

	return NewCollectionView(parent, workspaces, format, selectFunc)
}
