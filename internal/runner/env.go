package runner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jahvon/tuikit/io"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/vault"
	"github.com/jahvon/flow/types/executable"
)

func SetEnv(logger io.Logger, exec *executable.ExecutableEnvironment, promptedEnv map[string]string) error {
	var errs []error
	for _, param := range exec.Params {
		val, err := ResolveParameterValue(logger, param, promptedEnv)
		if err != nil {
			errs = append(errs, err)
		}

		if err := os.Setenv(param.EnvKey, val); err != nil {
			return err
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to set values for parameters: %v", errs)
	}
	return nil
}

func ResolveParameterValue(
	logger io.Logger,
	param executable.Parameter,
	promptedEnv map[string]string,
) (string, error) {
	switch {
	case param.Text == "" && param.SecretRef == "" && param.Prompt == "":
		return "", nil
	case param.Text != "":
		return param.Text, nil
	case param.Prompt != "":
		val, ok := promptedEnv[param.EnvKey]
		if !ok {
			return "", errors.New("failed to get value for parameter")
		}
		return val, nil
	case param.SecretRef != "":
		if err := vault.ValidateReference(param.SecretRef); err != nil {
			return "", err
		}
		v := vault.NewVault(logger)
		secret, err := v.GetSecret(param.SecretRef)
		if err != nil {
			return "", err
		}
		return secret.PlainTextString(), nil
	default:
		return "", errors.New("failed to get value for parameter")
	}
}

func BuildEnvList(
	logger io.Logger,
	exec *executable.ExecutableEnvironment,
	inputEnv map[string]string,
	defaultEnv map[string]string,
) ([]string, error) {
	envList := make([]string, 0)
	var errs []error

	for k, v := range defaultEnv {
		if _, ok := inputEnv[k]; !ok {
			envList = append(envList, fmt.Sprintf("%s=%s", k, v))
		}
	}
	for _, param := range exec.Params {
		val, err := ResolveParameterValue(logger, param, inputEnv)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		envList = append(envList, fmt.Sprintf("%s=%s", param.EnvKey, val))
	}
	for _, arg := range exec.Args {
		envList = append(envList, fmt.Sprintf("%s=%s", arg.EnvKey, arg.Value()))
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to get values for parameters: %v", errs)
	}
	return envList, nil
}

func BuildEnvMap(
	logger io.Logger,
	exec *executable.ExecutableEnvironment,
	inputEnv map[string]string,
	defaultEnv map[string]string,
) (map[string]string, error) {
	envMap := make(map[string]string)
	var errs []error

	for k, v := range defaultEnv {
		if _, ok := envMap[k]; !ok {
			envMap[k] = v
		}
	}
	for _, param := range exec.Params {
		val, err := ResolveParameterValue(logger, param, inputEnv)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		envMap[param.EnvKey] = val
	}
	for _, arg := range exec.Args {
		envMap[arg.EnvKey] = arg.Value()
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to get values for parameters: %v", errs)
	}
	return envMap, nil
}

func DefaultEnv(ctx *context.Context, executable *executable.Executable) map[string]string {
	envMap := make(map[string]string)
	envMap["FLOW_RUNNER"] = "true"
	envMap["FLOW_CURRENT_WORKSPACE"] = ctx.CurrentWorkspace.AssignedName()
	envMap["FLOW_CURRENT_NAMESPACE"] = ctx.Config.CurrentNamespace
	if ctx.ProcessTmpDir != "" {
		envMap["FLOW_TMP_DIRECTORY"] = ctx.ProcessTmpDir
	}
	envMap["FLOW_EXECUTABLE_NAME"] = executable.Name
	envMap["FLOW_DEFINITION_PATH"] = executable.FlowFilePath()
	envMap["FLOW_DEFINITION_DIR"] = filepath.Dir(executable.FlowFilePath())
	envMap["FLOW_WORKSPACE_PATH"] = executable.WorkspacePath()
	envMap["FLOW_CONFIG_PATH"] = filesystem.ConfigDirPath()
	envMap["FLOW_CACHE_PATH"] = filesystem.CachedDataDirPath()
	envMap["DISABLE_FLOW_INTERACTIVE"] = "true"
	return envMap
}
