package parallel

import (
	stdCtx "context"
	"fmt"
	"maps"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/flowexec/flow/internal/context"
	"github.com/flowexec/flow/internal/logger"
	"github.com/flowexec/flow/internal/runner"
	"github.com/flowexec/flow/internal/runner/engine"
	"github.com/flowexec/flow/internal/services/expr"
	"github.com/flowexec/flow/internal/services/store"
	argUtils "github.com/flowexec/flow/internal/utils/args"
	execUtils "github.com/flowexec/flow/internal/utils/executables"
	"github.com/flowexec/flow/types/executable"
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
	if err := runner.SetEnv(ctx.Config.CurrentVaultName(), e.Env(), inputEnv); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	if len(parallelSpec.Execs) > 0 {
		str, err := store.NewStore(store.Path())
		if err != nil {
			return err
		}
		if err := str.CreateBucket(store.EnvironmentBucket()); err != nil {
			return err
		}
		cacheData, err := str.GetAll()
		if err != nil {
			return err
		}
		if err := str.Close(); err != nil {
			logger.Log().Error(err, "unable to close store")
		}

		return handleExec(ctx, e, eng, parallelSpec, inputEnv, cacheData)
	}

	return fmt.Errorf("no parallel executables to run")
}

//nolint:funlen,gocognit
func handleExec(
	ctx *context.Context, parent *executable.Executable,
	eng engine.Engine,
	parallelSpec *executable.ParallelExecutableType,
	promptedEnv map[string]string,
	cacheData map[string]string,
) error {
	groupCtx, cancel := stdCtx.WithCancel(ctx.Ctx)
	defer cancel()
	group, _ := errgroup.WithContext(groupCtx)
	limit := parallelSpec.MaxThreads
	if limit == 0 {
		limit = len(parallelSpec.Execs)
	}
	group.SetLimit(limit)

	dataMap := expr.ExpressionEnv(ctx, parent, cacheData, promptedEnv)

	var execs []engine.Exec
	for i, refConfig := range parallelSpec.Execs {
		if refConfig.If != "" {
			if truthy, err := expr.IsTruthy(refConfig.If, &dataMap); err != nil {
				return err
			} else if !truthy {
				logger.Log().Debugf("skipping execution %d/%d", i+1, len(parallelSpec.Execs))
				continue
			}
		}
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
				logger.Log().Error(err, "unable to process arguments")
			}
			maps.Copy(execPromptedEnv, a)
		}

		switch {
		case exec.Exec != nil:
			fields := map[string]interface{}{"step": exec.Ref().String()}
			exec.Exec.SetLogFields(fields)
			if parallelSpec.Dir != "" && exec.Exec.Dir == "" {
				exec.Exec.Dir = parallelSpec.Dir
			}
		case exec.Parallel != nil:
			if parallelSpec.Dir != "" && exec.Parallel.Dir == "" {
				exec.Parallel.Dir = parallelSpec.Dir
			}
		case exec.Serial != nil:
			if parallelSpec.Dir != "" && exec.Serial.Dir == "" {
				exec.Serial.Dir = parallelSpec.Dir
			}
		}

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
