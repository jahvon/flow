package request

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/itchyny/gojq"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/services/rest"
)

var log = io.Log()

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

func (r *requestRunner) Exec(_ *context.Context, executable *config.Executable) error {
	requestSpec := executable.Type.Request
	envMap, err := runner.ParametersToEnvMap(&requestSpec.ParameterizedExecutable)
	if err != nil {
		return fmt.Errorf("env setup failed - %w", err)
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
		return fmt.Errorf("request failed - %w", err)
	}

	if requestSpec.TransformResponse != "" {
		resp, err = executeJQQuery(requestSpec.TransformResponse, resp)
		if err != nil {
			return fmt.Errorf("jq execution failed - %w", err)
		}
	}

	if requestSpec.LogResponse {
		log.Info().Str("response", resp).Msgf("Successfully sent request to %s", requestSpec.URL)
	} else {
		log.Info().Msgf("Successfully sent request to %s", requestSpec.URL)
	}

	if requestSpec.ResponseFile != nil && requestSpec.ResponseFile.Filename != "" {
		targetDir, err := requestSpec.ResponseFile.ExpandDirectory(
			executable.WorkspacePath(),
			executable.DefinitionPath(),
			envMap,
		)
		if err != nil {
			return fmt.Errorf("unable to expand directory - %w", err)
		}
		defer requestSpec.ResponseFile.Finalize()

		err = writeResponseToFile(
			resp,
			filepath.Join(targetDir, requestSpec.ResponseFile.Filename),
			requestSpec.ResponseFile.SaveAs,
		)
		if err != nil {
			return fmt.Errorf("unable to write response to file - %w", err)
		} else {
			log.Info().Msgf("Successfully saved response to %s", requestSpec.ResponseFile.Filename)
		}
	}

	return nil
}

func isJSONString(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func executeJQQuery(query, resp string) (string, error) {
	var respMap map[string]interface{}
	err := json.Unmarshal([]byte(resp), &respMap)
	if err != nil {
		return "", fmt.Errorf("response is not a valid JSON string")
	}

	jqQuery, err := gojq.Parse(query)
	if err != nil {
		return "", err
	}

	iter := jqQuery.Run(respMap)
	result, ok := iter.Next()
	if !ok {
		return "", fmt.Errorf("unable to execute jq query")
	}
	if err, isErr := result.(error); isErr {
		return "", err
	}

	return fmt.Sprintf("%v", result), nil
}

func writeResponseToFile(resp, responseFile string, format config.OutputFormat) error {
	var formattedResp string
	switch format {
	case config.UNSET:
		formattedResp = resp
	case config.JSON:
		if !isJSONString(resp) {
			return fmt.Errorf("response is not a valid JSON string")
		}
		formattedResp = resp
	case config.FormattedJSON:
		var respMap map[string]interface{}
		err := json.Unmarshal([]byte(resp), &respMap)
		if err != nil {
			return fmt.Errorf("response is not a valid JSON string")
		}
		formattedStr, err := json.MarshalIndent(respMap, "", "  ")
		if err != nil {
			return err
		}
		formattedResp = string(formattedStr)
	case config.YAML:
		var respMap map[string]interface{}
		err := json.Unmarshal([]byte(resp), &respMap)
		if err != nil {
			return fmt.Errorf("response is not a valid JSON string")
		}
		yamlStr, err := yaml.Marshal(respMap)
		if err != nil {
			return err
		}
		formattedResp = string(yamlStr)
	default:
		return fmt.Errorf("unsupported output format - %s", format)
	}

	file, err := os.Create(responseFile)
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
