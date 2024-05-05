package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	stdio "io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/jahvon/tuikit/components"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/runner"
)

const appName = "flow renderer"

type renderRunner struct{}

func NewRunner() runner.Runner {
	return &renderRunner{}
}

func (r *renderRunner) Name() string {
	return "render"
}

func (r *renderRunner) IsCompatible(executable *config.Executable) bool {
	if executable == nil || executable.Type == nil || executable.Type.Render == nil {
		return false
	}
	return true
}

func (r *renderRunner) Exec(ctx *context.Context, executable *config.Executable, inputEnv map[string]string) error {
	if ctx.UserConfig.Interactive != nil && !ctx.UserConfig.Interactive.Enabled {
		return fmt.Errorf("unable to render when interactive mode is disabled")
	}

	renderSpec := executable.Type.Render
	if err := runner.SetEnv(ctx.Logger, &renderSpec.ExecutableEnvironment, inputEnv); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}
	envMap, err := runner.BuildEnvMap(
		ctx.Logger,
		&renderSpec.ExecutableEnvironment,
		inputEnv,
		runner.DefaultEnv(ctx, executable),
	)
	if err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}
	targetDir, isTmp, err := renderSpec.ExpandDirectory(
		ctx.Logger,
		executable.WorkspacePath(),
		executable.DefinitionPath(),
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

	tmpl, err := template.New(filepath.Base(renderSpec.TemplateFile)).
		Funcs(template.FuncMap{
			"env": func(key string) string {
				if val, ok := envMap[key]; ok {
					return val
				}
				return os.Getenv(key)
			},
		}).ParseFiles(contentFile)
	if err != nil {
		return errors.Wrapf(err, "unable to parse template file %s", contentFile)
	}

	var buff bytes.Buffer
	if err = tmpl.Execute(&buff, templateData); err != nil {
		return errors.Wrapf(err, "unable to execute template file %s", contentFile)
	}

	ctx.Logger.Infof("Rendering content from file %s", contentFile)
	filename := filepath.Base(contentFile)
	if err = components.RunMarkdownView(io.Theme(), appName, "file", filename, buff.String()); err != nil {
		return errors.Wrap(err, "unable to render content")
	}
	return nil
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
