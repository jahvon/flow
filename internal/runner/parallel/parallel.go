package parallel

import (
	"fmt"

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

func (r *parallelRunner) Exec(ctx *context.Context, executable *config.Executable) error {
	parallelSpec := executable.Type.Parallel
	if err := runner.SetEnv(&parallelSpec.ParameterizedExecutable); err != nil {
		return fmt.Errorf("unable to set parameters to env - %w", err)
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
		exec, err := ctx.ExecutableCache.GetExecutableByRef(ref)
		if err != nil {
			return fmt.Errorf("unable to get executable by ref %s - %w", ref, err)
		} else if executable == nil {
			return fmt.Errorf("unable to find executable with reference %s", ref)
		}
		group.Go(func() error {
			if parallelSpec.FailFast {
				return runner.Exec(ctx, exec)
			} else {
				err := runner.Exec(ctx, exec)
				if err != nil {
					errs = append(errs, err)
				}
				return nil
			}
		})
	}
	if err := group.Wait(); err != nil {
		return fmt.Errorf("unable to execute parallel executables - %w", err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("%d execution errors - %v", len(errs), errs)
	}
	return nil
}
