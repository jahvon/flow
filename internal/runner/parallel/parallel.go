package parallel

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
)

type parallelRunner struct{}

func NewRunner() runner.Runner {
	return &parallelRunner{}
}

func (r *parallelRunner) Name() string {
	return "parallel"
}

func (r *parallelRunner) IsCompatible(executable *config.Executable) bool {
	if executable == nil || executable.Type == nil || executable.Type.Parallel == nil {
		return false
	}
	return true
}

func (r *parallelRunner) Exec(
	ctx *context.Context,
	executable *config.Executable,
	promptedEnv map[string]string,
) error {
	parallelSpec := executable.Type.Parallel
	if err := runner.SetEnv(ctx.Logger, &parallelSpec.ExecutableEnvironment, promptedEnv); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	refs := parallelSpec.ExecutableRefs
	group, _ := errgroup.WithContext(ctx.Ctx)
	limit := parallelSpec.MaxThreads
	if limit == 0 {
		limit = len(refs)
	}
	group.SetLimit(limit)
	var errs []error
	for _, ref := range refs {
		ref = context.ExpandRef(ctx, ref)
		exec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, ref)
		if err != nil {
			return err
		}

		if exec.Type.Exec != nil {
			fields := map[string]interface{}{
				"executable": exec.ID(),
			}
			exec.Type.Exec.SetLogFields(fields)
		}

		group.Go(func() error {
			if parallelSpec.FailFast {
				return runner.Exec(ctx, exec, promptedEnv)
			} else {
				err := runner.Exec(ctx, exec, promptedEnv)
				if err != nil {
					errs = append(errs, err)
					ctx.Logger.Error(err, fmt.Sprintf("execution error for %s", ref))
				}
				return nil
			}
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
