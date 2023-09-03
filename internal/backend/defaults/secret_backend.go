package defaults

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/jahvon/flow/internal/backend"
	"github.com/jahvon/flow/internal/backend/consts"
	"github.com/jahvon/flow/internal/common"
)

const envFileSuffix = ".scrts"

type EnvFileSecretBackend struct {
}

func SecretBackend() backend.SecretBackend {
	return &EnvFileSecretBackend{}
}

func (d *EnvFileSecretBackend) Name() consts.BackendName {
	return consts.EnvFileBackendName
}

func (d *EnvFileSecretBackend) InitializeBackend() error {
	return nil
}

func (d *EnvFileSecretBackend) GetSecret(context, key string) (backend.Secret, error) {
	log.Trace().Msgf("getting secret for env file context %s and key %s", context, key)
	envFilePath, err := contextToEnvFilePath(context)
	if err != nil {
		return "", fmt.Errorf("failed to get env file path: %w", err)
	}

	envMap, err := envFileToMap(envFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to get env file map: %w", err)
	}

	value, ok := envMap[key]
	if !ok {
		return "", fmt.Errorf("key not found in env file: %s", key)
	}

	return backend.Secret(value), nil
}

func (d *EnvFileSecretBackend) SetSecret(context, key string, secret backend.Secret) error {
	log.Trace().Msgf("updating secret for env file context %s and key %s", context, key)
	envFilePath, err := contextToEnvFilePath(context)
	if err != nil {
		return fmt.Errorf("failed to get env file path: %w", err)
	}

	envMap, err := envFileToMap(envFilePath)
	if err != nil {
		return fmt.Errorf("failed to get env file map: %w", err)
	}

	envMap[key] = string(secret)

	if err := mapToEnvFile(envMap, envFilePath); err != nil {
		return fmt.Errorf("failed to write to env file: %w", err)
	}

	return nil
}

func (d *EnvFileSecretBackend) DeleteSecret(context, key string) error {
	log.Trace().Msgf("deleting secret for env file context %s and key %s", context, key)
	envFilePath, err := contextToEnvFilePath(context)
	if err != nil {
		return fmt.Errorf("failed to get env file path: %w", err)
	}

	envMap, err := envFileToMap(envFilePath)
	if err != nil {
		return fmt.Errorf("failed to get env file map: %w", err)
	}

	delete(envMap, key)

	if err := mapToEnvFile(envMap, envFilePath); err != nil {
		return fmt.Errorf("failed to write to env file: %w", err)
	}

	return nil
}

func contextToEnvFilePath(context string) (string, error) {
	if context == "" {
		return "", errors.New("missing context")
	}

	if err := common.EnsureDataDir(); err != nil {
		return "", err
	}

	return filepath.Join(common.DataDirPath(), context+envFileSuffix), nil
}

func envFileToMap(filePath string) (map[string]string, error) {
	if filePath == "" {
		return nil, errors.New("missing env file path")
	}
	envMap := make(map[string]string)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		key, value, found := strings.Cut(line, "=")
		if !found {
			return nil, fmt.Errorf("invalid line in env file: %s", line)
		}
		envMap[key] = value
	}

	return envMap, nil
}

func mapToEnvFile(envMap map[string]string, filePath string) error {
	if filePath == "" {
		return errors.New("missing env file path")
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open env file: %w", err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("unable to truncate env file - %v", err)
	}

	for key, value := range envMap {
		_, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return fmt.Errorf("failed to write to env file: %w", err)
		}
	}

	return nil
}
