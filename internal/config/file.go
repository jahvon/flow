package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/tbox/internal/common"
)

const dataDirName = ".tbox"

var (
	RootConfigPath       = common.DataDirPath() + "/config.yaml"
	DefaultWorkspacePath = common.DataDirPath() + "/default"
)

func writeConfigFile(config *RootConfig) error {
	file, err := os.OpenFile(RootConfigPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("unable to open config file - %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			
		}
	}(file)

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("unable to truncate config file - %v", err)
	}

	err = yaml.NewEncoder(file).Encode(config)
	if err != nil {
		return fmt.Errorf("unable to encode config file - %v", err)
	}

	return nil
}
