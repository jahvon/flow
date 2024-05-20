package executable

import (
	"fmt"

	"github.com/jahvon/tuikit/components"
	"github.com/jahvon/tuikit/styles"
	"github.com/samber/lo"
	"golang.design/x/clipboard"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/common"
)

func NewExecutableView(
	ctx *context.Context,
	exec config.Executable,
	format components.Format,
) components.TeaModel {
	container := ctx.InteractiveContainer
	var executableKeyCallbacks = []components.KeyCallback{
		{
			Key: "c", Label: "copy ref",
			Callback: func() error {
				err := clipboard.Init()
				if err != nil {
					return err
				}
				clipboard.Write(clipboard.FmtText, []byte(exec.Ref().String()))
				container.SetNotice("copied reference to clipboard", styles.NoticeLevelInfo)
				return nil
			},
		},
		{
			Key: "e", Label: "edit",
			Callback: func() error {
				if err := common.OpenInEditor(exec.DefinitionPath(), ctx.StdIn(), ctx.StdOut()); err != nil {
					container.HandleError(fmt.Errorf("unable to open executable: %w", err))
				}
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
	ctx *context.Context,
	executables config.ExecutableList,
	format components.Format,
) components.TeaModel {
	container := ctx.InteractiveContainer
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
		container.SetView(NewExecutableView(ctx, *exec, format))
		return nil
	}

	state := &components.TerminalState{
		Theme:  io.Theme(),
		Height: container.Height(),
		Width:  container.Width(),
	}
	return components.NewCollectionView(state, executables, format, selectFunc)
}
