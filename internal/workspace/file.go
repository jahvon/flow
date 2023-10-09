package workspace

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ConfigFileName = "workspace.yaml"
)

func SetDirectory(location string) error {
	if info, err := os.Stat(location); os.IsNotExist(err) {
		err = os.MkdirAll(location, 0750)
		if err != nil {
			return fmt.Errorf("unable to create workspace directory - %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("unable to check for workspace directory - %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("workspace path (%s) exists but is not a directory", location)
	}

	if configInfo, err := os.Stat(location + "/" + ConfigFileName); os.IsNotExist(err) {
		config := defaultConfig()
		err = writeConfigFile(location, config)
		if err != nil {
			return fmt.Errorf("unable to write workspace config file - %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("unable to check for workspace config file - %w", err)
	} else if configInfo.IsDir() {
		return fmt.Errorf("workspace config file (%s) exists but is a directory", location+"/"+ConfigFileName)
	}

	return nil
}

func writeConfigFile(workspacePath string, config *Config) error {
	file, err := os.Create(workspacePath + "/" + ConfigFileName)
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
