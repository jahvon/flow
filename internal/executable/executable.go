package executable

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/jahvon/flow/internal/errors"
	"github.com/jahvon/flow/internal/executable/consts"
	"github.com/jahvon/flow/internal/io"
)

var log = io.Log().With().Str("pkg", "executable").Logger()

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
		return fmt.Sprintf("%s/%s", e.workspace, e.Name)
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
	exec, found := lo.Find(*l, func(exec *Executable) bool {
		return (exec.Type == agent && exec.Name == name) || lo.Contains(exec.Aliases, name)
	})
	if found {
		return exec, nil
	}
	return nil, errors.ExecutableNotFoundError{Agent: agent, Name: name}
}

func (l *List) FilterByTags(tags []string) List {
	execs := lo.Filter(*l, func(exec *Executable, _ int) bool {
		return lo.Some(exec.Tags, tags)
	})
	log.Trace().Int("executables", len(execs)).Msgf("filtered executables by tags %v", tags)
	return execs
}

func (l *List) FilterByTag(tag string) List {
	execs := lo.Filter(*l, func(exec *Executable, _ int) bool {
		return lo.Contains(exec.Tags, tag)
	})
	log.Trace().Int("executables", len(execs)).Msgf("filtered executables by tag %s", tag)
	return execs
}

func (l *List) FilterByType(agent consts.AgentType) List {
	execs := lo.Filter(*l, func(exec *Executable, _ int) bool {
		return exec.Type == agent
	})
	log.Trace().Int("executables", len(execs)).Msgf("filtered executables by type %s", agent)
	return execs
}
