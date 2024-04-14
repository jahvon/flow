package run

import (
	"context"
	"fmt"
	stdio "io"
	"os"
	"path/filepath"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"github.com/jahvon/tuikit/io"

	"github.com/jahvon/flow/config"
)

// RunCmd executes a command in the current shell in a specific directory.
func RunCmd(
	commandStr, dir string,
	envList []string,
	logMode config.LogMode,
	logger io.Logger,
	logFields map[string]interface{},
) error {
	logger.Debugf("running command in dir (%s):\n%s", dir, strings.TrimSpace(commandStr))

	ctx := context.Background()
	parser := syntax.NewParser()
	reader := strings.NewReader(strings.TrimSpace(commandStr))
	prog, err := parser.Parse(reader, "")
	if err != nil {
		return fmt.Errorf("unable to parse command - %w", err)
	}

	if envList == nil {
		envList = make([]string, 0)
	}
	envList = append(os.Environ(), envList...)

	flattenedFields := make([]interface{}, 0)
	for k, v := range logFields {
		flattenedFields = append(flattenedFields, k, v)
	}
	runner, err := interp.New(
		interp.Dir(dir),
		interp.Env(expand.ListEnviron(envList...)),
		interp.StdIO(
			io.StdInReader{},
			stdOutWriter(logMode, logger, flattenedFields...),
			stdErrWriter(logMode, logger, flattenedFields...),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to create runner - %w", err)
	}

	err = runner.Run(ctx, prog)
	if err != nil {
		if code, isExit := interp.IsExitStatus(err); isExit {
			return fmt.Errorf("command exited with non-zero status %d", code)
		}
		return fmt.Errorf("encountered an error executing command - %w", err)
	}

	return nil
}

// RunFile executes a file in the current shell in a specific directory.
func RunFile(
	filename, dir string,
	envList []string,
	logMode config.LogMode,
	logger io.Logger,
	logFields map[string]interface{},
) error {
	logger.Debugf("executing file (%s)", filepath.Join(dir, filename))

	ctx := context.Background()
	fullPath := filepath.Join(dir, filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist - %s", fullPath)
	}
	file, err := os.OpenFile(filepath.Clean(fullPath), os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("unable to open file - %w", err)
	}
	defer file.Close()

	parser := syntax.NewParser()
	prog, err := parser.Parse(file, "")
	if err != nil {
		return fmt.Errorf("unable to parse file - %w", err)
	}

	if envList == nil {
		envList = make([]string, 0)
	}
	envList = append(os.Environ(), envList...)

	flattenedFields := make([]interface{}, 0)
	for k, v := range logFields {
		flattenedFields = append(flattenedFields, k, v)
	}
	runner, err := interp.New(
		interp.Env(expand.ListEnviron(envList...)),
		interp.StdIO(
			io.StdInReader{},
			stdOutWriter(logMode, logger, flattenedFields...),
			stdErrWriter(logMode, logger, flattenedFields...),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to create runner - %w", err)
	}

	err = runner.Run(ctx, prog)
	if err != nil {
		if code, isExit := interp.IsExitStatus(err); isExit {
			return fmt.Errorf("file execution exited with non-zero status %d", code)
		}
		return fmt.Errorf("encountered an error executing file - %w", err)
	}
	return nil
}

func stdOutWriter(mode config.LogMode, logger io.Logger, logFields ...any) stdio.Writer {
	switch mode {
	case config.NoLogMode:
		return stdio.Discard
	case config.StructuredLogMode:
		return io.StdOutWriter{LogFields: logFields, Logger: logger}
	case config.RawLogMode:
		return io.StdOutWriter{LogFields: logFields, Logger: logger, AsPlainText: true}
	default:
		logger.Errorx("unknown log mode", "mode", string(mode))
		return stdio.Discard
	}
}

func stdErrWriter(mode config.LogMode, logger io.Logger, logFields ...any) stdio.Writer {
	switch mode {
	case config.NoLogMode:
		return stdio.Discard
	case config.StructuredLogMode:
		return io.StdErrWriter{LogFields: logFields, Logger: logger}
	case config.RawLogMode:
		return io.StdErrWriter{LogFields: logFields, Logger: logger, AsPlainText: true}
	default:
		logger.Errorx("unknown log mode", "mode", string(mode))
		return stdio.Discard
	}
}
