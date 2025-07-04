package exec

import (
	"github.com/pkg/errors"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/engine"
	"github.com/jahvon/flow/internal/services/run"
	"github.com/jahvon/flow/types/executable"
)

type execRunner struct{}

func NewRunner() runner.Runner {
	return &execRunner{}
}

func (r *execRunner) Name() string {
	return "exec"
}

func (r *execRunner) IsCompatible(executable *executable.Executable) bool {
	if executable == nil || executable.Exec == nil {
		return false
	}
	return true
}

func (r *execRunner) Exec(
	ctx *context.Context,
	e *executable.Executable,
	_ engine.Engine,
	inputEnv map[string]string,
) error {
	execSpec := e.Exec
	defaultEnv := runner.DefaultEnv(ctx, e)
	envMap, err := runner.BuildEnvMap(ctx.Logger, ctx.Config.CurrentVaultName(), e.Env(), inputEnv, defaultEnv)
	if err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}
	envList, err := runner.BuildEnvList(ctx.Logger, ctx.Config.CurrentVaultName(), e.Env(), inputEnv, defaultEnv)
	if err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	targetDir, isTmp, err := execSpec.Dir.ExpandDirectory(
		ctx.Logger,
		e.WorkspacePath(),
		e.FlowFilePath(),
		ctx.ProcessTmpDir,
		envMap,
	)
	if err != nil {
		return errors.Wrap(err, "unable to expand directory")
	} else if isTmp {
		ctx.ProcessTmpDir = targetDir
	}

	logMode := execSpec.LogMode
	logFields := execSpec.GetLogFields()

	switch {
	case execSpec.Cmd == "" && execSpec.File == "":
		return errors.New("either cmd or file must be specified")
	case execSpec.Cmd != "" && execSpec.File != "":
		return errors.New("cannot set both cmd and file")
	case execSpec.Cmd != "":
		return run.RunCmd(execSpec.Cmd, targetDir, envList, logMode, ctx.Logger, ctx.StdIn(), logFields)
	case execSpec.File != "":
		return run.RunFile(execSpec.File, targetDir, envList, logMode, ctx.Logger, ctx.StdIn(), logFields)
	default:
		return errors.New("unable to determine how e should be run")
	}
}
