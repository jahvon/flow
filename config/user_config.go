package config

import (
	"fmt"
)

type UserConfig struct {
	Workspaces       map[string]string `json:"workspaces"       yaml:"workspaces"`
	CurrentWorkspace string            `json:"currentWorkspace" yaml:"currentWorkspace"`
	CurrentNamespace string            `json:"currentNamespace" yaml:"currentNamespace"`
}

func (c *UserConfig) Validate() error {
	if c.CurrentWorkspace == "" {
		return fmt.Errorf("current workspace is not set")
	}
	if _, wsFound := c.Workspaces[c.CurrentWorkspace]; !wsFound {
		return fmt.Errorf("current workspace %s does not exist", c.CurrentWorkspace)
	}

	return nil
}
