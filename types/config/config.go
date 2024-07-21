package config

import (
	"encoding/json"
	"fmt"

	tuikitIO "github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/utils"
)

//go:generate go run github.com/atombender/go-jsonschema@v0.16.0 -et --only-models -p config -o config.gen.go --capitalization URL  schema.yaml

func (c *Config) Validate() error {
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
	if c.WorkspaceMode != "" &&
		c.WorkspaceMode != ConfigWorkspaceModeFixed &&
		c.WorkspaceMode != ConfigWorkspaceModeDynamic {
		return fmt.Errorf("invalid workspace mode %s", c.WorkspaceMode)
	}
	for _, ws := range c.RemoteWorkspaces {
		if err := ws.Validate(); err != nil {
			return err
		}
	}
	if err := c.DefaultLogMode.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *Config) SetDefaults() {
	if c.Workspaces == nil {
		c.Workspaces = make(map[string]string)
	}
	if c.RemoteWorkspaces == nil {
		c.RemoteWorkspaces = make(map[string]RemoteWorkspace)
	}
	if c.CurrentWorkspace == "" && len(c.Workspaces) > 0 {
		c.CurrentWorkspace = maps.Keys(c.Workspaces)[0]
	}
	if c.WorkspaceMode == "" {
		c.WorkspaceMode = ConfigWorkspaceModeDynamic
	}
	if c.DefaultLogMode == "" {
		c.DefaultLogMode = tuikitIO.Logfmt
	}
}

func (c *Config) ShowTUI() bool {
	return c.Interactive != nil && c.Interactive.Enabled
}

func (c *Config) SendTextNotification() bool {
	return c.Interactive != nil && c.Interactive.Enabled &&
		c.Interactive.NotifyOnCompletion != nil && *c.Interactive.NotifyOnCompletion
}

func (c *Config) SendSoundNotification() bool {
	return c.Interactive != nil && c.Interactive.Enabled &&
		c.Interactive.SoundOnCompletion != nil && *c.Interactive.SoundOnCompletion
}

func (c *Config) YAML() (string, error) {
	yamlBytes, err := yaml.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user config - %w", err)
	}
	return string(yamlBytes), nil
}

func (c *Config) JSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal user config - %w", err)
	}
	return string(jsonBytes), nil
}

func (c *Config) Markdown() string {
	mkdwn := "# Global Configurations\n"
	mkdwn += fmt.Sprintf("## Current workspace\n%s\n", c.CurrentWorkspace)
	if c.WorkspaceMode == ConfigWorkspaceModeFixed {
		mkdwn += "**Workspace mode is set to fixed. This means that your working directory will have no impact on the " +
			"current workspace.**\n"
	} else if c.WorkspaceMode == ConfigWorkspaceModeDynamic {
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
			mkdwn += fmt.Sprintf("## Interactive UI config\n```yaml\n%s```\n", string(interactiveConfig))
		}
	}
	mkdwn += "## Registered Local Workspaces\n"
	for name, path := range c.Workspaces {
		mkdwn += fmt.Sprintf("- %s: %s\n", name, path)
	}

	if len(c.RemoteWorkspaces) > 0 {
		mkdwn += "## Registered Remote Workspaces\n"
		for name, ws := range c.RemoteWorkspaces {
			mkdwn += fmt.Sprintf("- name: %s\n", name)
			mkdwn += fmt.Sprintf("  url: %s\n", ws.URL)
			switch {
			case ws.Branch != "":
				mkdwn += fmt.Sprintf("  branch: %s\n", ws.Branch)
			case ws.Commit != "":
				mkdwn += fmt.Sprintf("  commit: %s\n", ws.Commit)
			case ws.Tag != "":
				mkdwn += fmt.Sprintf("  tag: %s\n", ws.Tag)
			}
			mkdwn += fmt.Sprintf("  pullOnSync: %s\n", ws.PullOnSync)
		}

	}

	if len(c.Templates) > 0 {
		mkdwn += "## Registered Templates\n"
		for name, path := range c.Templates {
			mkdwn += fmt.Sprintf("- %s: %s\n", name, path)
		}
	}

	return mkdwn
}

func (l RemoteWorkspace) Validate() error {
	if l.URL == "" {
		return fmt.Errorf("remote workspace URL is required")
	}
	if err := utils.ValidateOneOf("git version identifier", l.Branch, l.Commit, l.Tag); err != nil {
		return errors.Wrap(err, fmt.Sprintf("remote workspace %s", l.URL))
	}
	return nil
}
