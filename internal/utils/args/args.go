package args

import (
	"os"
	"strings"

	"github.com/jahvon/flow/types/executable"
)

func ParseArgs(args []string) (flagArgs map[string]string, posArgs []string) {
	flagArgs = make(map[string]string)
	posArgs = make([]string, 0)
	for i := 0; i < len(args); i++ {
		split := strings.Split(args[i], "=")
		if len(split) >= 2 {
			flagArgs[split[0]] = strings.Join(split[1:], "=")
			continue
		}
		posArgs = append(posArgs, args[i])
	}
	return
}

func ProcessArgs(
	executable *executable.Executable,
	execArgs []string,
	env map[string]string,
) (map[string]string, error) {
	if env != nil {
		// Expand environment variables in arguments
		for i, a := range execArgs {
			execArgs[i] = os.Expand(a, func(key string) string {
				return env[key]
			})
		}
	}
	flagArgs, posArgs := ParseArgs(execArgs)
	execEnv := executable.Env()
	if execEnv == nil || execEnv.Args == nil {
		return nil, nil //nolint:nilnil
	}
	if err := execEnv.Args.SetValues(flagArgs, posArgs); err != nil {
		return nil, err
	}
	if err := execEnv.Args.ValidateValues(); err != nil {
		return nil, err
	}
	return execEnv.Args.ToEnvMap(), nil
}
