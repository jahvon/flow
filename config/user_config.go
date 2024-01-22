package config

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type InteractiveConfig struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
	// +docsgen:notify
	// Whether to send a desktop notification when a command completes.
	NotifyOnCompletion bool `json:"notify" yaml:"notify"`
	// +docsgen:sound
	// Whether to play a sound when a command completes.
	SoundOnCompletion bool `json:"sound" yaml:"sound"`
}

type UserConfig struct {
	// +docsgen:workspaces
	// Map of workspace names to their paths.
	Workspaces map[string]string `json:"workspaces" yaml:"workspaces"`
	// +docsgen:currentWorkspace
	// The name of the current workspace. This should match a key in the `workspaces` map.
	CurrentWorkspace string `json:"currentWorkspace" yaml:"currentWorkspace"`
	// +docsgen:currentNamespace
	// The name of the current namespace. This is not required to be set.
	CurrentNamespace string `json:"currentNamespace" yaml:"currentNamespace"`
	// +docsgen:interactive
	// Configurations for the interactive UI.
	Interactive *InteractiveConfig `json:"interactive" yaml:"interactive"`
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
	if c.Interactive != nil {
		interactiveConfig, err := yaml.Marshal(c.Interactive)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal interactive config")
			mkdwn += "## Interactive UI config\nerror\n"
		} else {
			mkdwn += fmt.Sprintf("## Git UI config\n```yaml\n%s```\n", string(interactiveConfig))
		}
	}
	mkdwn += "## Registered Workspaces\n"
	for name, path := range c.Workspaces {
		mkdwn += fmt.Sprintf("- %s: %s\n", name, path)
	}
	return mkdwn
}
