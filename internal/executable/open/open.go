package open

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/jahvon/flow/internal/executable"
	"github.com/jahvon/flow/internal/executable/consts"
	"github.com/jahvon/flow/internal/services/open"
)

type agent struct {
}

type Spec struct {
	App     string `yaml:"app" mapstructure:"app"`
	Uri     string `yaml:"uri" mapstructure:"uri"`
	Wait    bool   `yaml:"wait" mapstructure:"wait"`
	Timeout string `yaml:"timeout" mapstructure:"timeout"`
}

func NewAgent() executable.Agent {
	return &agent{}
}
func (a *agent) Name() consts.AgentType {
	return consts.AgentTypeOpen
}

func (a *agent) Exec(exec executable.Executable) error {
	if exec.Spec == nil {
		return fmt.Errorf("'open' executable spec cannot be empty")
	}

	var openSpec Spec
	err := mapstructure.Decode(exec.Spec, &openSpec)
	if err != nil {
		return fmt.Errorf("unable to decode 'open' executable spec - %v", err)
	}

	if strings.HasPrefix(openSpec.Uri, "~/") {
		dir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("unable to get user home directory - %v", err)
		}
		openSpec.Uri = filepath.Join(dir, openSpec.Uri[2:])
	}

	err = executable.WithTimeout(openSpec.Timeout, func() error {
		if openSpec.App == "" {
			return open.Open(openSpec.Uri, openSpec.Wait)
		} else {
			return open.OpenWith(openSpec.App, openSpec.Uri, openSpec.Wait)
		}
	})

	return err
}
