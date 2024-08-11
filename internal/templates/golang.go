package templates

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jahvon/tuikit/components"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/utils"
	"github.com/jahvon/flow/types/executable"
	"github.com/jahvon/flow/types/workspace"
)

func ProcessTemplate(
	logger io.Logger,
	template *executable.Template,
	ws *workspace.Workspace,
	flowfileName, flowfilePath string,
) error {
	var data map[string]interface{}
	if template.Form != nil {
		if err := showForm(*template.Form); err != nil {
			return err
		}
		data = template.Form.MapInterface()
	}

	env := os.Environ()
	envMap := make(map[string]string)
	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		envMap[pair[0]] = pair[1]
	}
	flowfilePath = utils.ExpandDirectory(logger, flowfilePath, ws.Location(),  template.Location(), envMap)
	fullPath := filepath.Join(ws.Location(), flowfilePath)
	template.SetContext(fullPath)

	flowfile, err := templateToFlowfile(ws, template, flowfileName, flowfilePath, data)
	if err != nil {
		return err
	}

	if err := filesystem.WriteFlowFile(fullPath, flowfile); err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to write flowfile %s from template", flowfileName))
	}

	return nil
}

func showForm(fields executable.FormFields) error {
	if len(fields) == 0 {
		return nil
	}

	if err := fields.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid form fields: %v", err)
	}
	var inputs []*components.TextInput
	for _, f := range fields {
		inputs = append(inputs, &components.TextInput{
			Key:         f.Key,
			Prompt:      f.Prompt,
			Placeholder: f.Default,
		})
	}
	inputs, err := components.ProcessInputs(io.Theme(), inputs...)
	if err != nil {
		return fmt.Errorf("unable to show form: %v", err)
	}
	for _, input := range inputs {
		fields.Set(input.Key, input.Value())
	}
	if err := fields.ValidateValues(); err != nil {
		return err
	}
	return nil
}

func copyFlowFileTemplateArtifacts(t *executable.Template, cfgPath string) error {
	sourcePath := filepath.Dir(t.Location())
	sourceFiles, err := expandArtifactFiles(sourcePath, t.Artifacts)
	if err != nil {
		return errors.Wrap(err, "unable to expand artifact files")
	}

	for _, file := range sourceFiles {
		relPath, err := filepath.Rel(sourcePath, file)
		if err != nil {
			return errors.Wrap(err, "unable to get relative path")
		}
		destPath := filepath.Join(cfgPath, filepath.Base(relPath))
		if err := os.MkdirAll(filepath.Dir(destPath), 0750); err != nil {
			if !os.IsExist(err) {
				return errors.Wrap(err, "unable to create destination directory")
			}
			return errors.Wrap(err, "unable to create destination directory")
		}
		if err := CopyFile(file, destPath); err != nil {
			return errors.Wrap(err, "unable to copy file")
		}
	}
	return nil
}

func expandArtifactFiles(rootPath string, artifacts []executable.Artifact) ([]executable.Artifact, error) {
	var collectedFiles []string
	for _, file := range artifacts {
		if file.SrcDir != "" {

		}
		fullPath := filepath.Join(rootPath, file)
		//nolint:gocritic,nestif
		if info, err := os.Stat(fullPath); os.IsNotExist(err) {
			return nil, errors.Errorf("file does not exist: %s", fullPath)
		} else if err != nil {
			return nil, errors.Wrap(err, "unable to stat file")
		} else if info.IsDir() {
			err := filepath.WalkDir(fullPath, func(path string, entry fs.DirEntry, err error) error {
				if err != nil {
					return err
				} else if entry.IsDir() {
					return nil
				}
				collectedFiles = append(collectedFiles, path)
				return nil
			})
			if err != nil {
				return nil, errors.Wrap(err, "unable to walk directory")
			}
		} else {
			collectedFiles = append(collectedFiles, fullPath)
		}
	}
	return collectedFiles, nil
}


func templateToFlowfile(
	ws *workspace.Workspace,
	t *executable.Template,
	filename, path string,
	data map[string]interface{},
) (*executable.FlowFile, error) {
	data["FlowWorkspace"] = ws.AssignedName()
	data["FlowWorkspacePath"] = ws.Location()
	data["FlowFilePath"] = path
	tmpl, err := template.New(filename).Parse(t.Template)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to parse %s flowfile template in %s", filename, ws))
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to execute %s flowfile template in %s", filename, ws))
	}

	cfg := &executable.FlowFile{}
	if err := yaml.NewDecoder(&buf).Decode(cfg); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to decode %s flowfile template in %s", filename, ws)")
	}

	return cfg, nil
}

//
// func processAsGoTemplate(template *executable.Template, data map[string]interface{}) error {
//
// }
