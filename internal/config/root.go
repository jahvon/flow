package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/common"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/workspace"
)

var (
	log = io.Log()
)

type RootConfig struct {
	Workspaces       map[string]string `yaml:"workspaces"`
	CurrentWorkspace string            `yaml:"currentWorkspace"`

	currentWorkspaceConfig *workspace.Config
}

func (c *RootConfig) Validate() error {
	if c.CurrentWorkspace == "" {
		return fmt.Errorf("current workspace is not set")
	}
	if _, wsFound := c.Workspaces[c.CurrentWorkspace]; !wsFound {
		return fmt.Errorf("current workspace %s does not exist", c.CurrentWorkspace)
	}

	return nil
}

func (c *RootConfig) CurrentWorkspaceConfig() *workspace.Config {
	return c.currentWorkspaceConfig
}

func (c *RootConfig) setCurrentWorkspaceConfig() error {
	wsPath, found := c.Workspaces[c.CurrentWorkspace]
	if !found {
		return fmt.Errorf("current workspace %s is not registered", c.CurrentWorkspace)
	}
	wsCfg, err := workspace.LoadConfig(wsPath)
	if err != nil {
		return fmt.Errorf("unable to load current workspace config - %w", err)
	} else if wsCfg == nil {
		return fmt.Errorf("current workspace config is nil")
	} else if err := wsCfg.Validate(); err != nil {
		return fmt.Errorf("encountered workspace config validation error - %w", err)
	}
	c.currentWorkspaceConfig = wsCfg
	return nil
}

func LoadConfig() *RootConfig {
	if err := common.EnsureConfigDir(); err != nil {
		log.Panic().Err(err).Msg("encountered issue with flow data directory")
	}

	var config *RootConfig
	file, err := os.Open(RootConfigPath)
	if err != nil { //nolint: nestif
		if os.IsNotExist(err) {
			config, err = initializeDefaultConfig()
			if err != nil {
				log.Panic().Err(err).Msg("unable to initialize default config")
			}
		} else {
			log.Panic().Err(err).Msg("unable to open config file")
		}
	} else {
		config = &RootConfig{}
		err = yaml.NewDecoder(file).Decode(config)
		if err != nil {
			log.Panic().Err(err).Msg("unable to decode config file")
		}
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Panic().Err(err).Msg("unable to close config file")
		}
	}(file)

	if err := config.Validate(); err != nil {
		log.Panic().Err(err).Msg("encountered validation error")
	}

	if err := config.setCurrentWorkspaceConfig(); err != nil {
		log.Panic().Err(err).Msg("encountered issue setting current workspace config")
	}

	return config
}

func defaultConfig() *RootConfig {
	if err := workspace.SetDirectory(DefaultWorkspacePath); err != nil {
		log.Panic().Err(err).Msg("encountered issue with workspace directory")
	}

	return &RootConfig{
		Workspaces:       map[string]string{"default": DefaultWorkspacePath},
		CurrentWorkspace: "default",
	}
}

func initializeDefaultConfig() (*RootConfig, error) {
	config := defaultConfig()
	_, err := os.Create(RootConfigPath)
	if err != nil {
		return nil, fmt.Errorf("unable to create config file - %w", err)
	}
	err = writeConfigFile(config)
	if err != nil {
		return nil, fmt.Errorf("unable to write config file - %w", err)
	}
	return config, nil
}
