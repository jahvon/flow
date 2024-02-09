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
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/runner"
)

var log = io.Log()

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

func (r *renderRunner) Exec(ctx *context.Context, executable *config.Executable, promptedEnv map[string]string) error {
	if ctx.InteractiveContainer == nil {
		return fmt.Errorf("unable to render when interactive mode is disabled")
	}

	renderSpec := executable.Type.Render
	if err := runner.SetEnv(&renderSpec.ParameterizedExecutable, promptedEnv); err != nil {
		return fmt.Errorf("env setup failed\n%w", err)
	}
	envMap, err := runner.ParametersToEnvMap(&renderSpec.ParameterizedExecutable, promptedEnv)
	if err != nil {
		return fmt.Errorf("env setup failed\n%w", err)
	}
	targetDir, isTmp, err := renderSpec.ExpandDirectory(
		executable.WorkspacePath(),
		executable.DefinitionPath(),
		ctx.ProcessTmpDir,
		envMap,
	)

	if err != nil {
		return fmt.Errorf("unable to expand directory\n%w", err)
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
		return fmt.Errorf("unable to parse template file %s\n%w", contentFile, err)
	}

	var buff bytes.Buffer
	if err = tmpl.Execute(&buff, templateData); err != nil {
		return fmt.Errorf("unable to execute template file %s\n%w", contentFile, err)
	}

	log.Info().Msgf("Rendering content from file %s", contentFile)
	state := &components.TerminalState{
		Theme:  io.Styles(),
		Width:  ctx.InteractiveContainer.Width(),
		Height: ctx.InteractiveContainer.Height(),
	}
	ctx.InteractiveContainer.SetView(components.NewMarkdownView(state, buff.String()))
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
		return nil, fmt.Errorf("unable to open template data file %s\n%w", dataFilePath, err)
	}
	defer reader.Close()
	data, err := stdio.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to read template data file %s\n%w", dataFilePath, err)
	}
	extension := filepath.Ext(dataFilePath)
	switch extension {
	case ".json":
		if err = json.Unmarshal(data, &templateData); err != nil {
			return nil, fmt.Errorf("unable to unmarshal template data file %s\n%w", dataFilePath, err)
		}
	case ".yaml", ".yml":
		if err = yaml.Unmarshal(data, &templateData); err != nil {
			return nil, fmt.Errorf("unable to unmarshal template data file %s\n%w", dataFilePath, err)
		}
	}
	return templateData, nil
}
