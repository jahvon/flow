package filesystem

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
)

func WriteExecutableDefinitionTemplate(templatePath string, template *config.ExecutableDefinitionTemplate) error {
	file, err := os.Create(filepath.Clean(templatePath))
	if err != nil {
		return errors.Wrap(err, "unable to create template file")
	}
	defer file.Close()

	if err := yaml.NewEncoder(file).Encode(template); err != nil {
		return errors.Wrap(err, "unable to encode template file")
	}
	return nil
}

func RenderAndWriteExecutablesTemplate(
	definitionTemplate *config.ExecutableDefinitionTemplate,
	ws *config.WorkspaceConfig,
	name, subPath string,
) error {
	if err := EnsureExecutableDir(ws.Location(), subPath); err != nil {
		return errors.Wrap(err, "unable to ensure existence of executable directory")
	}

	executablesPath := filepath.Join(ws.Location(), subPath)
	definitionYaml, err := yaml.Marshal(definitionTemplate.ExecutableDefinition)
	if err != nil {
		return errors.Wrap(err, "unable to marshal executable definition")
	}
	templateData := definitionTemplate.Data.MapInterface()
	templateData["Workspace"] = ws.AssignedName()
	templateData["WorkspaceLocation"] = ws.Location()
	templateData["ExecutablePath"] = executablesPath
	t, err := template.New("definition").Parse(string(definitionYaml))
	if err != nil {
		return errors.Wrap(err, "unable to parse definition template")
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, templateData); err != nil {
		return errors.Wrap(err, "unable to execute definition template")
	}

	filename := strings.ToLower(name)
	filename = strings.ReplaceAll(filename, " ", "_")
	if !strings.HasSuffix(filename, ExecutableDefinitionExt) {
		filename += ExecutableDefinitionExt
	}
	file, err := os.Create(filepath.Clean(filepath.Join(executablesPath, filename)))
	if err != nil {
		return errors.Wrap(err, "unable to create rendered definition file")
	}
	defer file.Close()

	if _, err := file.Write(buf.Bytes()); err != nil {
		return errors.Wrap(err, "unable to write rendered definition file")
	}

	if err := copyExecutableDefinitionTemplateAssets(definitionTemplate, executablesPath); err != nil {
		return errors.Wrap(err, "unable to copy template assets")
	}

	return nil
}

func LoadExecutableDefinitionTemplate(templateFile string) (*config.ExecutableDefinitionTemplate, error) {
	file, err := os.Open(filepath.Clean(templateFile))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open template file")
	}
	defer file.Close()

	definitionTemplate := &config.ExecutableDefinitionTemplate{}
	err = yaml.NewDecoder(file).Decode(definitionTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode template file")
	}
	definitionTemplate.SetContext(templateFile)

	return definitionTemplate, nil
}

func copyExecutableDefinitionTemplateAssets(
	definitionTemplate *config.ExecutableDefinitionTemplate,
	destinationPath string,
) error {
	sourcePath := filepath.Dir(definitionTemplate.Location())
	sourceFiles, err := expandArtifactFiles(sourcePath, definitionTemplate.Artifacts)
	if err != nil {
		return errors.Wrap(err, "unable to expand artifact files")
	}

	for _, file := range sourceFiles {
		relPath, err := filepath.Rel(sourcePath, file)
		if err != nil {
			return errors.Wrap(err, "unable to get relative path")
		}
		destPath := filepath.Join(destinationPath, relPath)
		if err := os.MkdirAll(filepath.Dir(destPath), 0750); err != nil {
			return errors.Wrap(err, "unable to create destination directory")
		}
		if err := CopyFile(file, destPath); err != nil {
			return errors.Wrap(err, "unable to copy file")
		}
	}
	return nil
}

func expandArtifactFiles(rootPath string, artifacts []string) ([]string, error) {
	var collectedFiles []string
	for _, file := range artifacts {
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
