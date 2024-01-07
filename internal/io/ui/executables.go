package ui

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/jahvon/flow/config"
)

func NewExecutableView(
	app *Application,
	exec config.Executable,
	format config.OutputFormat,
) ViewBuilder {
	var executableKeyCallbacks = []KeyCallback{
		{
			Key: "r", Label: "run",
			Callback: func() error {
				// TODO: implement
				return nil
			},
		},
		{
			Key: "e", Label: "edit",
			Callback: func() error {
				openInEditor(app, exec.DefinitionPath())
				return nil
			},
		},
	}
	return NewEntityView(
		app,
		&exec,
		format,
		executableKeyCallbacks...,
	)
}

func NewExecutableListView(
	app *Application,
	executables config.ExecutableList,
	format config.OutputFormat,
) ViewBuilder {
	if len(executables.Items()) == 0 {
		app.HandleInternalError(fmt.Errorf("no workspaces found"))
	}

	selectFunc := func(filterVal string) error {
		exec, found := lo.Find(executables, func(e *config.Executable) bool {
			return e.ID() == filterVal
		})
		if !found {
			return fmt.Errorf("executable not found")
		}
		app.BuildAndSetView(NewExecutableView(app, *exec, format))
		return nil
	}

	return NewCollectionView(app, executables, format, selectFunc)
}
