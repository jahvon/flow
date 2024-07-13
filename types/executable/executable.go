package executable

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
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/errors"
	"github.com/jahvon/flow/internal/utils"
	"github.com/jahvon/flow/types/common"
)

//go:generate go run github.com/atombender/go-jsonschema@v0.16.0 -et --only-models -p executable -o executable.gen.go --capitalization URI --capitalization URL executable_schema.yaml

const (
	TmpDirLabel    = "f:tmp"
	DefaultTimeout = 30 * time.Minute
)

type ExecutableList []*Executable

func (e Directory) ExpandDirectory(
	logger io.Logger,
	wsPath, execPath, processTmpDir string,
	env map[string]string,
) (dir string, isTmpDir bool, err error) {
	if e == TmpDirLabel {
		if processTmpDir != "" {
			return processTmpDir, true, nil
		}

		td, err := os.MkdirTemp("", "flow")
		if err != nil {
			return "", false, err
		}
		return td, true, nil
	}

	return utils.ExpandDirectory(logger, string(e), wsPath, execPath, env), false, nil
}

type ExecutableEnvironment struct {
	Params ParameterList `json:"params" yaml:"params"`
	Args   ArgumentList  `json:"args"   yaml:"args"`
}

func (e *ExecExecutableType) SetLogFields(fields map[string]interface{}) {
	e.logFields = fields
}

func (e *ExecExecutableType) GetLogFields() map[string]interface{} {
	return e.logFields
}

type enrichedExecutableList struct {
	Executables []*enrichedExecutable `json:"executables" yaml:"executables"`
}
type enrichedExecutable struct {
	ID   string      `json:"id"   yaml:"id"`
	Spec *Executable `json:"spec" yaml:"spec"`
}

func (e *Executable) SetContext(workspaceName, workspacePath, namespace, flowFilePath string) {
	e.workspace = workspaceName
	e.workspacePath = workspacePath
	e.namespace = namespace
	e.flowFilePath = flowFilePath
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
	return execMarkdown(e)
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
	typeFields := []any{
		e.Exec,
		e.Launch,
		e.Request,
		e.Render,
		e.Serial,
		e.Parallel,
	}
	var execType any
	for _, field := range typeFields {
		if field == nil {
			continue
		}

		isPtr := reflect.ValueOf(field).Kind() == reflect.Ptr && !reflect.ValueOf(field).IsNil()
		isVal := reflect.ValueOf(field).Kind() != reflect.Ptr && !reflect.ValueOf(field).IsZero()
		if isPtr || isVal {
			execType = field
			break
		}
	}
	v := reflect.ValueOf(execType)
	if v.Kind() != reflect.Ptr {
		return nil
	}
	typeElem := v.Elem()
	execEnv := new(ExecutableEnvironment)
	for field := 0; field < typeElem.NumField(); field++ {
		if typeElem.Field(field).Kind() == reflect.Slice && !typeElem.Field(field).IsZero() {
			switch typeElem.Field(field).Interface().(type) {
			case ParameterList:
				execEnv.Params, _ = typeElem.Field(field).Interface().(ParameterList)
			case ArgumentList:
				execEnv.Args, _ = typeElem.Field(field).Interface().(ArgumentList)
			}
		}
	}
	return execEnv
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

func (e *Executable) FlowFilePath() string {
	return e.flowFilePath
}

func (e *Executable) SetDefaults() {
	if e.Verb == "" {
		e.Verb = "exec"
	}
	if e.Visibility == nil || *e.Visibility == "" {
		v := ExecutableVisibility(common.VisibilityPrivate)
		e.Visibility = &v
	}
	if e.Timeout == 0 {
		e.Timeout = DefaultTimeout
	}
}

func (e *Executable) Validate() error {
	if e.Verb == "" {
		return fmt.Errorf("verb cannot be empty")
	} else if err := e.Verb.Validate(); err != nil {
		return err
	}
	if e.Name == "" {
		return fmt.Errorf("name cannot be empty")
	} else if strings.Contains(e.Name, " ") {
		return fmt.Errorf("name cannot contain spaces")
	}

	err := utils.ValidateOneOf(
		"executable type",
		e.Exec,
		e.Launch,
		e.Request,
		e.Render,
		e.Serial,
		e.Parallel,
	)
	if err != nil {
		return err
	}

	if e.workspace == "" {
		return fmt.Errorf("workspace was not set")
	}
	if e.flowFilePath == "" {
		return fmt.Errorf("flowFile path was not set")
	}

	return nil
}

func (e *Executable) NameEquals(name string) bool {
	return e.Name == name || slices.Contains(e.Aliases, name)
}

func (e *Executable) MergeTags(tags common.Tags) {
	e.Tags = slices.Compact(append(e.Tags, tags...))
}

// IsVisibleFromWorkspace returns true if the executable should be shown in terminal output for the given workspace.
func (e *Executable) IsVisibleFromWorkspace(workspaceFilter string) bool {
	matchesWsFiler := e.workspace == workspaceFilter || workspaceFilter == "" || workspaceFilter == "*"
	if e.Visibility == nil {
		return matchesWsFiler
	}
	switch common.Visibility(*e.Visibility) {
	case common.VisibilityPrivate:
		return matchesWsFiler
	case common.VisibilityPublic:
		return true
	case common.VisibilityInternal, common.VisibilityHidden:
		return false
	default:
		return false
	}
}

// IsExecutableFromWorkspace returns true if the executable can be executed from the given workspace.
func (e *Executable) IsExecutableFromWorkspace(workspaceFilter string) bool {
	matchesWsFiler := e.workspace == workspaceFilter || workspaceFilter == "" || workspaceFilter == "*"
	if e.Visibility == nil {
		return matchesWsFiler
	}
	switch common.Visibility(*e.Visibility) {
	case common.VisibilityPrivate, common.VisibilityInternal:
		return matchesWsFiler
	case common.VisibilityPublic:
		return true
	case common.VisibilityHidden:
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
			Header: exec.Ref().String(),
			Desc:   exec.Description,
			ID:     exec.Ref().String(),
		}
		if t := common.Tags(exec.Tags); len(t) > 0 {
			item.Desc = fmt.Sprintf("[%s]\n", t.PreviewString()) + exec.Description
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
	var exec *Executable
	for _, e := range filteredList {
		if e.NameEquals(name) {
			exec = e
			break
		}
	}
	if exec != nil {
		return exec, nil
	}
	return nil, errors.ExecutableNotFoundError{Verb: string(verb), Name: name}
}

func (l ExecutableList) FilterByTags(tags common.Tags) ExecutableList {
	if len(tags) == 0 {
		return l
	}

	filteredExecs := make(ExecutableList, 0)
	for _, exec := range l {
		t := common.Tags(exec.Tags)
		if t.HasAnyTag(tags) {
			filteredExecs = append(filteredExecs, exec)
		}
	}
	return filteredExecs
}

func (l ExecutableList) FilterByVerb(verb Verb) ExecutableList {
	if verb == "" || verb == "*" {
		return l
	}

	if err := verb.Validate(); err != nil {
		return ExecutableList{}
	}

	filteredExecs := make(ExecutableList, 0)
	for _, exec := range l {
		if exec.Verb.Equals(verb) {
			filteredExecs = append(filteredExecs, exec)
		}
	}
	return filteredExecs
}

func (l ExecutableList) FilterBySubstring(str string) ExecutableList {
	if str == "" {
		return l
	}

	filteredExecs := make(ExecutableList, 0)
	for _, exec := range l {
		ref := exec.Ref().String()
		// TODO: Include aliases in search
		if strings.Contains(ref, str) || strings.Contains(exec.Description, str) {
			filteredExecs = append(filteredExecs, exec)
		}
	}
	return filteredExecs
}

func (l ExecutableList) FilterByWorkspace(ws string) ExecutableList {
	executables := make(ExecutableList, 0)
	for _, exec := range l {
		if exec.IsVisibleFromWorkspace(ws) {
			executables = append(executables, exec)
		}
	}

	if ws == "" || ws == "*" {
		return executables
	}

	filteredExecs := make(ExecutableList, 0)
	for _, exec := range executables {
		if exec.workspace == ws {
			filteredExecs = append(filteredExecs, exec)
		}
	}
	return filteredExecs
}

func (l ExecutableList) FilterByNamespace(ns string) ExecutableList {
	if ns == "" || ns == "*" {
		return l
	}

	filteredExecs := make(ExecutableList, 0)
	for _, exec := range l {
		if exec.namespace == ns {
			filteredExecs = append(filteredExecs, exec)
		}
	}
	return filteredExecs
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
	if namespace == "" || namespace == "*" {
		return fmt.Sprintf("%s/%s", workspace, name)
	}
	return fmt.Sprintf("%s/%s:%s", workspace, namespace, name)
}
