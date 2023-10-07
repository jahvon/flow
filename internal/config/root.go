package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/backend"
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
	Backends         *backend.Config   `yaml:"backends"`

	currentWorkspaceConfig *workspace.Config
}

func (c *RootConfig) Validate() error {
	if c.CurrentWorkspace == "" {
		return fmt.Errorf("current workspace is not set")
	}
	if _, wsFound := c.Workspaces[c.CurrentWorkspace]; !wsFound {
		return fmt.Errorf("current workspace %s does not exist", c.CurrentWorkspace)
	}

	if c.Backends == nil {
		return fmt.Errorf("backends are not set")
	}
	if err := c.Backends.Validate(); err != nil {
		return err
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
		return fmt.Errorf("unable to load current workspace config - %v", err)
	} else if wsCfg == nil {
		return fmt.Errorf("current workspace config is nil")
	} else if err := wsCfg.Validate(); err != nil {
		return fmt.Errorf("encountered workspace config validation error - %v", err)
	}
	c.currentWorkspaceConfig = wsCfg
	return nil
}

func LoadConfig() *RootConfig {
	if err := common.EnsureDataDir(); err != nil {
		log.Fatal().Err(err).Msg("encountered issue with flow data directory")
	}

	var config *RootConfig
	file, err := os.Open(RootConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			config = defaultConfig()
			file, err = os.Create(RootConfigPath)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to create config file")
			}
			err = writeConfigFile(config)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to write config file")
			}
		} else {
			log.Fatal().Err(err).Msg("unable to open config file")
		}
	} else {
		config = &RootConfig{}
		err = yaml.NewDecoder(file).Decode(config)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to decode config file")
		}
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("unable to close config file")
		}
	}(file)

	if err := config.Validate(); err != nil {
		log.Fatal().Err(err).Msg("encountered validation error")
	}

	if err := config.setCurrentWorkspaceConfig(); err != nil {
		log.Fatal().Err(err).Msg("encountered issue setting current workspace config")
	}

	return config
}

func defaultConfig() *RootConfig {
	if err := workspace.SetDirectory(DefaultWorkspacePath); err != nil {
		log.Fatal().Err(err).Msg("encountered issue with workspace directory")
	}

	return &RootConfig{
		Workspaces:       map[string]string{"default": DefaultWorkspacePath},
		CurrentWorkspace: "default",
		Backends:         backend.NewConfig(),
	}
}
