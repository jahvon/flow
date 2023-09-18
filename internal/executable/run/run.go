package run

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/jahvon/flow/internal/executable"
	"github.com/jahvon/flow/internal/executable/consts"
	"github.com/jahvon/flow/internal/services/run"
)

type agent struct {
}

type Spec struct {
	CommandStr string `json:"cmd"`
	Dir        string `json:"dir"`
}

func NewAgent() executable.Agent {
	return &agent{}
}
func (a *agent) Name() consts.AgentType {
	return consts.AgentTypeRun
}

func (a *agent) Exec(spec map[string]interface{}, _ *executable.Preference) error {
	if spec == nil {
		return fmt.Errorf("'run' executable spec cannot be empty")
	}

	var runSpec Spec
	err := mapstructure.Decode(spec, &runSpec)
	if err != nil {
		return fmt.Errorf("unable to decode 'run' executable spec - %v", err)
	}

	if runSpec.Dir != "" {
		return run.RunIn(runSpec.CommandStr, runSpec.Dir)
	}
	return run.Run(runSpec.CommandStr)
}
