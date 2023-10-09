package run

import (
	"context"
	"fmt"
	"os"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"github.com/jahvon/flow/internal/io"
)

var log = io.Log()

func RunCmd(commandStr string) error {
	return RunCmdIn(commandStr, "")
}

func RunCmdIn(commandStr, dir string) error {
	log.Trace().Msgf("running command (%s) in dir (%s)", commandStr, dir)

	ctx := context.Background()
	parser := syntax.NewParser()
	reader := strings.NewReader(strings.TrimSpace(commandStr))
	prog, err := parser.Parse(reader, "")
	if err != nil {
		return fmt.Errorf("unable to parse command - %v", err)
	}

	runner, err := interp.New(
		interp.Dir(dir),
		// TODO: accept configured env vars
		interp.Env(expand.ListEnviron(os.Environ()...)),
		interp.StdIO(
			io.StdInReader{},
			io.StdOutWriter{LogAsDebug: false},
			io.StdErrWriter{LogAsDebug: false},
		),
	)
	err = runner.Run(ctx, prog)
	if err != nil {
		return fmt.Errorf("encountered an error executing command - %v", err)
	}
	return nil
}
