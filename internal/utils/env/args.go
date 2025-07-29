package env

import (
	"os"
	"strings"

	"github.com/flowexec/flow/types/executable"
)

func BuildArgsEnvMap(
	args executable.ArgumentList,
	execArgs []string,
	env map[string]string,
) (map[string]string, error) {
	al, err := resolveArgValues(args, execArgs, env)
	if err != nil {
		return nil, err
	}
	return argsToEnvMap(al), nil
}

func parseArgs(args []string) (flagArgs map[string]string, posArgs []string) {
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

func resolveArgValues(
	args executable.ArgumentList,
	execArgs []string,
	env map[string]string,
) (executable.ArgumentList, error) {
	if len(args) == 0 {
		return nil, nil
	}

	if env != nil {
		// Expand environment variables in arguments
		for i, a := range execArgs {
			execArgs[i] = os.Expand(a, func(key string) string {
				return env[key]
			})
		}
	}
	flagArgs, posArgs := parseArgs(execArgs)
	if err := setArgValues(args, flagArgs, posArgs, env); err != nil {
		return nil, err
	}
	return args, nil
}

func setArgValues(
	args executable.ArgumentList,
	flagArgs map[string]string,
	posArgs []string,
	env map[string]string,
) error {
	for i, arg := range args {
		if arg.EnvKey != "" {
			if val, found := env[arg.EnvKey]; found {
				// Use the input value if provided
				arg.Set(val)
				args[i] = arg
				continue
			}
		}

		if arg.Flag != "" {
			if val, ok := flagArgs[arg.Flag]; ok {
				arg.Set(val)
				args[i] = arg
			}
		} else if arg.Pos != nil && *arg.Pos != 0 {
			if *arg.Pos <= len(posArgs) {
				arg.Set(posArgs[*arg.Pos-1])
				args[i] = arg
			}
		}
	}
	return args.ValidateValues()
}

func argsToEnvMap(args executable.ArgumentList) map[string]string {
	envMap := make(map[string]string)
	for _, arg := range args {
		if arg.OutputFile != "" && arg.EnvKey == "" {
			continue
		}
		envMap[arg.EnvKey] = arg.Value()
	}
	return envMap
}

func filterArgsWithOutputFile(args executable.ArgumentList) executable.ArgumentList {
	var outputArgs executable.ArgumentList
	for _, arg := range args {
		if arg.OutputFile != "" {
			outputArgs = append(outputArgs, arg)
		}
	}

	return outputArgs
}
