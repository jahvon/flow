package executable

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jahvon/tuikit/components"
	"github.com/jahvon/tuikit/styles"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/common"
	"github.com/jahvon/flow/types/executable"
)

func NewExecutableView(
	ctx *context.Context,
	exec executable.Executable,
	format components.Format,
	runFunc func(string) error,
) components.TeaModel {
	container := ctx.InteractiveContainer
	var executableKeyCallbacks = []components.KeyCallback{
		{
			Key: "r", Label: "run",
			Callback: func() error {
				ctx.InteractiveContainer.Shutdown()
				return runFunc(exec.Ref().String())
			},
		},
		{
			Key: "c", Label: "copy ref",
			Callback: func() error {
				if err := clipboard.WriteAll(exec.Ref().String()); err != nil {
					container.HandleError(fmt.Errorf("unable to copy reference to clipboard: %w", err))
				} else {
					container.SetNotice("copied reference to clipboard", styles.NoticeLevelInfo)
				}
				return nil
			},
		},
		{
			Key: "e", Label: "edit",
			Callback: func() error {
				if err := common.OpenInEditor(exec.FlowFilePath(), ctx.StdIn(), ctx.StdOut()); err != nil {
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
	executables executable.ExecutableList,
	format components.Format,
	runFunc func(string) error,
) components.TeaModel {
	container := ctx.InteractiveContainer
	if len(executables.Items()) == 0 {
		container.HandleError(fmt.Errorf("no workspaces found"))
	}

	selectFunc := func(filterVal string) error {
		s := strings.Split(filterVal, " ")
		if len(s) != 2 {
			return fmt.Errorf("invalid filter value")
		}
		verb, id := s[0], s[1]
		exec, err := executables.FindByVerbAndID(executable.Verb(verb), id)
		if err != nil {
			return fmt.Errorf("executable not found")
		}
		container.SetView(NewExecutableView(ctx, *exec, format, runFunc))
		return nil
	}

	state := &components.TerminalState{
		Theme:  io.Theme(),
		Height: container.Height(),
		Width:  container.Width(),
	}
	return components.NewCollectionView(state, executables, format, selectFunc)
}
