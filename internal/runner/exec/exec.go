package exec

import (
	"fmt"

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
	if executable == nil || executable.Type == nil || executable.Type.Parallel == nil {
		return false
	}
	return true
}

func (r *execRunner) Exec(_ *context.Context, executable *config.Executable) error {
	execSpec := executable.Type.Exec
	envMap, err := runner.ParametersToEnvMap(&execSpec.ParameterizedExecutable)
	if err != nil {
		return fmt.Errorf("unable to convert parameters to env map - %w", err)
	}
	envList, err := runner.ParametersToEnvList(&execSpec.ParameterizedExecutable)
	if err != nil {
		return fmt.Errorf("unable to convert parameters to env list - %w", err)
	}

	targetDir := execSpec.ExpandDirectory(executable.WorkspacePath(), executable.DefinitionPath(), envMap)
	logMode := execSpec.LogMode

	switch {
	case execSpec.Command == "" && execSpec.File == "":
		return fmt.Errorf("either cmd or file must be specified")
	case execSpec.Command != "" && execSpec.File != "":
		return fmt.Errorf("cannot set both cmd and file")
	case execSpec.Command != "":
		return run.RunCmd(execSpec.Command, targetDir, envList, logMode)
	case execSpec.File != "":
		return run.RunFile(execSpec.File, targetDir, envList, logMode)
	default:
		return fmt.Errorf("unable to determine how executable should be run")
	}
}
