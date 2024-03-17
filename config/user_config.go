package config

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type WorkspaceMode string

const (
	WorkspaceModeFixed   WorkspaceMode = "fixed"
	WorkspaceModeDynamic WorkspaceMode = "dynamic"
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
	// +docsgen:workspaceMode
	// The mode of the workspace. This can be either `fixed` or `dynamic`.
	// In `fixed` mode, the current workspace used at runtime is always the one set in the currentWorkspace config field.
	// In `dynamic` mode, the current workspace used at runtime is determined by the current directory.
	// If the current directory is within a workspace, that workspace is used.
	WorkspaceMode WorkspaceMode `json:"workspaceMode" yaml:"workspaceMode"`
	// +docsgen:currentNamespace
	// The name of the current namespace. This is not required to be set.
	CurrentNamespace string `json:"currentNamespace" yaml:"currentNamespace"`
	// +docsgen:interactive
	// Configurations for the interactive UI.
	Interactive *InteractiveConfig `json:"interactive" yaml:"interactive"`
	// +docsgen:templates
	// A map of executable definition template names to their paths.
	Templates map[string]string `json:"templates,omitempty" yaml:"templates,omitempty"`
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
	if c.WorkspaceMode != "" && c.WorkspaceMode != WorkspaceModeFixed && c.WorkspaceMode != WorkspaceModeDynamic {
		return fmt.Errorf("invalid workspace mode %s", c.WorkspaceMode)
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

func (c *UserConfig) JSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal user config - %w", err)
	}
	return string(jsonBytes), nil
}

func (c *UserConfig) Markdown() string {
	mkdwn := "# Global Configurations\n"
	mkdwn += fmt.Sprintf("## Current workspace\n%s\n", c.CurrentWorkspace)
	if c.WorkspaceMode == WorkspaceModeFixed {
		mkdwn += "**Workspace mode is set to fixed. This means that your working directory will have no impact on the " +
			"current workspace.**\n"
	} else if c.WorkspaceMode == WorkspaceModeDynamic {
		mkdwn += "**Workspace mode is set to dynamic. This means that your current workspace is also determined by " +
			"your working directory.**\n"
	}

	if c.CurrentNamespace != "" {
		mkdwn += fmt.Sprintf("## Current namespace\n%s\n", c.CurrentNamespace)
	}
	if c.Interactive != nil {
		interactiveConfig, err := yaml.Marshal(c.Interactive)
		if err != nil {
			mkdwn += "## Interactive UI config\nerror\n"
		} else {
			mkdwn += fmt.Sprintf("## Git UI config\n```yaml\n%s```\n", string(interactiveConfig))
		}
	}
	mkdwn += "## Registered Workspaces\n"
	for name, path := range c.Workspaces {
		mkdwn += fmt.Sprintf("- %s: %s\n", name, path)
	}

	if len(c.Templates) > 0 {
		mkdwn += "## Registered Templates\n"
		for name, path := range c.Templates {
			mkdwn += fmt.Sprintf("- %s: %s\n", name, path)
		}
	}

	return mkdwn
}
