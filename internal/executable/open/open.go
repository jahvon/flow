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
	App     string `mapstructure:"app"     yaml:"app"`
	URI     string `mapstructure:"uri"     yaml:"uri"`
	Wait    bool   `mapstructure:"wait"    yaml:"wait"`
	Timeout string `mapstructure:"timeout" yaml:"timeout"`
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
		return fmt.Errorf("unable to decode 'open' executable spec - %w", err)
	}

	if strings.HasPrefix(openSpec.URI, "~/") {
		dir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("unable to get user home directory - %w", err)
		}
		openSpec.URI = filepath.Join(dir, openSpec.URI[2:])
	} else if strings.HasPrefix(openSpec.URI, "//") {
		_, workspacePath, _ := exec.GetContext()
		openSpec.URI = strings.Replace(openSpec.URI, "//", workspacePath+"/", 1)
	}

	err = executable.WithTimeout(openSpec.Timeout, func() error {
		if openSpec.App == "" {
			return open.Open(openSpec.URI, openSpec.Wait)
		} else {
			return open.OpenWith(openSpec.App, openSpec.URI, openSpec.Wait)
		}
	})

	return err
}
