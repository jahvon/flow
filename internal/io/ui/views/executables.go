package views

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io/ui/types"
)

func NewExecutableView(
	parent types.ParentView,
	exec config.Executable,
	format config.OutputFormat,
) types.ViewBuilder {
	var executableKeyCallbacks = []types.KeyCallback{
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
				openInEditor(parent, exec.DefinitionPath())
				return nil
			},
		},
	}
	return NewEntityView(
		parent,
		&exec,
		format,
		executableKeyCallbacks...,
	)
}

func NewExecutableListView(
	parent types.ParentView,
	executables config.ExecutableList,
	format config.OutputFormat,
) types.ViewBuilder {
	if len(executables.Items()) == 0 {
		parent.HandleInternalError(fmt.Errorf("no workspaces found"))
	}

	selectFunc := func(filterVal string) error {
		exec, found := lo.Find(executables, func(e *config.Executable) bool {
			return e.ID() == filterVal
		})
		if !found {
			return fmt.Errorf("executable not found")
		}
		parent.BuildAndSetView(NewExecutableView(parent, *exec, format))
		return nil
	}

	return NewCollectionView(parent, executables, format, selectFunc)
}
