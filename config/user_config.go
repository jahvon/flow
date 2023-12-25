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
	UIEnabled        bool              `json:"uiEnabled"        yaml:"uiEnabled"`
	AppPreferences   struct {
		Open string `json:"open"  yaml:"open"`
		Edit string `json:"edit"  yaml:"edit"`
	} `json:"appPreferences" yaml:"appPreferences"`
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
		return "", fmt.Errorf("failed to marshal user config - %v", err)
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
		return "", fmt.Errorf("failed to marshal user config - %v", err)
	}
	return string(jsonBytes), nil
}

func (c *UserConfig) Map() map[string]string {
	fields := make(map[string]string)
	fields["Current workspace"] = c.CurrentWorkspace
	fields["Current namespace"] = c.CurrentNamespace
	if c.CurrentNamespace == "" {
		fields["Current namespace"] = "*"
	}
	fields["UI enabled"] = fmt.Sprintf("%t", c.UIEnabled)
	return fields
}
