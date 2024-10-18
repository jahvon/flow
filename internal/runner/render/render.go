package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	stdio "io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/jahvon/tuikit/views"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/types/executable"
)

type renderRunner struct{}

func NewRunner() runner.Runner {
	return &renderRunner{}
}

func (r *renderRunner) Name() string {
	return "render"
}

func (r *renderRunner) IsCompatible(executable *executable.Executable) bool {
	if executable == nil || executable.Render == nil {
		return false
	}
	return true
}

func (r *renderRunner) Exec(ctx *context.Context, e *executable.Executable, inputEnv map[string]string) error {
	if !ctx.Config.ShowTUI() {
		return fmt.Errorf("unable to render when interactive mode is disabled")
	}

	renderSpec := e.Render
	if err := runner.SetEnv(ctx.Logger, e.Env(), inputEnv); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}
	envMap, err := runner.BuildEnvMap(ctx.Logger, e.Env(), inputEnv, runner.DefaultEnv(ctx, e))
	if err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}
	targetDir, isTmp, err := renderSpec.Dir.ExpandDirectory(
		ctx.Logger,
		e.WorkspacePath(),
		e.FlowFilePath(),
		ctx.ProcessTmpDir,
		envMap,
	)
	if err != nil {
		return errors.Wrap(err, "unable to expand directory")
	} else if isTmp {
		ctx.ProcessTmpDir = targetDir
	}

	contentFile := filepath.Clean(filepath.Join(targetDir, renderSpec.TemplateFile))
	var templateData map[string]interface{}
	if renderSpec.TemplateDataFile != "" {
		templateData, err = readDataFile(targetDir, renderSpec.TemplateDataFile)
		if err != nil {
			return err
		}
	}

	tmpl, err := template.New(filepath.Base(renderSpec.TemplateFile)).Funcs(sprig.TxtFuncMap()).ParseFiles(contentFile)
	if err != nil {
		return errors.Wrapf(err, "unable to parse template file %s", contentFile)
	}

	var buff bytes.Buffer
	if err = tmpl.Execute(&buff, templateData); err != nil {
		return errors.Wrapf(err, "unable to execute template file %s", contentFile)
	}

	ctx.Logger.Infof("Rendering content from file %s", contentFile)
	filename := filepath.Base(contentFile)
	ctx.TUIContainer.SetState("file", filename)
	return ctx.TUIContainer.SetView(views.NewMarkdownView(ctx.TUIContainer.RenderState(), buff.String()))
}

func readDataFile(dir, path string) (map[string]interface{}, error) {
	templateData := map[string]interface{}{}
	dataFilePath := filepath.Clean(filepath.Join(dir, path))
	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template data file %s does not exist", dataFilePath)
	}
	reader, err := os.Open(dataFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open template data file %s", dataFilePath)
	}
	defer reader.Close()
	data, err := stdio.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read template data file %s", dataFilePath)
	}
	extension := filepath.Ext(dataFilePath)
	switch extension {
	case ".json":
		if err = json.Unmarshal(data, &templateData); err != nil {
			return nil, errors.Wrapf(err, "unable to unmarshal template data file %s", dataFilePath)
		}
	case ".yaml", ".yml":
		if err = yaml.Unmarshal(data, &templateData); err != nil {
			return nil, errors.Wrapf(err, "unable to unmarshal template data file %s", dataFilePath)
		}
	}
	return templateData, nil
}
