package open

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/jahvon/flow/internal/executable"
	"github.com/jahvon/flow/internal/executable/consts"
	"github.com/jahvon/flow/internal/services/open"
)

type agent struct {
}

type Spec struct {
	App  string `json:"app"`
	Uri  string `json:"uri"`
	Wait bool   `json:"wait"`
}

func NewAgent() executable.Agent {
	return &agent{}
}
func (a *agent) Name() consts.AgentType {
	return consts.AgentTypeOpen
}

func (a *agent) Exec(spec map[string]interface{}, _ *executable.Preference) error {
	if spec == nil {
		return fmt.Errorf("'open' executable spec cannot be empty")
	}

	var openSpec Spec
	err := mapstructure.Decode(spec, &openSpec)
	if err != nil {
		return fmt.Errorf("unable to decode 'open' executable spec - %v", err)
	}

	if openSpec.App == "" {
		return open.Open(openSpec.Uri, openSpec.Wait)
	} else {
		return open.OpenWith(openSpec.App, openSpec.Uri, openSpec.Wait)
	}
}
