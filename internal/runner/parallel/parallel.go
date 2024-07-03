package parallel

import (
	"fmt"
	"maps"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/jahvon/flow/config"
	argUtils "github.com/jahvon/flow/config/args"
	"github.com/jahvon/flow/config/executables"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/utils"
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

	if err := utils.ValidateOneOf(
		"executable list",
		parallelSpec.ExecutableRefs, parallelSpec.Executables,
	); err != nil {
		return err
	}

	if len(parallelSpec.ExecutableRefs) > 0 {
		return handleExecRef(ctx, parallelSpec, promptedEnv)
	} else if len(parallelSpec.Executables) > 0 {
		return handleExec(ctx, executable, parallelSpec, promptedEnv)
	}

	return fmt.Errorf("no parallel executables to run")
}

func handleExecRef(
	ctx *context.Context, parallelSpec *config.ParallelExecutableType, promptedEnv map[string]string,
) error {
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
		exec, err := executables.ExecutableForRef(ctx, ref)
		if err != nil {
			return err
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

//nolint:gocognit
func handleExec(
	ctx *context.Context, parent *config.Executable,
	parallelSpec *config.ParallelExecutableType, promptedEnv map[string]string,
) error {
	group, _ := errgroup.WithContext(ctx.Ctx)
	limit := parallelSpec.MaxThreads
	if limit == 0 {
		limit = len(parallelSpec.Executables)
	}
	group.SetLimit(limit)
	var errs []error
	for i, refConfig := range parallelSpec.Executables {
		var exec *config.Executable
		switch {
		case len(refConfig.Ref) > 0:
			var err error
			exec, err = executables.ExecutableForRef(ctx, refConfig.Ref)
			if err != nil {
				return err
			}
		case refConfig.Cmd != "":
			exec = executables.ExecutableForCmd(parent, refConfig.Cmd, i)
		}

		if len(refConfig.Arguments) > 0 {
			a, err := argUtils.ProcessArgs(exec, refConfig.Arguments)
			if err != nil {
				ctx.Logger.Error(err, "unable to process arguments")
			}
			maps.Copy(promptedEnv, a)
		}

		group.Go(func() error {
			for {
				if err := runner.Exec(ctx, exec, promptedEnv); err != nil {
					switch {
					case refConfig.Retries == 0 && parallelSpec.FailFast:
						return errors.Wrapf(err, "execution error ref='%s'", exec.Ref())
					case refConfig.Retries == 0 && !parallelSpec.FailFast:
						errs = append(errs, err)
						ctx.Logger.Error(err, fmt.Sprintf("execution error ref='%s'", exec.Ref()))
					case refConfig.Retries != 0 && refConfig.AttemptedMaxTimes() && parallelSpec.FailFast:
						return fmt.Errorf("retries exceeded ref='%s' max=%d", exec.Ref(), refConfig.Retries)
					case refConfig.Retries != 0 && refConfig.AttemptedMaxTimes() && !parallelSpec.FailFast:
						errs = append(errs, err)
						ctx.Logger.Error(err, fmt.Sprintf("retries exceeded ref='%s' max=%d", exec.Ref(), refConfig.Retries))
					case refConfig.Retries != 0 && !refConfig.AttemptedMaxTimes():
						refConfig.RecordAttempt()
						ctx.Logger.Warnf("retrying ref='%s'", exec.Ref())
					default:
						return errors.Wrapf(err, "unexpected error handling ref='%s'", exec.Ref())
					}
					continue
				}
				break
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
