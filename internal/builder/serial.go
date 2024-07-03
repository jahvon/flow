package builder

import (
	"time"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
)

const (
	RefConfigCmd = "echo 'hello from a ref config!'"
)

var (
	serialBaseDesc = "Serial flows are executed in the order they are defined, one after the other."
)

func SerialExecByRef(ctx *context.Context, name, definitionPath string) (*config.Executable, config.ExecutableList) {
	e1 := SimpleExec(ctx, "serial-exec-1", definitionPath)
	e2 := SimpleExec(ctx, "serial-exec-2", definitionPath)
	e3 := SimpleExec(ctx, "serial-exec-3", definitionPath)
	e := &config.Executable{
		Verb:       "start",
		Name:       name,
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: serialBaseDesc +
			"\n\n- The environment defined on the root executable is inherited by all the child executables." +
			"\n- Setting `f:tmp` as the directory on child executables will create a shared temporary directory for execution.",
		Type: &config.ExecutableTypeSpec{
			Serial: &config.SerialExecutableType{
				ExecutableRefs: []config.Ref{e1.Ref(), e2.Ref(), e3.Ref()},
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e, config.ExecutableList{e1, e2, e3}
}

func SerialExecWithExit(ctx *context.Context, name, definitionPath string) *config.Executable {
	e1 := SimpleExec(ctx, "serial-exec-1", definitionPath)
	e2 := ExecWithExitCode(ctx, "serial-exec-2", definitionPath, 1)
	e3 := SimpleExec(ctx, "serial-exec-3", definitionPath)
	e := &config.Executable{
		Verb:       "start",
		Name:       name,
		Aliases:    []string{"serial-exit"},
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: serialBaseDesc +
			"\n\n The `failFast` option can be set to `true` to stop the flow if a child executable fails.",
		Type: &config.ExecutableTypeSpec{
			Serial: &config.SerialExecutableType{
				FailFast:       true,
				ExecutableRefs: []config.Ref{e1.Ref(), e2.Ref(), e3.Ref()},
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func SerialExecByRefConfig(
	ctx *context.Context, name, definitionPath string,
) (*config.Executable, config.ExecutableList) {
	e1 := SimpleExec(ctx, "serial-exec-1", definitionPath)
	e2 := ExecWithArgs(ctx, "serial-exec-2", definitionPath, config.ArgumentList{{Pos: 1, EnvKey: "ARG1"}})
	e := &config.Executable{
		Verb:       "run",
		Name:       name,
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: serialBaseDesc +
			"\n\n- The `executables` field can be used to define the child executables with more options. " +
			"This includes defining an executable inline, retries, arguments, and more.",
		Timeout: 10 * time.Second,
		Type: &config.ExecutableTypeSpec{
			Serial: &config.SerialExecutableType{
				Executables: []config.SerialRefConfig{
					{
						Ref: e1.Ref(),
					},
					{
						Ref:       e2.Ref(),
						Arguments: []string{"arg1"},
					},
					{
						Cmd: RefConfigCmd,
					},
				},
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e, config.ExecutableList{e1, e2}
}
