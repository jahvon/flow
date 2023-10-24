package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	ConfigFileName = "workspace.yaml"
)

func writeConfigFile(workspacePath string, config *Config) error {
	file, err := os.Create(filepath.Clean(workspacePath + "/" + ConfigFileName))
	if err != nil {
		return fmt.Errorf("unable to create workspace config file - %w", err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("unable to truncate config file - %w", err)
	}

	err = yaml.NewEncoder(file).Encode(config)
	if err != nil {
		return fmt.Errorf("unable to encode workspace config file - %w", err)
	}

	return nil
}
