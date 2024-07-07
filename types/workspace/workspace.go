package workspace

import (
	"encoding/json"
	"fmt"

	"github.com/jahvon/tuikit/types"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/utils"
	"github.com/jahvon/flow/types/common"
)

//go:generate go run github.com/atombender/go-jsonschema@v0.16.0 -et --only-models -p workspace -o workspace.gen.go schema.yaml

type WorkspaceList []Workspace

type enrichedWorkspaceConfigList struct {
	Workspaces WorkspaceList `json:"workspaces" yaml:"workspaces"`
}

func (w *Workspace) AssignedName() string {
	return w.assignedName
}

func (w *Workspace) Location() string {
	return w.location
}

func (w *Workspace) SetContext(name, location string) {
	w.assignedName = name
	w.location = location
}

func (w *Workspace) YAML() (string, error) {
	yamlBytes, err := yaml.Marshal(w)
	if err != nil {
		return "", fmt.Errorf("failed to marshal workspace config - %w", err)
	}
	return string(yamlBytes), nil
}

func (w *Workspace) JSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal workspace config - %w", err)
	}
	return string(jsonBytes), nil
}

func (w *Workspace) Markdown() string {
	var mkdwn string
	if w.DisplayName != "" {
		mkdwn = fmt.Sprintf("# [Workspace] %s\n", w.DisplayName)
	} else {
		mkdwn = fmt.Sprintf("# [Workspace] %s\n", w.AssignedName())
	}

	mkdwn += fmt.Sprintf("## Location\n%s\n", w.Location())
	if w.Description != "" {
		mkdwn += fmt.Sprintf("## Description\n%s\n", w.Description)
	}
	if w.Tags != nil && len(w.Tags) > 0 {
		mkdwn += "## Tags\n"
		for _, tag := range w.Tags {
			mkdwn += fmt.Sprintf("- %s\n", tag)
		}
	}
	if w.Executables != nil {
		execs, err := yaml.Marshal(w.Executables)
		if err != nil {
			mkdwn += "## Executables\nerror\n"
		} else {
			mkdwn += fmt.Sprintf("## Executables\n```yaml\n%s```\n", string(execs))
		}
	}
	return mkdwn
}

func DefaultWorkspaceConfig(name string) *Workspace {
	return &Workspace{DisplayName: name}
}

func (l WorkspaceList) YAML() (string, error) {
	enriched := enrichedWorkspaceConfigList{Workspaces: l}
	yamlBytes, err := yaml.Marshal(enriched)
	if err != nil {
		return "", fmt.Errorf("failed to marshal workspace config list - %w", err)
	}
	return string(yamlBytes), nil
}

func (l WorkspaceList) JSON() (string, error) {
	enriched := enrichedWorkspaceConfigList{Workspaces: l}
	jsonBytes, err := json.MarshalIndent(enriched, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal workspace config list - %w", err)
	}
	return string(jsonBytes), nil
}

func (l WorkspaceList) FindByName(name string) *Workspace {
	for _, ws := range l {
		if ws.AssignedName() == name {
			return &ws
		}
	}
	return nil
}

func (l WorkspaceList) Items() []*types.CollectionItem {
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
		if ws.Tags != nil && len(ws.Tags) > 0 {
			tags := common.Tags(ws.Tags)
			d := fmt.Sprintf("[%s]\n", tags.PreviewString()) + ws.Description
			ws.Description = d
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

func (l WorkspaceList) Singular() string {
	return "workspace"
}

func (l WorkspaceList) Plural() string {
	return "workspaces"
}
