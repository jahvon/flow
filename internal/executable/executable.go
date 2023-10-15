package executable

import (
	"fmt"

	"github.com/jahvon/flow/internal/errors"
	"github.com/jahvon/flow/internal/executable/consts"
	"github.com/jahvon/flow/internal/io"
)

var log = io.Log()

type Preference map[string]interface{}

type Executable struct {
	Type        consts.AgentType       `yaml:"type"`
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Aliases     []string               `yaml:"aliases"`
	Tags        []string               `yaml:"tags"`
	Spec        map[string]interface{} `yaml:"spec"`

	workspace     string
	workspacePath string
	namespace     string
}

func (e *Executable) SetContext(workspace, workspacePath, namespace string) {
	e.workspace = workspace
	e.workspacePath = workspacePath
	e.namespace = namespace
}

func (e *Executable) GetContext() (workspace, workspacePath, namespace string) {
	return e.workspace, e.workspacePath, e.namespace
}

func (e *Executable) ID() string {
	if e.namespace == "" {
		return fmt.Sprintf("%s:%s", e.workspace, e.Name)
	} else if e.workspace == "" {
		log.Debug().
			Str("workspace", e.workspace).Str("namespace", e.namespace).
			Msgf("missing context for %s", e.Name)
		return "missing-context"
	}

	return fmt.Sprintf("%s/%s:%s", e.workspace, e.namespace, e.Name)
}

func (e *Executable) Validate() error {
	if e.Type == "" {
		return fmt.Errorf("type cannot be empty")
	}

	if e.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if e.Spec == nil {
		return fmt.Errorf("spec cannot be empty")
	}

	return nil
}

func (e *Executable) MergeTags(tags []string) {
	// TODO: dedupe tags
	e.Tags = append(e.Tags, tags...)
}

type List []*Executable

func (l *List) FindByTypeAndName(agent consts.AgentType, name string) (*Executable, error) {
	for _, exec := range *l {
		if exec.Type != agent {
			continue
		}

		if exec.Name == name {
			return exec, nil
		} else if len(exec.Aliases) > 0 {
			for _, alias := range exec.Aliases {
				if alias == name {
					return exec, nil
				}
			}
		}
	}
	return nil, errors.ExecutableNotFound(agent, name)
}

func (l *List) FilterByTag(tag string) []*Executable {
	var executables []*Executable
	for _, exec := range *l {
		for _, execTag := range exec.Tags {
			if execTag == tag {
				executables = append(executables, exec)
			}
		}
	}
	return executables
}

func (l *List) FilterByType(agent consts.AgentType) []*Executable {
	var executables []*Executable
	for _, exec := range *l {
		if exec.Type == agent {
			executables = append(executables, exec)
		}
	}
	return executables
}
