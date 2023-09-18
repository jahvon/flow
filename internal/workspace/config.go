package workspace

import (
	"fmt"
	"os"

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
	file, err := os.Open(workspacePath + "/" + ConfigFileName)
	if err != nil {
		return nil, fmt.Errorf("unable to open workspace config file - %v", err)
	}
	defer file.Close()

	config := &Config{}
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode workspace config file - %v", err)
	}

	return config, nil
}

func defaultConfig() *Config {
	return &Config{
		ExecutablePreferences: make(map[string]executable.Preference),
	}
}
