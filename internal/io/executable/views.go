package executable

import (
	"fmt"

	"github.com/jahvon/tuikit/components"
	"github.com/jahvon/tuikit/styles"
	"github.com/samber/lo"
	"golang.design/x/clipboard"

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
		{
			Key: "c", Label: "copy ref",
			Callback: func() error {
				err := clipboard.Init()
				if err != nil {
					panic(err)
				}
				clipboard.Write(clipboard.FmtText, []byte(exec.Ref().String()))
				container.SetNotice("copied reference to clipboard", styles.NoticeLevelInfo)
				return nil
			},
		},
		{
			Key: "e", Label: "edit",
			Callback: func() error {
				common.DeprecatedOpenInEditor(container, exec.DefinitionPath())
				return nil
			},
		},
	}
	state := &components.TerminalState{
		Theme:  io.Theme(),
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
		Theme:  io.Theme(),
		Height: container.Height(),
		Width:  container.Width(),
	}
	return components.NewCollectionView(state, executables, format, selectFunc)
}
