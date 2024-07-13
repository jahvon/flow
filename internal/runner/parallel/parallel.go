package parallel

import (
	stdCtx "context"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
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

func (r *parallelRunner) Exec(ctx *context.Context, e *executable.Executable, promptedEnv map[string]string) error {
	parallelSpec := e.Parallel
	if err := runner.SetEnv(ctx.Logger, e.Env(), promptedEnv); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	refs := parallelSpec.Refs
	groupCtx, cancel := stdCtx.WithCancel(ctx.Ctx)
	defer cancel()
	group, _ := errgroup.WithContext(groupCtx)
	limit := parallelSpec.MaxThreads
	if limit == 0 {
		limit = 5
	}
	group.SetLimit(limit)
	var errs []error
	for _, ref := range refs {
		ref = context.ExpandRef(ctx, ref)
		exec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, ref)
		if err != nil {
			return err
		}

		if exec.Exec != nil {
			fields := map[string]interface{}{
				"e": exec.ID(),
			}
			exec.Exec.SetLogFields(fields)
		}

		group.Go(func() error {
			if parallelSpec.FailFast {
				if err := runner.Exec(ctx, exec, promptedEnv); err != nil {
					cancel()
					return err
				}
			} else {
				err := runner.Exec(ctx, exec, promptedEnv)
				if err != nil {
					errs = append(errs, err)
					ctx.Logger.Error(err, fmt.Sprintf("execution error for %s", ref))
				}
			}
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return errors.Wrap(err, "parallel execution error")
	}

	if len(errs) > 0 {
		return fmt.Errorf("%d execution errors - %v", len(errs), errs)
	}
	return nil
}
