package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/tbox/internal/io"
	"github.com/jahvon/tbox/internal/workspace"
)

const (
	dataDirName = ".tbox"
)

var (
	log                  = io.Log()
	RootConfigPath       = DataDirPath() + "/config.yaml"
	DefaultWorkspacePath = DataDirPath() + "/default"
)

type RootConfig struct {
	Workspaces       map[string]string `yaml:"workspaces"`
	CurrentWorkspace string            `yaml:"currentWorkspace"`
	GlobalVariables  map[string]string `yaml:"globalVariables"`
}

func LoadConfig() *RootConfig {
	if err := ensureDataDir(); err != nil {
		log.Fatal().Err(err).Msg("encountered issue with data directory")
	}

	file, err := os.Open(RootConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			config := defaultConfig()
			err = writeConfigFile(config)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to write default config file")
			}
			return config
		} else {
			log.Fatal().Err(err).Msg("unable to open config file")
		}
	}
	defer file.Close()

	config := &RootConfig{}
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to decode config file")
	}

	return config
}

func SetCurrentWorkspace(config *RootConfig, name string) error {
	if _, ok := config.Workspaces[name]; !ok {
		return fmt.Errorf("workspace %s does not exist", name)
	}

	config.CurrentWorkspace = name
	return writeConfigFile(config)
}

func AddWorkspace(config *RootConfig, name, location string) error {
	workspaceDir := filepath.Join(location, name)
	if err := workspace.CreateWorkspaceDirectory(workspaceDir); err != nil {
		return err
	}

	config.Workspaces[name] = workspaceDir
	return writeConfigFile(config)
}

func RemoveWorkspace(config *RootConfig, name string) error {
	if _, ok := config.Workspaces[name]; !ok {
		return fmt.Errorf("workspace %s does not exist", name)
	}

	delete(config.Workspaces, name)
	return writeConfigFile(config)
}

func DataDirPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to get home directory")
	}
	return filepath.Join(dirname, dataDirName)
}

func ensureDataDir() error {
	if _, err := os.Stat(DataDirPath()); os.IsNotExist(err) {
		err = os.MkdirAll(DataDirPath(), 0755)
		if err != nil {
			return fmt.Errorf("unable to create data directory - %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("unable to check for data directory - %v", err)
	}
	return nil
}

func defaultConfig() *RootConfig {
	if err := workspace.CreateWorkspaceDirectory(DefaultWorkspacePath); err != nil {
		log.Fatal().Err(err).Msg("encountered issue with workspace directory")
	}

	return &RootConfig{
		Workspaces:       map[string]string{"default": DefaultWorkspacePath},
		CurrentWorkspace: "default",
	}
}

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
