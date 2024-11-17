package parallel

import (
	"fmt"
	"maps"

	"github.com/pkg/errors"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/engine"
	argUtils "github.com/jahvon/flow/internal/utils/args"
	execUtils "github.com/jahvon/flow/internal/utils/executables"
	"github.com/jahvon/flow/types/executable"
)

type parallelRunner struct{}

func NewRunner() runner.Runner {
	return &parallelRunner{}
}

func (r *parallelRunner) Name() string {
	return "parallel"
}

func (r *parallelRunner) IsCompatible(executable *executable.Executable) bool {
	if executable == nil || executable.Parallel == nil {
		return false
	}
	return true
}

func (r *parallelRunner) Exec(
	ctx *context.Context,
	e *executable.Executable,
	eng engine.Engine,
	inputEnv map[string]string,
) error {
	parallelSpec := e.Parallel
	if err := runner.SetEnv(ctx.Logger, e.Env(), inputEnv); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	if len(parallelSpec.Execs) > 0 {
		return handleExec(ctx, e, eng, parallelSpec, inputEnv)
	}

	return fmt.Errorf("no parallel executables to run")
}

func handleExec(
	ctx *context.Context, parent *executable.Executable,
	eng engine.Engine,
	parallelSpec *executable.ParallelExecutableType, promptedEnv map[string]string,
) error {
	var execs []engine.Exec
	for i, refConfig := range parallelSpec.Execs {
		var exec *executable.Executable
		switch {
		case len(refConfig.Ref) > 0:
			var err error
			exec, err = execUtils.ExecutableForRef(ctx, refConfig.Ref)
			if err != nil {
				return err
			}
		case refConfig.Cmd != "":
			exec = execUtils.ExecutableForCmd(parent, refConfig.Cmd, i)
		default:
			return errors.New("parallel executable must have a ref or cmd")
		}

		execPromptedEnv := make(map[string]string)
		maps.Copy(promptedEnv, execPromptedEnv)
		if len(refConfig.Args) > 0 {
			a, err := argUtils.ProcessArgs(exec, refConfig.Args, execPromptedEnv)
			if err != nil {
				ctx.Logger.Error(err, "unable to process arguments")
			}
			maps.Copy(execPromptedEnv, a)
		}

		fields := map[string]interface{}{
			"step": exec.ID(),
		}
		exec.Exec.SetLogFields(fields)

		runExec := func() error {
			err := runner.Exec(ctx, exec, eng, execPromptedEnv)
			if err != nil {
				return err
			}
			return nil
		}

		execs = append(execs, engine.Exec{ID: exec.Ref().String(), Function: runExec, MaxRetries: refConfig.Retries})
	}
	results := eng.Execute(
		ctx.Ctx, execs,
		engine.WithMode(engine.Parallel),
		engine.WithFailFast(parent.Parallel.FailFast),
		engine.WithMaxThreads(parent.Parallel.MaxThreads),
	)
	if results.HasErrors() {
		return errors.New(results.String())
	}
	return nil
}
