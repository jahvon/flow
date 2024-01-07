package ui

import (
	"fmt"
	"path/filepath"

	"github.com/samber/lo"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/services/open"
)

func NewWorkspaceView(app *Application, ws config.WorkspaceConfig, format config.OutputFormat) ViewBuilder {
	var workspaceKeyCallbacks = []KeyCallback{
		{
			Key: "o", Label: "open",
			Callback: func() error {
				if err := open.Open(ws.Location(), false); err != nil {
					log.Err(err).Msg("unable to open workspace")
					app.HandleInternalError(err)
				}
				return nil
			},
		},
		{
			Key: "e", Label: "edit",
			Callback: func() error {
				openInEditor(app, filepath.Join(ws.Location(), file.WorkspaceConfigFileName))
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
					app.HandleInternalError(err)
				}
				app.SetHeader(
					WithCurrentWorkspace(ws.AssignedName()),
					WithNotice("workspace updated", NoticeLevelInfo),
				)
				app.Update(syncAppMsg)
				return nil
			},
		},
	}

	return NewEntityView(app, &ws, format, workspaceKeyCallbacks...)
}

func NewWorkspaceListView(
	app *Application,
	workspaces config.WorkspaceConfigList,
	format config.OutputFormat,
) ViewBuilder {
	if len(workspaces.Items()) == 0 {
		app.HandleInternalError(fmt.Errorf("no workspaces found"))
	}

	selectFunc := func(filterVal string) error {
		ws, found := lo.Find(workspaces, func(s config.WorkspaceConfig) bool {
			return s.AssignedName() == filterVal || s.DisplayName == filterVal
		})
		if !found {
			return fmt.Errorf("workspace not found")
		}
		app.BuildAndSetView(NewWorkspaceView(app, ws, format))
		return nil
	}

	return NewCollectionView(app, workspaces, format, selectFunc)
}
