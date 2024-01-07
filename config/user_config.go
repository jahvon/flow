package config

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type UserConfig struct {
	Workspaces       map[string]string `json:"workspaces"       yaml:"workspaces"`
	CurrentWorkspace string            `json:"currentWorkspace" yaml:"currentWorkspace"`
	CurrentNamespace string            `json:"currentNamespace" yaml:"currentNamespace"`
	InteractiveUI    bool              `json:"interactive"      yaml:"interactive"`
}

func (c *UserConfig) Validate() error {
	if c.CurrentWorkspace == "" {
		if _, found := c.Workspaces["default"]; found {
			c.CurrentWorkspace = "default"
		} else {
			return fmt.Errorf("current workspace is not set")
		}
	}
	if _, wsFound := c.Workspaces[c.CurrentWorkspace]; !wsFound {
		return fmt.Errorf("current workspace %s does not exist", c.CurrentWorkspace)
	}

	return nil
}

func (c *UserConfig) YAML() (string, error) {
	yamlBytes, err := yaml.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user config - %w", err)
	}
	return string(yamlBytes), nil
}

func (c *UserConfig) JSON(pretty bool) (string, error) {
	var jsonBytes []byte
	var err error
	if pretty {
		jsonBytes, err = json.MarshalIndent(c, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(c)
	}
	if err != nil {
		return "", fmt.Errorf("failed to marshal user config - %w", err)
	}
	return string(jsonBytes), nil
}

func (c *UserConfig) Markdown() string {
	mkdwn := "# Global Configurations\n"
	mkdwn += fmt.Sprintf("## Current workspace\n%s\n", c.CurrentWorkspace)
	if c.CurrentNamespace != "" {
		mkdwn += fmt.Sprintf("## Current namespace\n%s\n", c.CurrentNamespace)
	}
	interactive := "disabled"
	if c.InteractiveUI {
		interactive = "enabled"
	}
	mkdwn += fmt.Sprintf("## Interactive UI\n%s\n", interactive)
	mkdwn += "## Registered Workspaces\n"
	for name, path := range c.Workspaces {
		mkdwn += fmt.Sprintf("- %s: %s\n", name, path)
	}
	return mkdwn
}
