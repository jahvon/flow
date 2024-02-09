package exec

import (
	"github.com/pkg/errors"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/services/run"
)

type execRunner struct{}

func NewRunner() runner.Runner {
	return &execRunner{}
}

func (r *execRunner) Name() string {
	return "exec"
}

func (r *execRunner) IsCompatible(executable *config.Executable) bool {
	if executable == nil || executable.Type == nil || executable.Type.Exec == nil {
		return false
	}
	return true
}

func (r *execRunner) Exec(ctx *context.Context, executable *config.Executable, promptedEnv map[string]string) error {
	execSpec := executable.Type.Exec
	promptedEnv = applyBaseEnv(ctx, executable, promptedEnv)
	envMap, err := runner.ParametersToEnvMap(ctx.Logger, &execSpec.ParameterizedExecutable, promptedEnv)
	if err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}
	envList, err := runner.ParametersToEnvList(ctx.Logger, &execSpec.ParameterizedExecutable, promptedEnv)
	if err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	targetDir, isTmp, err := execSpec.ExpandDirectory(
		ctx.Logger,
		executable.WorkspacePath(),
		executable.DefinitionPath(),
		ctx.ProcessTmpDir,
		envMap,
	)
	if err != nil {
		return errors.Wrap(err, "unable to expand directory")
	} else if isTmp {
		ctx.ProcessTmpDir = targetDir
	}

	logMode := execSpec.LogMode
	var logFields map[string]interface{}

	if logMode == config.StructuredLogMode {
		logFields = execSpec.GetLogFields()
	}

	switch {
	case execSpec.Command == "" && execSpec.File == "":
		return errors.New("either cmd or file must be specified")
	case execSpec.Command != "" && execSpec.File != "":
		return errors.New("cannot set both cmd and file")
	case execSpec.Command != "":
		return run.RunCmd(execSpec.Command, targetDir, envList, logMode, ctx.Logger, logFields)
	case execSpec.File != "":
		return run.RunFile(execSpec.File, targetDir, envList, logMode, ctx.Logger, logFields)
	default:
		return errors.New("unable to determine how executable should be run")
	}
}

func applyBaseEnv(ctx *context.Context, executable *config.Executable, envMap map[string]string) map[string]string {
	if envMap == nil {
		envMap = make(map[string]string)
	}
	envMap["FLOW_RUNNER"] = "true"
	envMap["FLOW_CURRENT_WORKSPACE"] = ctx.UserConfig.CurrentWorkspace
	envMap["FLOW_CURRENT_NAMESPACE"] = ctx.UserConfig.CurrentNamespace
	envMap["FLOW_EXECUTABLE_NAME"] = executable.Name
	envMap["FLOW_DEFINITION_PATH"] = executable.DefinitionPath()
	envMap["FLOW_WORKSPACE_PATH"] = executable.WorkspacePath()
	envMap["DISABLE_FLOW_INTERACTIVE"] = "true"
	return envMap
}
