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
	Timeout string                `yaml:"timeout" mapstructure:"timeout"`
	Dir     string                `yaml:"dir" mapstructure:"dir"`
	Params  []parameter.Parameter `yaml:"params" mapstructure:"params"`

	CommandStr string `yaml:"cmd" mapstructure:"cmd"`
	ShFile     string `yaml:"file" mapstructure:"file"`
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
		return fmt.Errorf("unable to decode 'run' executable spec - %v", err)
	}

	params := runSpec.Params
	envList, err := parameter.ParameterListToEnvList(params)
	if err != nil {
		return fmt.Errorf("unable to convert parameters to env list - %v", err)
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
		if runSpec.CommandStr != "" && runSpec.ShFile != "" {
			return fmt.Errorf("cannot set both cmd and file")
		} else if runSpec.CommandStr != "" {
			return run.RunCmd(runSpec.CommandStr, targetDir, envList)
		} else if runSpec.ShFile != "" {
			return run.RunFile(runSpec.ShFile, targetDir, envList)
		} else {
			return fmt.Errorf("either cmd or file must be specified")
		}
	})

	return err
}
