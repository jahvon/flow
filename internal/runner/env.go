package runner

import (
	"errors"
	"fmt"
	"os"

	"github.com/samber/lo"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/vault"
)

func SetEnv(exec *config.ParameterizedExecutable) error {
	var errs []error
	for _, param := range exec.Parameters {
		val, err := ResolveParameterValue(param)
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

func ResolveParameterValue(param config.Parameter) (string, error) {
	switch {
	case param.Text == "" && param.SecretRef == "" && param.Prompt == "":
		return "", nil
	case param.Text != "":
		return param.Text, nil
	case param.Prompt != "":
		response := io.Ask(param.Prompt)
		return response, nil
	case param.SecretRef != "":
		if err := vault.ValidateReference(param.SecretRef); err != nil {
			return "", err
		}
		v := vault.NewVault()
		secret, err := v.GetSecret(param.SecretRef)
		if err != nil {
			return "", err
		}
		return secret.PlainTextString(), nil
	default:
		return "", errors.New("failed to get value for parameter")
	}
}

func ParametersToEnvList(exec *config.ParameterizedExecutable) ([]string, error) {
	params := exec.Parameters
	var errs []error
	env := lo.Map(params, func(param config.Parameter, _ int) string {
		key := param.EnvKey
		value, err := ResolveParameterValue(param)
		if err != nil {
			errs = append(errs, err)
			return ""
		}
		return fmt.Sprintf("%s=%s", key, value)
	})

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to get values for parameters: %v", errs)
	}
	return env, nil
}

func ParametersToEnvMap(exec *config.ParameterizedExecutable) (map[string]string, error) {
	params := exec.Parameters
	var errs []error
	env := lo.SliceToMap(params, func(param config.Parameter) (string, string) {
		val, err := ResolveParameterValue(param)
		if err != nil {
			errs = append(errs, err)
			return param.EnvKey, ""
		}

		return param.EnvKey, val
	})

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to get values for parameters: %v", errs)
	}
	return env, nil
}
