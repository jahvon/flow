package config

import (
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/utils"
)

type WorkspaceConfig struct {
	DisplayName string                    `json:"displayName"           yaml:"displayName"`
	Description string                    `json:"description,omitempty" yaml:"description,omitempty"`
	Tags        Tags                      `json:"tags,omitempty"        yaml:"tags,omitempty"`
	Git         *GitConfig                `json:"git,omitempty"         yaml:"git,omitempty"`
	Executables *ExecutableLocationConfig `json:"executables,omitempty" yaml:"executables,omitempty"`

	assignedName string
	location     string
}

type WorkspaceConfigList []WorkspaceConfig

type GitConfig struct {
	Enabled    bool `json:"enabled"              yaml:"enabled"`
	PullOnSync bool `json:"pullOnSync,omitempty" yaml:"pullOnSync,omitempty"`
}

type ExecutableLocationConfig struct {
	Included []string `json:"included,omitempty" yaml:"included,omitempty"`
	Excluded []string `json:"excluded,omitempty" yaml:"excluded,omitempty"`
}

type enrichedWorkspaceConfigList struct {
	Workspaces WorkspaceConfigList `json:"workspaces" yaml:"workspaces"`
}

func (c *WorkspaceConfig) AssignedName() string {
	return c.assignedName
}

func (c *WorkspaceConfig) Location() string {
	return c.location
}

func (c *WorkspaceConfig) SetContext(name, location string) {
	c.assignedName = name
	c.location = location
}

func (c *WorkspaceConfig) YAML() (string, error) {
	yamlBytes, err := yaml.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal workspace config - %w", err)
	}
	return string(yamlBytes), nil
}

func (c *WorkspaceConfig) JSON(pretty bool) (string, error) {
	var jsonBytes []byte
	var err error
	if pretty {
		jsonBytes, err = json.MarshalIndent(c, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(c)
	}
	if err != nil {
		return "", fmt.Errorf("failed to marshal workspace config - %w", err)
	}
	return string(jsonBytes), nil
}

func (c *WorkspaceConfig) Markdown() string {
	var mkdwn string
	if c.DisplayName != "" {
		mkdwn = fmt.Sprintf("# [Workspace] %s\n", c.DisplayName)
	} else {
		mkdwn = fmt.Sprintf("# [Workspace] %s\n", c.AssignedName())
	}

	mkdwn += fmt.Sprintf("## Location\n%s\n", c.Location())
	if c.Description != "" {
		mkdwn += fmt.Sprintf("## Description\n%s\n", c.Description)
	}
	if len(c.Tags) > 0 {
		mkdwn += "## Tags\n"
		lo.ForEach(c.Tags, func(tag string, _ int) {
			mkdwn += fmt.Sprintf("- %s\n", tag)
		})
	}
	if c.Git != nil {
		gitConfig, err := yaml.Marshal(c.Git)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal git config")
			mkdwn += "## Git config\nerror\n"
		} else {
			mkdwn += fmt.Sprintf("## Git config\n```yaml\n%s```\n", string(gitConfig))
		}
	}
	if c.Executables != nil {
		execs, err := yaml.Marshal(c.Executables)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal executables")
			mkdwn += "## Executables\nerror\n"
		} else {
			mkdwn += fmt.Sprintf("## Executables\n```yaml\n%s```\n", string(execs))
		}
	}
	return mkdwn
}

func (l WorkspaceConfigList) YAML() (string, error) {
	enriched := enrichedWorkspaceConfigList{Workspaces: l}
	yamlBytes, err := yaml.Marshal(enriched)
	if err != nil {
		return "", fmt.Errorf("failed to marshal workspace config list - %w", err)
	}
	return string(yamlBytes), nil
}

func (l WorkspaceConfigList) JSON(pretty bool) (string, error) {
	var jsonBytes []byte
	var err error
	enriched := enrichedWorkspaceConfigList{Workspaces: l}
	if pretty {
		jsonBytes, err = json.MarshalIndent(enriched, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(enriched)
	}
	if err != nil {
		return "", fmt.Errorf("failed to marshal workspace config list - %w", err)
	}
	return string(jsonBytes), nil
}

func (l WorkspaceConfigList) Items() []CollectionItem {
	items := make([]CollectionItem, 0)
	for _, ws := range l {
		name := ws.AssignedName()
		if ws.DisplayName != "" {
			name = ws.DisplayName
		}

		var location string
		if ws.Location() == "" {
			location = "unk"
		} else {
			location = utils.PathFromWd(ws.Location())
		}

		item := CollectionItem{
			Header:      name,
			SubHeader:   location,
			Description: ws.Description,
			Tags:        ws.Tags,
		}
		items = append(items, item)
	}
	return items
}

func (l WorkspaceConfigList) Singular() string {
	return "workspace"
}

func (l WorkspaceConfigList) Plural() string {
	return "workspaces"
}
