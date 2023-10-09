package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/common"
)

var (
	RootConfigPath       = common.ConfigDirPath() + "/config.yaml"
	DefaultWorkspacePath = common.CachedDataDirPath() + "/default"
)

func writeConfigFile(config *RootConfig) error {
	file, err := os.OpenFile(RootConfigPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("unable to open config file - %v", err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("unable to truncate config file - %v", err)
	}

	err = yaml.NewEncoder(file).Encode(config)
	if err != nil {
		return fmt.Errorf("unable to encode config file - %v", err)
	}

	return nil
}
