package run

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"github.com/jahvon/flow/internal/io"
)

var log = io.Log()

// RunCmd executes a command in the current shell in a specific directory.
func RunCmd(commandStr, dir string, envList []string) error {
	log.Trace().Msgf("running command (%s) in dir (%s)", commandStr, dir)

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

	runner, err := interp.New(
		interp.Dir(dir),
		interp.Env(expand.ListEnviron(envList...)),
		interp.StdIO(
			io.StdInReader{},
			io.StdOutWriter{LogAsDebug: false},
			io.StdErrWriter{LogAsDebug: false},
		),
	)
	if err != nil {
		return fmt.Errorf("unable to create runner - %w", err)
	}

	err = runner.Run(ctx, prog)
	if err != nil {
		return fmt.Errorf("encountered an error executing command - %w", err)
	}
	return nil
}

// RunFile executes a file in the current shell in a specific directory.
func RunFile(filename, dir string, envList []string) error {
	log.Trace().Msgf("executing file (%s)", filename)

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

	runner, err := interp.New(
		interp.Env(expand.ListEnviron(envList...)),
		interp.StdIO(
			io.StdInReader{},
			io.StdOutWriter{LogAsDebug: false},
			io.StdErrWriter{LogAsDebug: false},
		),
	)
	if err != nil {
		return fmt.Errorf("unable to create runner - %w", err)
	}

	err = runner.Run(ctx, prog)
	if err != nil {
		return fmt.Errorf("encountered an error executing file - %w", err)
	}
	return nil
}
