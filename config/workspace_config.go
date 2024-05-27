package config

import (
	"encoding/json"
	"fmt"

	"github.com/jahvon/tuikit/types"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/utils"
)

type WorkspaceConfig struct {
	// +docsgen:displayName
	// The display name of the workspace. This is used in the interactive UI.
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
	Enabled bool `json:"enabled" yaml:"enabled"`

	// +docsgen:pullOnSync
	// Whether to pull the latest changes from the remote when syncing.
	PullOnSync bool `json:"pullOnSync,omitempty" yaml:"pullOnSync,omitempty"`
}

type ExecutableLocationConfig struct {
	// +docsgen:included
	// A list of directories to include in the executable search.
	Included []string `json:"included,omitempty" yaml:"included,omitempty"`
	// +docsgen:excluded
	// A list of directories to exclude from the executable search.
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

func (c *WorkspaceConfig) JSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
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
			mkdwn += "## Git config\nerror\n"
		} else {
			mkdwn += fmt.Sprintf("## Git config\n```yaml\n%s```\n", string(gitConfig))
		}
	}
	if c.Executables != nil {
		execs, err := yaml.Marshal(c.Executables)
		if err != nil {
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

func (l WorkspaceConfigList) JSON() (string, error) {
	enriched := enrichedWorkspaceConfigList{Workspaces: l}
	jsonBytes, err := json.MarshalIndent(enriched, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal workspace config list - %w", err)
	}
	return string(jsonBytes), nil
}

func (l WorkspaceConfigList) Items() []*types.CollectionItem {
	items := make([]*types.CollectionItem, 0)
	for _, ws := range l {
		name := ws.AssignedName()
		if ws.DisplayName != "" {
			name = ws.DisplayName
		}

		var location string
		if ws.Location() == "" {
			location = "unk"
		} else {
			var err error
			location, err = utils.PathFromWd(ws.Location())
			if err != nil {
				location = ws.Location()
			}
		}
		if len(ws.Tags) > 0 {
			ws.Description = fmt.Sprintf("[%s]\n", ws.Tags.PreviewString()) + ws.Description
		}

		item := types.CollectionItem{
			Header:    name,
			SubHeader: location,
			Desc:      ws.Description,
			ID:        name,
		}
		items = append(items, &item)
	}
	return items
}

func (l WorkspaceConfigList) Singular() string {
	return "workspace"
}

func (l WorkspaceConfigList) Plural() string {
	return "workspaces"
}
