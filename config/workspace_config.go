package config

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/utils"
)

type WorkspaceConfig struct {
	DisplayName string             `json:"displayName"           yaml:"displayName"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
	Tags        Tags               `json:"tags,omitempty"        yaml:"tags,omitempty"`
	Git         *GitConfig         `json:"git,omitempty"         yaml:"git,omitempty"`
	Executables *ExecutablesConfig `json:"executables,omitempty" yaml:"executables,omitempty"`

	assignedName string
	location     string
}

type WorkspaceConfigList []WorkspaceConfig

type GitConfig struct {
	Enabled    bool `json:"enabled"              yaml:"enabled"`
	PullOnSync bool `json:"pullOnSync,omitempty" yaml:"pullOnSync,omitempty"`
}

type ExecutablesConfig struct {
	Included []string `json:"included,omitempty" yaml:"included,omitempty"`
	Excluded []string `json:"excluded,omitempty" yaml:"excluded,omitempty"`
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
		return "", fmt.Errorf("failed to marshal workspace config - %v", err)
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
		return "", fmt.Errorf("failed to marshal workspace config - %v", err)
	}
	return string(jsonBytes), nil
}

func (c *WorkspaceConfig) Map() map[string]string {
	fields := make(map[string]string)
	fields["Name"] = c.AssignedName()
	if c.AssignedName() != c.DisplayName && c.DisplayName != "" {
		fields["Display name"] = c.DisplayName
	}
	fields["Location"] = c.Location()
	if c.Description != "" {
		fields["Description"] = utils.WrapLines(c.Description, 20)
	}
	if c.Tags != nil {
		fields["Tags"] = c.Tags.String()
	}
	if c.Git != nil {
		gitConfig, err := yaml.Marshal(c.Git)
		if err != nil {
			fields["Git config"] = "error"
			log.Error().Err(err).Msg("failed to marshal git config")
		} else {
			fields["Git config"] = strings.TrimSpace(string(gitConfig))
		}
	}
	if c.Executables != nil {
		execs, err := yaml.Marshal(c.Executables)
		if err != nil {
			fields["Executables"] = "error"
			log.Error().Err(err).Msg("failed to marshal executables")
		} else {
			fields["Executables"] = string(execs)
		}
	}
	return fields
}

func (c *WorkspaceConfig) DetailsString() string {
	fields := c.Map()
	return utils.DetailsString(fields)
}

func (l WorkspaceConfigList) YAML() (string, error) {
	yamlBytes, err := yaml.Marshal(l)
	if err != nil {
		return "", fmt.Errorf("failed to marshal workspace config list - %v", err)
	}
	return string(yamlBytes), nil
}

func (l WorkspaceConfigList) JSON(pretty bool) (string, error) {
	var jsonBytes []byte
	var err error
	if pretty {
		jsonBytes, err = json.MarshalIndent(l, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(l)
	}
	if err != nil {
		return "", fmt.Errorf("failed to marshal workspace config list - %v", err)
	}
	return string(jsonBytes), nil
}

func (l WorkspaceConfigList) TableData() (header []string, rows [][]string) {
	header = []string{"Name", "Tags", "Description"}
	for _, ws := range l {
		name := ws.AssignedName()
		if ws.DisplayName != "" {
			name = ws.DisplayName
		}
		rows = append(
			rows,
			[]string{
				name,
				ws.Tags.PreviewString(),
				utils.ShortenString(ws.Description, 80),
			})
	}
	sort.Slice(rows, func(i, j int) bool {
		return strings.Compare(rows[i][0], rows[j][0]) < 0
	})
	return
}
