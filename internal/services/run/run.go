package run

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jahvon/flow/internal/io"
)

var log = io.Log()

func Run(commandStr string) error {
	return RunIn(commandStr, "")
}

func RunIn(commandStr, dir string) error {
	log.Trace().Msgf("running command (%s) in dir (%s)", commandStr, dir)

	splitCommandStr := strings.Split(commandStr, " ")
	if len(splitCommandStr) == 0 {
		log.Warn().Msgf("no command to execute")
		return nil
	} else if len(splitCommandStr) == 1 {
		return exec.Command(splitCommandStr[0]).Run()
	}
	cmd := exec.Command(splitCommandStr[0], splitCommandStr[1:]...)
	cmd.Stdout = io.StdOutWriter{LogAsDebug: false}
	cmd.Stderr = io.StdOutWriter{LogAsDebug: true}
	cmd.Stdin = io.StdInReader{}
	// TODO: accept configured env vars
	cmd.Env = os.Environ()
	if dir != "" {
		cmd.Dir = dir
	}
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("encountered an error executing command - %v", err)
	}
	return nil
}
