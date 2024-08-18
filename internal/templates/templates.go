package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	tuikitIO "github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/utils"
	"github.com/jahvon/flow/types/executable"
	"github.com/jahvon/flow/types/workspace"
)

func ProcessTemplate(
	ctx *context.Context,
	template *executable.Template,
	ws *workspace.Workspace,
	flowfileName, flowfilePath string,
) error {
	logger := ctx.Logger
	if flowfileName == "" {
		flowfileName = fmt.Sprintf("executables_%s", time.Now().Format("20060102150405"))
	}
	flowfileName = strings.ReplaceAll(strings.ToLower(flowfileName), " ", "_")
	if !strings.HasSuffix(flowfileName, filesystem.FlowFileExt) {
		flowfileName += filesystem.FlowFileExt
	}

	data := make(map[string]interface{})
	if template.Form != nil {
		if err := showForm(template.Form); err != nil {
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
	flowfilePath = utils.ExpandDirectory(logger, flowfilePath, ws.Location(), template.Location(), envMap)
	fullPath := filepath.Join(ws.Location(), flowfilePath)
	template.SetContext(flowfileName, fullPath)
	logger.Debugx(
		fmt.Sprintf("processing %s template", flowfileName),
		"template", template.Location(), "output", fullPath,
	)

	data["FlowWorkspace"] = ws.AssignedName()
	data["FlowWorkspacePath"] = ws.Location()
	data["FlowFileName"] = template.Name()
	data["FlowFilePath"] = flowfilePath

	var preRun []executable.ExecExecutableType
	for _, e := range template.PreRun {
		preRun = append(preRun, executable.ExecExecutableType(e))
	}
	if err := runExecutables(ctx, "pre-run", flowfilePath, preRun, envMap); err != nil {
		return err
	}

	if err := copyAllArtifacts(logger, template.Artifacts, ws.Location(), fullPath, data, envMap); err != nil {
		return err
	}

	flowfile, err := templateToFlowfile(template, data)
	if err != nil {
		return err
	}

	if err := filesystem.WriteFlowFile(fullPath, flowfile); err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to write flowfile %s from template", flowfileName))
	}

	var postRun []executable.ExecExecutableType
	for _, e := range template.PostRun {
		postRun = append(postRun, executable.ExecExecutableType(e))
	}
	if err := runExecutables(ctx, "post-run", flowfilePath, postRun, envMap); err != nil {
		return err
	}

	return nil
}

func runExecutables(
	ctx *context.Context,
	stage, flowfileDir string,
	execs []executable.ExecExecutableType,
	envMap map[string]string,
) error {
	ctx.Logger.Debugf("running %d %s executables", len(execs), stage)
	for i, exec := range execs {
		exec.SetLogFields(map[string]interface{}{
			"stage": stage,
			"step":  i + 1,
		})
		eCopy := exec
		e := executable.Executable{
			Verb: "exec",
			Name: fmt.Sprintf("%s-exec-%d", stage, i),
			Exec: &eCopy,
		}
		e.SetContext(
			ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
			"flow-internal", flowfileDir,
		)
		if err := runner.Exec(ctx, &e, envMap); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to execute %s executable %d", stage, i))
		}
	}
	return nil
}

func parseSourcePath(
	logger tuikitIO.Logger,
	name, flowfileDir, wsDir string,
	artifact executable.Artifact,
	data map[string]interface{},
	envMap map[string]string,
) (string, error) {
	var err error
	srcDir := flowfileDir
	if artifact.SrcDir != "" {
		srcDir = utils.ExpandDirectory(logger, artifact.SrcDir, wsDir, flowfileDir, envMap)
	}
	var sb *bytes.Buffer
	sb, err = processAsGoTemplate(name, filepath.Join(srcDir, artifact.SrcName), data)
	if err != nil {
		return "", errors.Wrap(err, "unable to process artifact as template")
	}
	return sb.String(), nil
}

func parseDestinationPath(
	logger tuikitIO.Logger,
	name, flowfileDir, relPath, wsDir string,
	artifact executable.Artifact,
	data map[string]interface{},
	envMap map[string]string,
) (string, error) {
	var err error
	dstDir := filepath.Dir(filepath.Join(flowfileDir, relPath))
	if artifact.DstDir != nil {
		dstDir = utils.ExpandDirectory(logger, *artifact.DstDir, wsDir, flowfileDir, envMap)
	}
	dstName := artifact.DstName
	var db *bytes.Buffer
	db, err = processAsGoTemplate(name, dstName, data)
	if err != nil {
		return "", errors.Wrap(err, "unable to process artifact as template")
	}
	dstName = db.String()
	return filepath.Join(dstDir, dstName), nil
}

func templateToFlowfile(
	t *executable.Template,
	data map[string]interface{},
) (*executable.FlowFile, error) {
	buf, err := processAsGoTemplate(t.Name(), t.Template, data)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("flowfile template %s", t.Name()))
	}

	flowfile := &executable.FlowFile{}
	if err := yaml.NewDecoder(buf).Decode(flowfile); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to decode %s flowfile template", t.Name()))
	}

	return flowfile, nil
}

func processAsGoTemplate(fileName, txt string, data map[string]interface{}) (*bytes.Buffer, error) {
	tmpl, err := template.New(fileName).Funcs(sprig.FuncMap()).Parse(txt)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to parse %s template", fileName))
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to execute %s template", fileName))
	}

	return &buf, nil
}

func goTemplateEvaluatedTrue(fileName, txt string, data map[string]interface{}) (bool, error) {
	t, err := template.New(fileName).Funcs(sprig.FuncMap()).Parse(txt)
	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("unable to parse %s template", fileName))
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return false, errors.Wrap(err, "unable to evaluate template")
	}
	return strconv.ParseBool(buf.String())
}
