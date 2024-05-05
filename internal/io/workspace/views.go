package workspace

import (
	"fmt"
	"path/filepath"

	"github.com/jahvon/tuikit/components"
	"github.com/jahvon/tuikit/styles"
	"github.com/samber/lo"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/common"
	"github.com/jahvon/flow/internal/services/open"
)

func NewWorkspaceView(
	container *components.ContainerView,
	ws config.WorkspaceConfig,
	format components.Format,
) components.TeaModel {
	var workspaceKeyCallbacks = []components.KeyCallback{
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
				common.DeprecatedOpenInEditor(container, filepath.Join(ws.Location(), file.WorkspaceConfigFileName))
				return nil
			},
		},
		{
			Key: "s", Label: "set",
			Callback: func() error {
				curCfg, err := file.LoadUserConfig()
				if err != nil {
					container.HandleError(err)
					return nil
				}
				curCfg.CurrentWorkspace = ws.AssignedName()
				if err := file.WriteUserConfig(curCfg); err != nil {
					container.HandleError(err)
				}
				container.SetContext(fmt.Sprintf("%s/*", ws.AssignedName()))
				container.SetNotice("workspace updatedd", styles.NoticeLevelInfo)
				return nil
			},
		},
	}

	state := &components.TerminalState{
		Theme:  io.Theme(),
		Height: container.Height(),
		Width:  container.Width(),
	}
	return components.NewEntityView(state, &ws, format, workspaceKeyCallbacks...)
}

func NewWorkspaceListView(
	container *components.ContainerView,
	workspaces config.WorkspaceConfigList,
	format components.Format,
) components.TeaModel {
	if len(workspaces.Items()) == 0 {
		container.HandleError(fmt.Errorf("no workspaces found"))
	}

	selectFunc := func(filterVal string) error {
		ws, found := lo.Find(workspaces, func(s config.WorkspaceConfig) bool {
			return s.AssignedName() == filterVal || s.DisplayName == filterVal
		})
		if !found {
			return fmt.Errorf("workspace not found")
		}

		container.SetView(NewWorkspaceView(container, ws, format))
		return nil
	}

	state := &components.TerminalState{
		Theme:  io.Theme(),
		Height: container.Height(),
		Width:  container.Width(),
	}
	return components.NewCollectionView(state, workspaces, format, selectFunc)
}
