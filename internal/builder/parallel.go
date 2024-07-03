package builder

import (
	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
)

var (
	parallelBaseDesc = "Parallel flows are executed concurrently."
)

func ParallelExecByRef(ctx *context.Context, name, definitionPath string) *config.Executable {
	e1 := ExecWithPauses(ctx, "parallel-exec-1", definitionPath)
	e2 := ExecWithPauses(ctx, "parallel-exec-2", definitionPath)
	e3 := ExecWithPauses(ctx, "parallel-exec-3", definitionPath)
	e := &config.Executable{
		Verb:       "start",
		Name:       name,
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: parallelBaseDesc +
			"\n\n- The environment defined on the root executable is inherited by all the child executables." +
			"\n- Setting `f:tmp` as the directory on child executables will create a shared temporary directory for execution.",
		Type: &config.ExecutableTypeSpec{
			Parallel: &config.ParallelExecutableType{
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

func ParallelExecWithExit(ctx *context.Context, name, definitionPath string) *config.Executable {
	e1 := ExecWithPauses(ctx, "parallel-exec-1", definitionPath)
	e2 := ExecWithExitCode(ctx, "parallel-exec-2", definitionPath, 1)
	e3 := ExecWithPauses(ctx, "parallel-exec-3", definitionPath)
	e := &config.Executable{
		Verb:       "start",
		Name:       name,
		Aliases:    []string{"parallel-exit"},
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: parallelBaseDesc +
			"\n\n The `failFast` option can be set to `true` to stop the flow if a child executable fails.",
		Type: &config.ExecutableTypeSpec{
			Parallel: &config.ParallelExecutableType{
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

func ParallelExecWithMaxThreads(ctx *context.Context, name, definitionPath string) *config.Executable {
	e := ParallelExecByRef(ctx, name, definitionPath)
	e.Description = parallelBaseDesc +
		"\n\nThe `maxThreads` option can be set to limit the number of concurrent executions."
	e.Type.Parallel.MaxThreads = 1
	return e
}

func ParallelExecByRefConfig(ctx *context.Context, name, definitionPath string) *config.Executable {
	e1 := ExecWithPauses(ctx, "parallel-exec-1", definitionPath)
	e2 := ExecWithArgs(ctx, "parallel-exec-2", definitionPath, config.ArgumentList{{Pos: 0, EnvKey: "ARG1"}})
	e := &config.Executable{
		Verb:       "run",
		Name:       name,
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: parallelBaseDesc +
			"\n\n- The environment defined on the root executable is inherited by all the child executables." +
			"\n- Setting `f:tmp` as the directory on child executables will create a shared temporary directory for execution.",
		Type: &config.ExecutableTypeSpec{
			Parallel: &config.ParallelExecutableType{
				Executables: []config.ParallelRefConfig{
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
	return e
}
