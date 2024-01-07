package config

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/errors"
	"github.com/jahvon/flow/internal/utils"
)

const tmpdir = "f:tmp"

type DirectoryScopedExecutable struct {
	Directory string `yaml:"dir"`

	finalizerFunc func()
}

func (e *DirectoryScopedExecutable) ExpandDirectory(wsPath, execPath string, env map[string]string) (string, error) {
	if e.Directory == tmpdir {
		file, err := os.CreateTemp("", "flow")
		if err != nil {
			return "", err
		}
		e.finalizerFunc = func() {
			if err := os.Remove(file.Name()); err != nil {
				log.Error().Err(err).Msgf("unable to remove temp file %s", file.Name())
			}
		}
		return file.Name(), nil
	}

	return utils.ExpandDirectory(e.Directory, wsPath, execPath, env), nil
}

func (e *DirectoryScopedExecutable) Finalize() {
	if e.finalizerFunc != nil {
		e.finalizerFunc()
	}
}

type ParameterizedExecutable struct {
	Parameters ParameterList `yaml:"params"`
}

type ExecExecutableType struct {
	DirectoryScopedExecutable `yaml:",inline"`
	ParameterizedExecutable   `yaml:",inline"`

	Command string  `yaml:"cmd"`
	File    string  `yaml:"file"`
	LogMode LogMode `yaml:"logMode"`

	logFields map[string]interface{}
}

func (e *ExecExecutableType) SetLogFields(fields map[string]interface{}) {
	e.logFields = fields
}

func (e *ExecExecutableType) GetLogFields() map[string]interface{} {
	return e.logFields
}

type LaunchExecutableType struct {
	ParameterizedExecutable `yaml:",inline"`

	App  string `yaml:"app"`
	URI  string `yaml:"uri"`
	Wait bool   `yaml:"wait"`
}

type SerialExecutableType struct {
	ParameterizedExecutable `yaml:",inline"`

	ExecutableRefs []Ref `yaml:"refs"`
	FailFast       bool  `yaml:"failFast"`
}

type ParallelExecutableType struct {
	ParameterizedExecutable `yaml:",inline"`

	ExecutableRefs []Ref `yaml:"refs"`
	MaxThreads     int   `yaml:"maxThreads"`
	FailFast       bool  `yaml:"failFast"`
}

type ExecutableTypeSpec struct {
	Exec     *ExecExecutableType     `yaml:"exec,omitempty"`
	Launch   *LaunchExecutableType   `yaml:"launch,omitempty"`
	Serial   *SerialExecutableType   `yaml:"serial,omitempty"`
	Parallel *ParallelExecutableType `yaml:"parallel,omitempty"`
}

type Executable struct {
	Verb        Verb                `yaml:"verb"`
	Name        string              `yaml:"name"`
	Aliases     []string            `yaml:"aliases"`
	Tags        Tags                `yaml:"tags"`
	Description string              `yaml:"description"`
	Visibility  VisibilityType      `yaml:"visibility"`
	Timeout     time.Duration       `yaml:"timeout"`
	Type        *ExecutableTypeSpec `yaml:"type"`

	workspace, namespace, workspacePath, definitionPath string
}

type ExecutableList []*Executable

type enrichedExecutableList struct {
	Executables []*enrichedExecutable `json:"executables" yaml:"executables"`
}
type enrichedExecutable struct {
	ID   string      `json:"id"   yaml:"id"`
	Spec *Executable `json:"spec" yaml:"spec"`
}

func (e *Executable) SetContext(workspaceName, workspacePath, namespace, definitionPath string) {
	e.workspace = workspaceName
	e.workspacePath = workspacePath
	e.namespace = namespace
	e.definitionPath = definitionPath
}

func (e *Executable) YAML() (string, error) {
	enriched := &enrichedExecutable{
		ID:   e.ID(),
		Spec: e,
	}
	yamlBytes, err := yaml.Marshal(enriched)
	if err != nil {
		return "", fmt.Errorf("failed to marshal executable - %w", err)
	}
	return string(yamlBytes), nil
}

func (e *Executable) JSON(pretty bool) (string, error) {
	var jsonBytes []byte
	var err error
	enriched := &enrichedExecutable{
		ID:   e.ID(),
		Spec: e,
	}
	if pretty {
		jsonBytes, err = json.MarshalIndent(enriched, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(enriched)
	}
	if err != nil {
		return "", fmt.Errorf("failed to marshal executable - %w", err)
	}
	return string(jsonBytes), nil
}

func (e *Executable) Markdown() string {
	var mkdwn string
	mkdwn += fmt.Sprintf("# [Executable] %s %s\n", e.Verb, e.ID())
	mkdwn += fmt.Sprintf("## Defined in\n%s\n", e.definitionPath)
	if len(e.Aliases) > 0 {
		mkdwn += "## Aliases\n"
		lo.ForEach(e.Aliases, func(alias string, _ int) {
			mkdwn += fmt.Sprintf("- %s\n", alias)
		})
	}
	if e.Description != "" {
		mkdwn += fmt.Sprintf("## Description\n%s\n", e.Description)
	}
	if len(e.Tags) > 0 {
		mkdwn += "## Tags\n"
		lo.ForEach(e.Tags, func(tag string, _ int) {
			mkdwn += fmt.Sprintf("- %s\n", tag)
		})
	}
	if e.Type != nil {
		typeSpec, err := yaml.Marshal(e.Type)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal type spec")
			mkdwn += "## Type spec\nerror\n"
		} else {
			mkdwn += fmt.Sprintf("## Type spec\n```yaml\n%s```\n", string(typeSpec))
		}
	}
	if e.Visibility != "" {
		mkdwn += fmt.Sprintf("## Visibility\n%s\n", e.Visibility)
	}
	if e.Timeout != 0 {
		mkdwn += fmt.Sprintf("## Timeout\n%s\n", e.Timeout.String())
	}
	return mkdwn
}

func (e *Executable) Ref() Ref {
	return Ref(fmt.Sprintf("%s %s", e.Verb, e.ID()))
}

func (e *Executable) ID() string {
	if e.workspace == "" {
		log.Debug().
			Str("namespace", e.namespace).
			Str("definitionFile", e.definitionPath).
			Msgf("missing workspace for %s", e.Name)
		return "unk"
	}

	return NewExecutableID(e.workspace, e.namespace, e.Name)
}

func (e *Executable) AliasesIDs() []string {
	if len(e.Aliases) == 0 {
		return nil
	}

	if e.workspace == "" {
		log.Debug().
			Str("namespace", e.namespace).
			Str("definitionFile", e.definitionPath).
			Msgf("missing workspace for %s", e.Name)
		return nil
	}
	aliases := make([]string, 0)
	for _, alias := range e.Aliases {
		aliases = append(aliases, NewExecutableID(e.workspace, e.namespace, alias))
	}
	return aliases
}

func (e *Executable) WorkspacePath() string {
	return e.workspacePath
}

func (e *Executable) DefinitionPath() string {
	return e.definitionPath
}

func (e *Executable) SetDefaults() {
	if e.Verb == "" {
		e.Verb = "exec"
	}
	if e.Visibility == "" {
		e.Visibility = VisibilityPrivate
	}
	if e.Timeout == 0 {
		e.Timeout = DefaultTimeout
	}

	if e.Type != nil && e.Type.Exec != nil && e.Type.Exec.LogMode == "" {
		e.Type.Exec.LogMode = StructuredLogMode
	}
}

func (e *Executable) Validate() error {
	if e.Verb == "" {
		return fmt.Errorf("verb cannot be empty")
	}
	if e.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if e.Type == nil {
		return fmt.Errorf("type cannot be empty")
	}
	err := utils.ValidateOneOf(
		"executable type",
		e.Type.Exec,
		e.Type.Launch,
		e.Type.Serial,
		e.Type.Parallel,
	)
	if err != nil {
		return err
	}

	if e.workspace == "" {
		return fmt.Errorf("workspace was not set")
	}
	if e.namespace == "" {
		return fmt.Errorf("namespace was not set")
	}
	if e.definitionPath == "" {
		return fmt.Errorf("definition path was not set")
	}

	return nil
}

func (e *Executable) NameEquals(name string) bool {
	return e.Name == name || lo.Contains(e.Aliases, name)
}

func (e *Executable) MergeTags(tags Tags) {
	e.Tags = lo.Uniq(append(e.Tags, tags...))
}

func (e *Executable) MergeVisibility(visibility VisibilityType) {
	curLevel := slices.Index(visibilityByLevel, e.Visibility)
	vLevel := slices.Index(visibilityByLevel, visibility)
	if vLevel > curLevel {
		e.Visibility = visibility
	}
}

// IsVisibleFromWorkspace returns true if the executable should be shown in terminal output for the given workspace.
func (e *Executable) IsVisibleFromWorkspace(workspaceFilter string) bool {
	switch e.Visibility {
	case VisibilityPrivate:
		return e.workspace == workspaceFilter || workspaceFilter == "" || workspaceFilter == "*"
	case VisibilityPublic:
		return true
	case VisibilityInternal, VisibilityHidden:
		return false
	default:
		return false
	}
}

// IsExecutableFromWorkspace returns true if the executable can be executed from the given workspace.
func (e *Executable) IsExecutableFromWorkspace(workspace string) bool {
	switch e.Visibility {
	case VisibilityPrivate, VisibilityInternal:
		return e.workspace == workspace
	case VisibilityPublic:
		return true
	case VisibilityHidden:
		return false
	default:
		return false
	}
}

func (l ExecutableList) YAML() (string, error) {
	enriched := &enrichedExecutableList{}
	for _, exec := range l {
		enriched.Executables = append(enriched.Executables, &enrichedExecutable{
			ID:   exec.ID(),
			Spec: exec,
		})
	}
	yamlBytes, err := yaml.Marshal(enriched)
	if err != nil {
		return "", fmt.Errorf("failed to marshal executable list - %w", err)
	}
	return string(yamlBytes), nil
}

func (l ExecutableList) JSON(pretty bool) (string, error) {
	var jsonBytes []byte
	var err error
	enriched := &enrichedExecutableList{}
	for _, exec := range l {
		enriched.Executables = append(enriched.Executables, &enrichedExecutable{
			ID:   exec.ID(),
			Spec: exec,
		})
	}
	if pretty {
		jsonBytes, err = json.MarshalIndent(enriched, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(enriched)
	}
	if err != nil {
		return "", fmt.Errorf("failed to marshal executable list - %w", err)
	}
	return string(jsonBytes), nil
}

func (l ExecutableList) Items() []CollectionItem {
	items := make([]CollectionItem, 0)
	for _, exec := range l {
		item := CollectionItem{
			Header:      exec.ID(),
			SubHeader:   exec.Verb.String(),
			Description: exec.Description,
			Tags:        exec.Tags,
		}
		items = append(items, item)
	}
	return items
}

func (l ExecutableList) Singular() string {
	return "executable"
}

func (l ExecutableList) Plural() string {
	return "executables"
}

func (l ExecutableList) FindByVerbAndID(verb Verb, id string) (*Executable, error) {
	_, _, name := ParseExecutableID(id) // Assumes that ws and ns has already been filtered down
	if name == "" {
		return nil, errors.ExecutableNotFoundError{Verb: string(verb), Name: name}
	}
	filteredList := l.FilterByVerb(verb)
	exec, found := lo.Find(filteredList, func(exec *Executable) bool {
		return exec.NameEquals(name)
	})
	if found {
		return exec, nil
	}
	return nil, errors.ExecutableNotFoundError{Verb: string(verb), Name: name}
}

func (l ExecutableList) FilterByTags(tags Tags) ExecutableList {
	if len(tags) == 0 {
		return l
	}

	execs := lo.Filter(l, func(exec *Executable, _ int) bool {
		return exec.Tags.HasAnyTag(tags)
	})
	log.Trace().Int("executables", len(execs)).Msgf("filtered executables by tags %v", tags)
	return execs
}

func (l ExecutableList) FilterByVerb(verb Verb) ExecutableList {
	if verb == "" || verb == "*" {
		return l
	}

	if err := verb.Validate(); err != nil {
		return ExecutableList{}
	}

	execs := lo.Filter(l, func(exec *Executable, _ int) bool {
		return exec.Verb.Equals(verb)
	})

	log.Trace().Int("executables", len(execs)).Msgf("filtered executables by verb %s", verb)
	return execs
}

func (l ExecutableList) FilterByWorkspace(ws string) ExecutableList {
	executables := lo.Filter(l, func(executable *Executable, _ int) bool {
		return executable.IsVisibleFromWorkspace(ws)
	})

	if ws == "" || ws == "*" {
		return executables
	}

	executables = lo.Filter(executables, func(executable *Executable, _ int) bool {
		return executable.workspace == ws
	})

	log.Trace().Int("executables", len(executables)).Msgf("filtered executables by workspace %s", ws)
	return executables
}

func (l ExecutableList) FilterByNamespace(ns string) ExecutableList {
	if ns == "" || ns == "*" {
		return l
	}

	executables := lo.Filter(l, func(executable *Executable, _ int) bool {
		return executable.namespace == ns
	})

	log.Trace().Int("executables", len(executables)).Msgf("filtered executables by namespace %s", ns)
	return executables
}

func ParseExecutableID(id string) (workspace, namespace, name string) {
	parts := strings.Split(id, "/")
	switch len(parts) {
	case 1:
		return "", "", parts[0]
	case 2:
		subparts := strings.Split(parts[1], ":")
		if len(subparts) == 1 {
			return parts[0], "*", subparts[0]
		} else if len(subparts) == 2 {
			return parts[0], subparts[0], subparts[1]
		}
	}

	return "", "", ""
}

func NewExecutableID(workspace, namespace, name string) string {
	if namespace == "" {
		return fmt.Sprintf("%s/%s", workspace, name)
	}
	return fmt.Sprintf("%s/%s:%s", workspace, namespace, name)
}
