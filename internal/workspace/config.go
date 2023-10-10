package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/executable"
)

type Config struct {
	ExecutablePreferences map[string]executable.Preference `yaml:"executablePreferences"`
}

func (c *Config) Validate() error {
	return nil
}

func LoadConfig(workspacePath string) (*Config, error) {
	file, err := os.Open(filepath.Clean(workspacePath + "/" + ConfigFileName))
	if err != nil {
		if os.IsNotExist(err) {
			if err := SetDirectory(workspacePath); err != nil {
				return nil, fmt.Errorf("unable to set workspace directory - %w", err)
			}
			return defaultConfig(), nil
		}
		return nil, fmt.Errorf("unable to open workspace config file - %w", err)
	}
	defer file.Close()

	config := &Config{}
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode workspace config file - %w", err)
	}

	return config, nil
}

func defaultConfig() *Config {
	return &Config{
		ExecutablePreferences: make(map[string]executable.Preference),
	}
}
