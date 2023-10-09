package workspace

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/errors"
	"github.com/jahvon/flow/internal/executable"
	"github.com/jahvon/flow/internal/executable/consts"
)

const (
	definitionExt = ".flow"
)

type Definition struct {
	Namespace   string          `yaml:"namespace"`
	Tags        []string        `yaml:"tags"`
	Executables executable.List `yaml:"executables"`

	workspace     string
	workspacePath string
}

func (d *Definition) SetContext(workspace, workspacePath string) {
	d.workspace = workspace
	d.workspacePath = workspacePath
	for _, exec := range d.Executables {
		exec.SetContext(workspace, workspacePath, d.Namespace)
	}
}

func (d *Definition) HasTag(tag string) bool {
	for _, t := range d.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

type DefinitionList []*Definition

func (l *DefinitionList) FilterByNamespace(namespace string) DefinitionList {
	var definitions []*Definition
	for _, definition := range *l {
		if definition.Namespace == namespace {
			definitions = append(definitions, definition)
		}
	}
	return definitions
}

func (l *DefinitionList) FilterByTag(tag string) DefinitionList {
	var definitions []*Definition
	for _, definition := range *l {
		for _, definitionTag := range definition.Tags {
			if definitionTag == tag {
				definitions = append(definitions, definition)
			}
		}
	}
	return definitions
}

// LookupExecutableByTypeAndName searches for an executable by type and name.
// If the executable is found, the namespace and executable are returned.
func (l *DefinitionList) LookupExecutableByTypeAndName(agent consts.AgentType, name string) (string, *executable.Executable, error) {
	for _, definition := range *l {
		exec, err := definition.Executables.FindByTypeAndName(agent, name)
		if err != nil {
			return "", nil, err
		} else if exec != nil {
			return definition.Namespace, exec, nil
		}
	}
	return "", nil, errors.ExecutableNotFound(agent, name)
}

func LoadDefinitions(workspace, workspacePath string) (DefinitionList, error) {
	definitionFiles, err := findDefinitionFiles(workspacePath)
	if err != nil {
		return nil, err
	}

	var definitions []*Definition
	for _, definitionFile := range definitionFiles {
		definition, err := loadDefinition(definitionFile)
		if err != nil {
			return nil, err
		}
		definition.SetContext(workspace, workspacePath)
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

func loadDefinition(definitionFile string) (*Definition, error) {
	file, err := os.Open(definitionFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open definition file - %v", err)
	}
	defer file.Close()

	config := &Definition{}
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode definition file - %v", err)
	}

	return config, nil
}

func findDefinitionFiles(root string) ([]string, error) {
	var definitions []string
	walkDirFunc := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(entry.Name()) == definitionExt {
			definitions = append(definitions, path)
		}
		return nil
	}
	if err := filepath.WalkDir(root, walkDirFunc); err != nil {
		return nil, err
	}
	return definitions, nil
}
