package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/common"
)

var (
	RootConfigPath       = common.ConfigDirPath() + "/config.yaml"
	DefaultWorkspacePath = common.CachedDataDirPath() + "/default"
)

func writeConfigFile(config *RootConfig) error {
	file, err := os.OpenFile(filepath.Clean(RootConfigPath), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("unable to open config file - %w", err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("unable to truncate config file - %w", err)
	}

	err = yaml.NewEncoder(file).Encode(config)
	if err != nil {
		return fmt.Errorf("unable to encode config file - %w", err)
	}

	return nil
}
