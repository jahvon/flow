package executable

import (
	"fmt"

	"github.com/jahvon/tuikit/components"
	"github.com/samber/lo"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/common"
)

func NewExecutableView(
	container *components.ContainerView,
	exec config.Executable,
	format components.Format,
) components.TeaModel {
	var executableKeyCallbacks = []components.KeyCallback{
		// TODO: implement run key that will run the executable
		//	- this will require come exec command logic to be moved to the runner package
		// {
		// 	Key: "r", Label: "run",
		// 	Callback: func() error {
		// 		return nil
		// 	},
		// },
		{
			Key: "e", Label: "edit",
			Callback: func() error {
				common.OpenInEditor(container, exec.DefinitionPath())
				return nil
			},
		},
	}
	state := &components.TerminalState{
		Theme:  io.Styles(),
		Height: container.Height(),
		Width:  container.Width(),
	}
	return components.NewEntityView(
		state,
		&exec,
		format,
		executableKeyCallbacks...,
	)
}

func NewExecutableListView(
	container *components.ContainerView,
	executables config.ExecutableList,
	format components.Format,
) components.TeaModel {
	if len(executables.Items()) == 0 {
		container.HandleError(fmt.Errorf("no workspaces found"))
	}

	selectFunc := func(filterVal string) error {
		exec, found := lo.Find(executables, func(e *config.Executable) bool {
			return e.ID() == filterVal
		})
		if !found {
			return fmt.Errorf("executable not found")
		}
		container.SetView(NewExecutableView(container, *exec, format))
		return nil
	}

	state := &components.TerminalState{
		Theme:  io.Styles(),
		Height: container.Height(),
		Width:  container.Width(),
	}
	return components.NewCollectionView(state, executables, format, selectFunc)
}
