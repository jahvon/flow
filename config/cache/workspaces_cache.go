package cache

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/io"
)

const wsCacheKey = "workspace"

type WorkspaceCacheData struct {
	// Map of workspace name to workspace config
	Workspaces map[string]*config.WorkspaceConfig `yaml:"workspaces"`
	// Map of workspace name to workspace path
	WorkspaceLocations map[string]string `yaml:"workspaceLocations"`
}

type WorkspaceCache struct {
	Data *WorkspaceCacheData
}

var (
	log = io.Log().With().Str("scope", "discovery/cache").Logger()
)

func NewWorkspaceCache() *WorkspaceCache {
	if workspaceCache == nil {
		workspaceCache = &WorkspaceCache{}
	}
	return workspaceCache
}

func (c *WorkspaceCache) Update() error {
	if c.Data == nil {
		log.Debug().Msg("Initializing workspace cache data")
	} else {
		log.Debug().Msg("Updating workspace cache data")
	}

	userCfg := file.LoadUserConfig()
	if userCfg == nil {
		return fmt.Errorf("failed to load user config")
	}

	cacheData := &WorkspaceCacheData{
		Workspaces:         make(map[string]*config.WorkspaceConfig),
		WorkspaceLocations: make(map[string]string),
	}
	for name, path := range userCfg.Workspaces {
		wsCfg, err := file.LoadWorkspaceConfig(name, path)
		if err != nil {
			return fmt.Errorf("failed loading workspace config: %w", err)
		} else if wsCfg == nil {
			log.Err(err).Msgf("config not found for workspace %s", name)
			continue
		}
		cacheData.Workspaces[name] = wsCfg
		cacheData.WorkspaceLocations[name] = path
	}
	log.Trace().Int("workspaces", len(cacheData.Workspaces)).Msg("Successfully loaded workspace configs")

	data, err := yaml.Marshal(cacheData)
	if err != nil {
		return fmt.Errorf("unable to encode cache data - %w", err)
	}

	err = file.WriteLatestCachedData(wsCacheKey, data)
	if err != nil {
		return fmt.Errorf("unable to write cache data - %w", err)
	}

	c.Data = cacheData
	log.Trace().Msg("Successfully updated workspace cache data")

	return nil
}

func (c *WorkspaceCache) Get() (*WorkspaceCacheData, error) {
	if c.Data != nil {
		return c.Data, nil
	}

	cacheData, err := file.LoadLatestCachedData(wsCacheKey)
	if err != nil {
		return nil, fmt.Errorf("unable to load workspace cache data - %w", err)
	} else if cacheData == nil {
		if err := c.Update(); err != nil {
			return nil, fmt.Errorf("unable to get updated workspace cache data - %w", err)
		}
	}

	c.Data = &WorkspaceCacheData{}
	if err := yaml.Unmarshal(cacheData, c.Data); err != nil {
		return nil, fmt.Errorf("unable to decode workspace cache data - %w", err)
	}
	log.Trace().Int("workspaces", len(c.Data.Workspaces)).Msg("Fetched workspace cache data")

	return c.Data, nil
}
