package run

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/jahvon/flow/internal/executable"
	"github.com/jahvon/flow/internal/executable/consts"
	"github.com/jahvon/flow/internal/executable/parameter"
	"github.com/jahvon/flow/internal/services/run"
)

type agent struct {
}

type Spec struct {
	Timeout string                `mapstructure:"timeout" yaml:"timeout"`
	Dir     string                `mapstructure:"dir"     yaml:"dir"`
	Params  []parameter.Parameter `mapstructure:"params"  yaml:"params"`

	CommandStr string `mapstructure:"cmd"  yaml:"cmd"`
	ShFile     string `mapstructure:"file" yaml:"file"`
}

func NewAgent() executable.Agent {
	return &agent{}
}
func (a *agent) Name() consts.AgentType {
	return consts.AgentTypeRun
}

func (a *agent) Exec(exec executable.Executable) error {
	if exec.Spec == nil {
		return fmt.Errorf("'run' executable spec cannot be empty")
	}

	var runSpec Spec
	err := mapstructure.Decode(exec.Spec, &runSpec)
	if err != nil {
		return fmt.Errorf("unable to decode 'run' executable spec - %w", err)
	}

	params := runSpec.Params
	envList, err := parameter.ParameterListToEnvList(params)
	if err != nil {
		return fmt.Errorf("unable to convert parameters to env list - %w", err)
	}

	targetDir := runSpec.Dir
	if targetDir == "" {
		_, workspacePath, _ := exec.GetContext()
		targetDir = workspacePath
	} else if strings.HasPrefix(targetDir, "//") {
		_, workspacePath, _ := exec.GetContext()
		targetDir = strings.Replace(targetDir, "//", workspacePath+"/", 1)
	}

	err = executable.WithTimeout(runSpec.Timeout, func() error {
		switch {
		case runSpec.CommandStr == "" && runSpec.ShFile == "":
			return fmt.Errorf("either cmd or file must be specified")
		case runSpec.CommandStr != "" && runSpec.ShFile != "":
			return fmt.Errorf("cannot set both cmd and file")
		case runSpec.CommandStr != "":
			return run.RunCmd(runSpec.CommandStr, targetDir, envList)
		case runSpec.ShFile != "":
			return run.RunFile(runSpec.ShFile, targetDir, envList)
		default:
			return fmt.Errorf("unable to determine how executable should be run")
		}
	})

	return err
}
