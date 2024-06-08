package serial

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/utils"
)

type serialRunner struct{}

func NewRunner() runner.Runner {
	return &serialRunner{}
}

func (r *serialRunner) Name() string {
	return "serial"
}

func (r *serialRunner) IsCompatible(executable *config.Executable) bool {
	if executable == nil || executable.Type == nil || executable.Type.Serial == nil {
		return false
	}
	return true
}

func (r *serialRunner) Exec(ctx *context.Context, executable *config.Executable, promptedEnv map[string]string) error {
	serialSpec := executable.Type.Serial
	if err := runner.SetEnv(ctx.Logger, &serialSpec.ExecutableEnvironment, promptedEnv); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	if err := utils.ValidateOneOf(
		"executable list",
		serialSpec.ExecutableRefs, serialSpec.Executables,
	); err != nil {
		return err
	}

	if len(serialSpec.ExecutableRefs) > 0 {
		return handleExecRef(ctx, serialSpec, promptedEnv)
	} else if len(serialSpec.Executables) > 0 {
		return handleExec(ctx, executable, serialSpec, promptedEnv)
	}
	return fmt.Errorf("no executables to run")
}

func handleExecRef(ctx *context.Context, serialSpec *config.SerialExecutableType, promptedEnv map[string]string) error {
	order := serialSpec.ExecutableRefs
	var errs []error
	for i, executableRef := range order {
		ctx.Logger.Debugf("executing %s (%d/%d)", executableRef, i+1, len(order))
		exec, err := executableFromRef(ctx, executableRef)
		if err != nil {
			return err
		}

		if exec.Type.Exec != nil {
			fields := map[string]interface{}{
				"executable": exec.ID(),
			}
			exec.Type.Exec.SetLogFields(fields)
		}

		if err := runner.Exec(ctx, exec, promptedEnv); err != nil {
			if serialSpec.FailFast {
				return errors.Wrapf(err, "execution error for %s", executableRef)
			}
			errs = append(errs, err)
			ctx.Logger.Error(err, fmt.Sprintf("execution error for %s", executableRef))
		}
	}
	return nil
}

func handleExec(
	ctx *context.Context,
	parent *config.Executable,
	serialSpec *config.SerialExecutableType,
	promptedEnv map[string]string,
) error {
	var errs []error
	for i, refConfig := range serialSpec.Executables {
		var exec *config.Executable
		switch {
		case len(refConfig.Ref) > 0:
			var err error
			exec, err = executableFromRef(ctx, refConfig.Ref)
			if err != nil {
				return err
			}
			if exec.Type.Exec != nil {
				fields := map[string]interface{}{
					"executable": exec.ID(),
				}
				exec.Type.Exec.SetLogFields(fields)
			}
		case refConfig.Cmd != "":
			vis := config.VisibilityInternal
			exec = &config.Executable{
				Verb:       "exec",
				Name:       fmt.Sprintf("%s-cmd-%d", parent.Name, i),
				Visibility: &vis,
				Type: &config.ExecutableTypeSpec{
					Exec: &config.ExecExecutableType{
						Command: refConfig.Cmd,
					},
				},
			}
			fields := map[string]interface{}{"executable": exec.ID()}
			exec.Type.Exec.SetLogFields(fields)
			exec.SetContext(parent.Workspace(), parent.WorkspacePath(), parent.Namespace(), parent.DefinitionPath())
		}
		ctx.Logger.Debugf("executing %s (%d/%d)", exec.Ref(), i+1, len(serialSpec.Executables))

		for {
			if err := runner.Exec(ctx, exec, promptedEnv); err != nil {
				switch {
				case refConfig.Retries == 0 && serialSpec.FailFast:
					return errors.Wrapf(err, "execution error ref='%s'", exec.Ref())
				case refConfig.Retries == 0 && !serialSpec.FailFast:
					errs = append(errs, err)
					ctx.Logger.Error(err, fmt.Sprintf("execution error ref='%s'", exec.Ref()))
				case refConfig.Retries != 0 && refConfig.AttemptedMaxTimes() && serialSpec.FailFast:
					return fmt.Errorf("retries exceeded ref='%s' max=%d", exec.Ref(), refConfig.Retries)
				case refConfig.Retries != 0 && refConfig.AttemptedMaxTimes() && !serialSpec.FailFast:
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

		if i < len(serialSpec.Executables) && refConfig.ProceedPrompt {
			ctx.Logger.Println("Do you want to proceed with the next execution? (y/n)")
			if !inputConfirmed(os.Stdin) {
				return fmt.Errorf("stopping runner early (%d/%d)", i+1, len(serialSpec.Executables))
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%d execution errors - %v", len(errs), errs)
	}
	return nil
}

func executableFromRef(ctx *context.Context, ref config.Ref) (*config.Executable, error) {
	exec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, ref)
	if err != nil {
		return nil, err
	} else if exec == nil {
		return nil, fmt.Errorf("unable to find executable with reference %s", ref)
	}
	return exec, nil
}

func inputConfirmed(in *os.File) bool {
	var response string
	_, _ = fmt.Fscanf(in, response)
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
