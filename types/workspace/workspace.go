package workspace

import (
	"encoding/json"
	"fmt"

	"github.com/jahvon/tuikit/types"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/utils"
	"github.com/jahvon/flow/types/common"
)

//go:generate go tool go-jsonschema -et --only-models -p workspace -o workspace.gen.go schema.yaml

type WorkspaceList []*Workspace

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
	return workspaceMarkdown(w)
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
			return ws
		}
	}
	return nil
}

func (l WorkspaceList) Items() []*types.EntityInfo {
	items := make([]*types.EntityInfo, 0)
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
			tags := common.Tags(ws.Tags)
			d := fmt.Sprintf("[%s]\n", tags.PreviewString()) + ws.Description
			ws.Description = d
		}

		item := types.EntityInfo{
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
