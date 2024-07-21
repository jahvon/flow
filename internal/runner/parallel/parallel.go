package parallel

import (
	stdCtx "context"
	"fmt"
	"maps"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/utils"
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

func (r *parallelRunner) Exec(ctx *context.Context, e *executable.Executable, promptedEnv map[string]string) error {
	parallelSpec := e.Parallel
	if err := runner.SetEnv(ctx.Logger, e.Env(), promptedEnv); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	if err := utils.ValidateOneOf("executable list", parallelSpec.Refs, parallelSpec.Execs); err != nil {
		return err
	}

	if len(parallelSpec.Refs) > 0 {
		return handleExecRef(ctx, parallelSpec, promptedEnv)
	} else if len(parallelSpec.Execs) > 0 {
		return handleExec(ctx, e, parallelSpec, promptedEnv)
	}

	return fmt.Errorf("no parallel executables to run")
}

func handleExecRef(
	ctx *context.Context, parallelSpec *executable.ParallelExecutableType, promptedEnv map[string]string,
) error {
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
		exec, err := execUtils.ExecutableForRef(ctx, ref)
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
					ctx.Logger.Errorx("execution error", "err", err, "ref", exec.Ref())
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

//nolint:gocognit
func handleExec(
	ctx *context.Context, parent *executable.Executable,
	parallelSpec *executable.ParallelExecutableType, promptedEnv map[string]string,
) error {
	groupCtx, cancel := stdCtx.WithCancel(ctx.Ctx)
	defer cancel()
	group, _ := errgroup.WithContext(groupCtx)
	limit := parallelSpec.MaxThreads
	if limit == 0 {
		limit = len(parallelSpec.Execs)
	}
	group.SetLimit(limit)
	var errs []error
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
		}

		execPromptedEnv := make(map[string]string)
		maps.Copy(promptedEnv, execPromptedEnv)
		if len(refConfig.Args) > 0 {
			a, err := argUtils.ProcessArgs(exec, refConfig.Args)
			if err != nil {
				ctx.Logger.Error(err, "unable to process arguments")
			}
			maps.Copy(execPromptedEnv, a)
		}

		if exec.Exec != nil {
			fields := map[string]interface{}{
				"e": exec.ID(),
			}
			exec.Exec.SetLogFields(fields)
		}

		group.Go(func() error {
			var attempts int
		retryLoop:
			for {
				attempts++
				if err := runner.Exec(ctx, exec, execPromptedEnv); err != nil {
					switch {
					case refConfig.Retries == 0 && parallelSpec.FailFast:
						return errors.Wrapf(err, "execution error ref='%s'", exec.Ref())
					case refConfig.Retries == 0 && !parallelSpec.FailFast:
						errs = append(errs, err)
						ctx.Logger.Errorx("execution error", "err", err, "ref", exec.Ref())
						break retryLoop
					case refConfig.Retries != 0 && attempts-1 >= refConfig.Retries && parallelSpec.FailFast:
						return fmt.Errorf("retries exceeded ref='%s' max=%d", exec.Ref(), refConfig.Retries)
					case refConfig.Retries != 0 && attempts-1 >= refConfig.Retries && !parallelSpec.FailFast:
						errs = append(errs, err)
						ctx.Logger.Errorx(
							"retries exceeded", "err", err, "ref", exec.Ref(), "max", refConfig.Retries,
						)
						break retryLoop
					case refConfig.Retries != 0 && attempts-1 < refConfig.Retries:
						ctx.Logger.Warnx("retrying", "ref", exec.Ref())
					default:
						return errors.Wrapf(err, "unexpected error handling ref='%s'", exec.Ref())
					}
				} else {
					break retryLoop
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
