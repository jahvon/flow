package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/jahvon/tuikit/io"
	"github.com/jahvon/tuikit/types"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/errors"
	"github.com/jahvon/flow/internal/utils"
)

const tmpdir = "f:tmp"

type ExecutableDirectory struct {
	// +docsgen:dir
	// The directory to execute the command in.
	// If unset, the directory of the executable definition will be used.
	// If set to `f:tmp`, a temporary directory will be created for the process.
	// If prefixed with `./`, the path will be relative to the current working directory.
	// If prefixed with `//`, the path will be relative to the workspace root.
	// Environment variables in the path will be expended at runtime.
	Directory string `yaml:"dir,omitempty"`
}

func (e *ExecutableDirectory) ExpandDirectory(
	logger *io.Logger,
	wsPath, execPath, processTmpDir string,
	env map[string]string,
) (dir string, isTmpDir bool, err error) {
	if e.Directory == tmpdir {
		if processTmpDir != "" {
			return processTmpDir, true, nil
		}

		file, err := os.CreateTemp("", "flow")
		if err != nil {
			return "", false, err
		}
		return file.Name(), true, nil
	}

	return utils.ExpandDirectory(logger, e.Directory, wsPath, execPath, env), false, nil
}

type ExecutableEnvironment struct {
	// +docsgen:params
	// List of parameters to pass to the executable.
	Parameters ParameterList `yaml:"params,omitempty"`
	// +docgen:args
	// List of arguments to pass to the executable.
	Args ArgumentList `yaml:"args,omitempty"`
}

type ExecExecutableType struct {
	ExecutableDirectory   `yaml:",inline"`
	ExecutableEnvironment `yaml:",inline"`

	Command string  `yaml:"cmd,omitempty"`
	File    string  `yaml:"file,omitempty"`
	LogMode LogMode `yaml:"logMode,omitempty"`

	logFields map[string]interface{}
}

func (e *ExecExecutableType) SetLogFields(fields map[string]interface{}) {
	e.logFields = fields
}

func (e *ExecExecutableType) GetLogFields() map[string]interface{} {
	return e.logFields
}

type LaunchExecutableType struct {
	ExecutableEnvironment `yaml:",inline"`

	App  string `yaml:"app,omitempty"`
	URI  string `yaml:"uri,omitempty"`
	Wait bool   `yaml:"wait,omitempty"`
}

type RequestResponseFile struct {
	ExecutableDirectory `yaml:",inline"`

	Filename string `yaml:"filename,omitempty"`
	SaveAs   string `yaml:"saveAs,omitempty"`
}

type RequestExecutableType struct {
	ExecutableEnvironment `yaml:",inline"`

	Method  string            `yaml:"method,omitempty"`
	URL     string            `yaml:"url,omitempty"`
	Body    string            `yaml:"body,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
	Timeout time.Duration     `yaml:"timeout,omitempty"`

	ResponseFile      *RequestResponseFile `yaml:"responseFile,omitempty"`
	TransformResponse string               `yaml:"transformResponse,omitempty"`
	LogResponse       bool                 `yaml:"logResponse,omitempty"`
	ValidStatusCodes  []int                `yaml:"validStatusCodes,omitempty"`
}

type RenderExecutableType struct {
	ExecutableDirectory   `yaml:",inline"`
	ExecutableEnvironment `yaml:",inline"`

	TemplateFile     string `yaml:"templateFile,omitempty"`
	TemplateDataFile string `yaml:"templateDataFile,omitempty"`
}

type SerialExecutableType struct {
	ExecutableEnvironment `yaml:",inline"`

	// +docsgen:refs
	// List of executables references
	ExecutableRefs []Ref `yaml:"refs,omitempty"`
	FailFast       bool  `yaml:"failFast,omitempty"`
}

type ParallelExecutableType struct {
	ExecutableEnvironment `yaml:",inline"`

	ExecutableRefs []Ref `yaml:"refs,omitempty"`
	MaxThreads     int   `yaml:"maxThreads,omitempty"`
	FailFast       bool  `yaml:"failFast,omitempty"`
}

type ExecutableTypeSpec struct {
	// +docsgen:exec
	// Standard executable type. Runs a command/file in a subprocess.
	Exec *ExecExecutableType `yaml:"exec,omitempty"`
	// +docsgen:launch
	// Launches an application or opens a URI.
	Launch *LaunchExecutableType `yaml:"launch,omitempty"`
	// +docsgen:request
	// Makes an HTTP request.
	Request *RequestExecutableType `yaml:"request,omitempty"`
	// +docsgen:render
	// Renders a Markdown template with provided data. Requires the Interactive UI.
	Render *RenderExecutableType `yaml:"render,omitempty"`
	// +docsgen:serial
	// Runs a list of executables in serial.
	Serial *SerialExecutableType `yaml:"serial,omitempty"`
	// +docsgen:parallel
	// Runs a list of executables in parallel.
	Parallel *ParallelExecutableType `yaml:"parallel,omitempty"`
}

type Executable struct {
	Verb        Verb           `yaml:"verb,omitempty"`
	Name        string         `yaml:"name,omitempty"`
	Aliases     []string       `yaml:"aliases,omitempty"`
	Tags        Tags           `yaml:"tags,omitempty"`
	Description string         `yaml:"description,omitempty"`
	Visibility  VisibilityType `yaml:"visibility,omitempty"`
	Timeout     time.Duration  `yaml:"timeout,omitempty"`
	// +docsgen:type
	// The type of executable. Only one type can be set.
	Type *ExecutableTypeSpec `yaml:"type,omitempty"`

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

func (e *Executable) JSON() (string, error) {
	enriched := &enrichedExecutable{
		ID:   e.ID(),
		Spec: e,
	}
	jsonBytes, err := json.MarshalIndent(enriched, "", "  ")
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
		return "unk"
	}

	return NewExecutableID(e.workspace, e.namespace, e.Name)
}

func (e *Executable) Env() *ExecutableEnvironment {
	v := reflect.ValueOf(e.Type)
	if v.Kind() != reflect.Ptr {
		return nil
	}
	typeElem := v.Elem()
	for field := 0; field < typeElem.NumField(); field++ {
		if typeElem.Field(field).Kind() == reflect.Ptr && !typeElem.Field(field).IsNil() {
			elem := typeElem.Field(field).Elem()
			envField := elem.FieldByName("ExecutableEnvironment")
			if envField.IsValid() {
				return envField.Addr().Interface().(*ExecutableEnvironment)
			}
		}
	}
	return nil
}

func (e *Executable) AliasesIDs() []string {
	if len(e.Aliases) == 0 {
		return nil
	}

	if e.workspace == "" {
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
		e.Type.Request,
		e.Type.Render,
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

func (l ExecutableList) JSON() (string, error) {
	enriched := &enrichedExecutableList{}
	for _, exec := range l {
		enriched.Executables = append(enriched.Executables, &enrichedExecutable{
			ID:   exec.ID(),
			Spec: exec,
		})
	}
	jsonBytes, err := json.MarshalIndent(enriched, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal executable list - %w", err)
	}
	return string(jsonBytes), nil
}

func (l ExecutableList) Items() []*types.CollectionItem {
	items := make([]*types.CollectionItem, 0)
	for _, exec := range l {
		item := &types.CollectionItem{
			Header:    exec.ID(),
			SubHeader: exec.Verb.String(),
			Desc:      exec.Description,
		}
		if len(exec.Tags) > 0 {
			item.Desc = fmt.Sprintf("[%s]\n", exec.Tags.PreviewString()) + exec.Description
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
	return executables
}

func (l ExecutableList) FilterByNamespace(ns string) ExecutableList {
	if ns == "" || ns == "*" {
		return l
	}

	executables := lo.Filter(l, func(executable *Executable, _ int) bool {
		return executable.namespace == ns
	})
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
