package filesystem

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/types/executable"
	"github.com/jahvon/flow/types/workspace"
)

func WriteFlowFileTemplate(templatePath string, template *executable.Template) error {
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

func WriteFlowFileFromTemplate(
	cfgTemplate *executable.Template,
	ws *workspace.Workspace,
	name, subPath string,
) error {
	if err := EnsureExecutableDir(ws.Location(), subPath); err != nil {
		return errors.Wrap(err, "unable to ensure existence of executable directory")
	}

	executablesPath := filepath.Join(ws.Location(), subPath)
	cfgYaml, err := yaml.Marshal(cfgTemplate.Template)
	if err != nil {
		return errors.Wrap(err, "unable to marshal flowfile template")
	}
	templateData := cfgTemplate.Form.MapInterface()
	templateData["Workspace"] = ws.AssignedName()
	templateData["WorkspaceLocation"] = ws.Location()
	templateData["ExecutablePath"] = executablesPath
	t, err := template.New("config").Parse(string(cfgYaml))
	if err != nil {
		return errors.Wrap(err, "unable to parse flowfile template")
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, templateData); err != nil {
		return errors.Wrap(err, "unable to execute config template")
	}

	filename := strings.ToLower(name)
	filename = strings.ReplaceAll(filename, " ", "_")
	if !strings.HasSuffix(filename, FlowFileExt) {
		filename += FlowFileExt
	}
	file, err := os.Create(filepath.Clean(filepath.Join(executablesPath, filename)))
	if err != nil {
		return errors.Wrap(err, "unable to create rendered config file")
	}
	defer file.Close()

	if _, err := file.Write(buf.Bytes()); err != nil {
		return errors.Wrap(err, "unable to write rendered config file")
	}

	if err := copyFlowFileTemplateAssets(cfgTemplate, executablesPath); err != nil {
		return errors.Wrap(err, "unable to copy template assets")
	}

	return nil
}

func LoadFlowFileTemplate(templateFile string) (*executable.FlowFileTemplate, error) {
	file, err := os.Open(filepath.Clean(templateFile))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open template file")
	}
	defer file.Close()

	cfgTemplate := &executable.FlowFileTemplate{}
	err = yaml.NewDecoder(file).Decode(cfgTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode template file")
	}
	cfgTemplate.SetContext(templateFile)

	return cfgTemplate, nil
}
