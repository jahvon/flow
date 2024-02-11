package request

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/services/rest"
)

type requestRunner struct{}

func NewRunner() runner.Runner {
	return &requestRunner{}
}

func (r *requestRunner) Name() string {
	return "request"
}

func (r *requestRunner) IsCompatible(executable *config.Executable) bool {
	if executable == nil || executable.Type == nil || executable.Type.Request == nil {
		return false
	}
	return true
}

func (r *requestRunner) Exec(ctx *context.Context, executable *config.Executable, promptedEnv map[string]string) error {
	requestSpec := executable.Type.Request
	envMap, err := runner.ParametersToEnvMap(ctx.Logger, &requestSpec.ParameterizedExecutable, promptedEnv)
	if err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	url := expandEnvVars(envMap, requestSpec.URL)
	body := expandEnvVars(envMap, requestSpec.Body)
	for key, value := range requestSpec.Headers {
		requestSpec.Headers[key] = expandEnvVars(envMap, value)
	}
	restRequest := rest.Request{
		URL:     url,
		Method:  requestSpec.Method,
		Headers: requestSpec.Headers,
		Body:    body,
		Timeout: requestSpec.Timeout,
	}
	resp, err := rest.SendRequest(&restRequest, requestSpec.ValidStatusCodes)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}

	if requestSpec.TransformResponse != "" {
		resp, err = executeJQQuery(requestSpec.TransformResponse, resp)
		if err != nil {
			return errors.Wrap(err, "unable to transform response")
		}
	}

	logger := ctx.Logger
	if requestSpec.LogResponse {
		logger.Infox(fmt.Sprintf("Successfully sent request to %s", requestSpec.URL), "response", resp)
	} else {
		logger.Infof("Successfully sent request to %s", requestSpec.URL)
	}

	if requestSpec.ResponseFile != nil && requestSpec.ResponseFile.Filename != "" {
		targetDir, isTmp, err := requestSpec.ResponseFile.ExpandDirectory(
			ctx.Logger,
			executable.WorkspacePath(),
			executable.DefinitionPath(),
			ctx.ProcessTmpDir,
			envMap,
		)
		if err != nil {
			return errors.Wrap(err, "unable to expand directory")
		} else if isTmp {
			ctx.ProcessTmpDir = targetDir
		}

		err = writeResponseToFile(
			resp,
			filepath.Join(targetDir, requestSpec.ResponseFile.Filename),
			requestSpec.ResponseFile.SaveAs,
		)
		if err != nil {
			return errors.Wrap(err, "unable to save response")
		} else {
			logger.Infof("Successfully saved response to %s", requestSpec.ResponseFile.Filename)
		}
	}

	return nil
}

func executeJQQuery(query, resp string) (string, error) {
	var respMap map[string]interface{}
	err := json.Unmarshal([]byte(resp), &respMap)
	if err != nil {
		return "", errors.New("response is not a valid JSON string")
	}

	jqQuery, err := gojq.Parse(query)
	if err != nil {
		return "", err
	}

	iter := jqQuery.Run(respMap)
	result, ok := iter.Next()
	if !ok {
		return "", errors.New("unable to execute jq query")
	}
	if err, isErr := result.(error); isErr {
		return "", err
	}

	return fmt.Sprintf("%v", result), nil
}

func writeResponseToFile(resp, responseFile string, format string) error {
	var formattedResp string
	switch strings.ToLower(format) {
	case "", "raw":
		formattedResp = resp
	case "json":
		var js map[string]interface{}
		if json.Unmarshal([]byte(resp), &js) != nil {
			return errors.New("response is not a valid JSON string")
		}
		formattedResp = resp
	case "indent-json", "indented-json", "formatted-json":
		var respMap map[string]interface{}
		err := json.Unmarshal([]byte(resp), &respMap)
		if err != nil {
			return errors.New("response is not a valid JSON string")
		}
		formattedStr, err := json.MarshalIndent(respMap, "", "  ")
		if err != nil {
			return err
		}
		formattedResp = string(formattedStr)
	case "yaml", "yml":
		var respMap map[string]interface{}
		err := json.Unmarshal([]byte(resp), &respMap)
		if err != nil {
			return errors.New("response is not a valid JSON string")
		}
		yamlStr, err := yaml.Marshal(respMap)
		if err != nil {
			return err
		}
		formattedResp = string(yamlStr)
	default:
		return fmt.Errorf("unsupported output format - %s", format)
	}

	file, err := os.Create(filepath.Clean(responseFile))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(formattedResp)
	if err != nil {
		return err
	}
	return nil
}

func expandEnvVars(envMap map[string]string, value string) string {
	if envMap == nil || value == "" {
		return value
	}
	return os.Expand(value, func(envVar string) string {
		return envMap[envVar]
	})
}
