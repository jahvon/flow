package run

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/jahvon/flow/internal/executable"
	"github.com/jahvon/flow/internal/executable/consts"
	"github.com/jahvon/flow/internal/executable/parameter"
	"github.com/jahvon/flow/internal/services/run"
)

type agent struct {
}

type Spec struct {
	Timeout    string                `yaml:"timeout" mapstructure:"timeout"`
	CommandStr string                `yaml:"cmd" mapstructure:"cmd"`
	Dir        string                `yaml:"dir" mapstructure:"dir"`
	Params     []parameter.Parameter `yaml:"params" mapstructure:"params"`
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

	err = executable.WithTimeout(runSpec.Timeout, func() error {
		if runSpec.Dir != "" {
			return run.RunCmdIn(runSpec.CommandStr, runSpec.Dir)
		}
		return run.RunCmd(runSpec.CommandStr)
	})

	return err
}
